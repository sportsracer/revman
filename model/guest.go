package model

import (
	"fmt"
	"math/rand"
	"sort"
	//"errors"
	
	"me/revman/util"
)

const (
	// how many recently seen/bought offers does the guest remember? 
	cognitiveLoad = 10
	// weight of bought/seen offers and external value in perceived value calculation
	boughtWeight = 0.5
	seenWeight = 0.25
	externalWeight = 0.25
	// we're not fully rational ;)
	priceRandomness = 0.10
)

type Guest struct {
	// guest's perceived value of a hotel room
	externalValue float32
	// offers this guest has bought
	boughtOffers []*Offer
	// offers this guest has seen
	seenOffers []*Offer
	// number of platforms this guest usually checks
	numPlatforms int
	// number of offers this guest usually looks at per platform
	numOffers int
	// increases likelihood that this guest returns to the same hotel
	loyalty float32
}

func lastOffers(offers []*Offer, num int) []*Offer {
	var i = len(offers) - num
	if i < 0 {
		i = 0
	}
	return offers[i:]
}

// make []*Offer iterable

type offerSlice []*Offer

func (xs offerSlice) iterate(out chan<- interface{}) {
	for _, x := range xs {
		out <- x
	}
	close(out)
}

func (xs offerSlice) Iter() <-chan interface{} {
	ch := make(chan interface{})
	go xs.iterate(ch)
	return ch
}

func (g *Guest) getPreferredHotel() *Hotel {
	
	getHotel := func(o interface{}) interface{} {
		return o.(*Offer).Hotel
	}
	var offers offerSlice = lastOffers(g.boughtOffers, cognitiveLoad)
	hotels := util.Map(getHotel, offers)

	hotel, count := util.MaxOccur(hotels)

	if count <= 1 {
		return nil
	}
	return hotel.(*Hotel)
}

func (g *Guest) getPreferredPlatform() *Platform {

	getPlatform := func(o interface{}) interface{} {
		return o.(*Offer).Platform
	}
	var offers offerSlice = lastOffers(g.boughtOffers, cognitiveLoad)
	platforms := util.Map(getPlatform, offers)

	platform, count := util.MaxOccur(platforms)

	if count <= 1 {
		return nil
	}
	return platform.(*Platform)
}

func (g *Guest) GetValue() float32 {

	var sum float32
	var totalWeight float32

	getPrice := func(o interface{}) interface{} {
		return o.(*Offer).Price
	}

	var boughtOffers offerSlice = lastOffers(g.boughtOffers, cognitiveLoad)
	var avgBought, boughtErr = util.Avg(util.Map(getPrice, boughtOffers))
	if boughtErr == nil {
		sum += avgBought * boughtWeight
		totalWeight += boughtWeight
	}

	var seenOffers offerSlice = lastOffers(g.seenOffers, cognitiveLoad)
	var avgSeen, seenErr = util.Avg(util.Map(getPrice, seenOffers))
	if seenErr == nil {
		sum += avgSeen * seenWeight
		totalWeight += seenWeight
	}

	{
		sum += g.externalValue * externalWeight
		totalWeight += externalWeight
	}

	return sum / totalWeight
}

func shufflePlatforms(platforms []*Platform) []*Platform {
	indices := rand.Perm(len(platforms))
	shuffled := make([]*Platform, len(platforms))
	for i, v := range indices {
		shuffled[v] = platforms[i]
	}
	return shuffled
}

func getAdjustedPrice(o *Offer, value float32, preferredHotel *Hotel, loyalty float32) float32 {
	price := o.Price
	if (o.Hotel == preferredHotel) {
		price /= loyalty
	}
	price -= o.Hotel.getReputationValue()
	price *= (1.0 - priceRandomness) + rand.Float32() * priceRandomness // spice it up a bit ;)
	return price
}

// make offers sortable by value for this guest

type offerByAdjustedPrice struct {
	offers []*Offer
	value float32
	preferredHotel *Hotel
	loyalty float32
}

func (o *offerByAdjustedPrice) Len() int {
	return len(o.offers)
}

func (o *offerByAdjustedPrice) Less(i, j int) bool {
	v1 := getAdjustedPrice(o.offers[i], o.value, o.preferredHotel, o.loyalty)
	v2 := getAdjustedPrice(o.offers[j], o.value, o.preferredHotel, o.loyalty)
	return v1 < v2
}

func (o *offerByAdjustedPrice) Swap(i, j int) {
	o.offers[i], o.offers[j] = o.offers[j], o.offers[i]
}

func (g *Guest) BuyOffer(platforms []*Platform) *Offer {
	
	// choose platforms
	consideredPlatforms := make([]*Platform, 0)
	if preferredPlatform := g.getPreferredPlatform(); preferredPlatform != nil {
		consideredPlatforms = append(consideredPlatforms, preferredPlatform)
	}
	platforms = shufflePlatforms(platforms)
	for i := len(consideredPlatforms); i < g.numPlatforms && i < len(platforms); i++ {
		consideredPlatforms = append(consideredPlatforms, platforms[i])
	}

	// choose offers
	consideredOffers := make([]*Offer, 0)
	// add offers of preferred hotel on all considered platforms
	preferredHotel := g.getPreferredHotel()
	if preferredHotel != nil {
		for _, platform := range consideredPlatforms {
			if offer := platform.getOfferByHotel(preferredHotel); offer != nil {
				consideredOffers = append(consideredOffers, offer)
			}
		}
	}
	// add cheapest x offers per platform
	for _, platform := range consideredPlatforms {
		offers := platform.getCheapestOffers(g.numOffers)
		consideredOffers = append(consideredOffers, offers...)
	}

	g.seenOffers = append(g.seenOffers, consideredOffers...)

	value := g.GetValue()
	maxPrice := value * 1.1
	filterOffer := func(o interface{}) bool {
		return o.(*Offer).Price <= maxPrice
	}

	acceptableOffers := make([]*Offer, 0)
	for _, offer := range consideredOffers {
		if filterOffer(offer) {
			acceptableOffers = append(acceptableOffers, offer)
		}
	}
	
	offers := &offerByAdjustedPrice {
		offers: acceptableOffers,
		value: g.GetValue(),
		preferredHotel: preferredHotel,
		loyalty: g.loyalty,
	}
	sort.Sort(offers)

	if len(offers.offers) == 0 {
		return nil
	}
	g.boughtOffers = append(g.boughtOffers, offers.offers[0])
	return offers.offers[0]
}

func MakeGuest() *Guest {
	return &Guest {
		externalValue: 100.0,
		boughtOffers: make([]*Offer, 0),
		seenOffers: make([]*Offer, 0),
		numPlatforms: 1,
		numOffers: 10,
		loyalty: 1.1,
	}
}

func (g *Guest) String() string {
	return fmt.Sprintf("Guest (numPlatforms: %d, numOffers: %d)", g.numPlatforms, g.numOffers)
}
