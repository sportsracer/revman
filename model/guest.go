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
	cognitiveLoad = 50
	// weight of bought/seen offers and external value in perceived value calculation
	boughtWeight = 0.05
	seenWeight = 0.2
	valueWeight = 0.75
	// we're not fully rational ;)
	priceRandomness = 0.05
)

type Guest struct {
	// guest's perceived value of a hotel room
	value float32
	// maximum price this guest is willing to pay, based on economic situation, comparable prices etc
	maxValue float32
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

	if count <= cognitiveLoad / 5 {
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

	if count <= cognitiveLoad / 5 {
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
		sum += g.value * valueWeight
		totalWeight += valueWeight
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

	// update value
	g.value = g.GetValue()

	filterOffer := func(o interface{}) bool {
		max := func() float32 { if g.maxValue < g.value { return g.maxValue }; return g.value }()
		return getAdjustedPrice(o.(*Offer), g.value, preferredHotel, g.loyalty) <= max
	}
	acceptableOffers := make([]*Offer, 0)
	for _, offer := range consideredOffers {
		if filterOffer(offer) {
			acceptableOffers = append(acceptableOffers, offer)
		}
	}
	
	offers := &offerByAdjustedPrice {
		offers: acceptableOffers,
		value: g.value,
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
		value: 100.0,
		maxValue: 90.0 + rand.Float32() * 110.0,
		boughtOffers: make([]*Offer, 0),
		seenOffers: make([]*Offer, 0),
		numPlatforms: rand.Intn(4) + 1,
		numOffers: rand.Intn(5) + 5,
		loyalty: 1.0 + rand.Float32() * 0.1,
	}
}

func (g *Guest) String() string {
	return fmt.Sprintf("Guest (numPlatforms: %d, numOffers: %d)", g.numPlatforms, g.numOffers)
}
