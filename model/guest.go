package model

import (
	"fmt"
	"math/rand"
	"sort"

	"github.com/sportsracer/revman/util"
)

const (
	// how many recently seen/bought offers does the guest remember?
	cognitiveLoad = 200
	// weight of bought/seen offers and external value in perceived value calculation
	boughtWeight = 0.05
	seenWeight   = 0.2
	valueWeight  = 0.75
	// we're not fully rational ;)
	valueRandomness = 0.1
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

// Model loyalty to a specific hotel. Return this guest's preferred hotel, if they have one. They need to have stayed at this hotel a lot
func (g *Guest) getPreferredHotel() *Hotel {
	offers := util.TakeFirst(util.Iter(g.boughtOffers), cognitiveLoad)
	hotels := util.Map(offers, func(o *Offer) *Hotel {
		return o.Hotel
	})
	hotel, count := util.MaxOccur(hotels)

	if count <= cognitiveLoad/5 {
		return nil
	}
	return *hotel
}

// Model loyalty to a specific platform. Return this guest's preferred platform, if they have one.
func (g *Guest) getPreferredPlatform() *Platform {
	offers := util.TakeFirst(util.Iter(g.boughtOffers), cognitiveLoad)
	platforms := util.Map(offers, func(o *Offer) *Platform {
		return o.Platform
	})
	platform, count := util.MaxOccur(platforms)

	if count <= cognitiveLoad/5 {
		return nil
	}
	return *platform
}

// Get this guest's adjusted value perception, based on offers they have seen or bought.
func (g *Guest) GetValue() float32 {
	var sum, totalWeight float32

	getPrice := func(o *Offer) float32 {
		return o.Price
	}

	boughtOffers := util.TakeFirst(util.Iter(g.boughtOffers), cognitiveLoad)
	avgBought, boughtErr := util.Avg(util.Map(boughtOffers, getPrice))
	if boughtErr == nil {
		sum += avgBought * boughtWeight
		totalWeight += boughtWeight
	}

	seenOffers := util.TakeFirst(util.Iter(g.seenOffers), cognitiveLoad)
	avgSeen, seenErr := util.Avg(util.Map(seenOffers, getPrice))
	if seenErr == nil {
		sum += avgSeen * seenWeight
		totalWeight += seenWeight
	}

	sum += g.value * valueWeight
	totalWeight += valueWeight
	return sum / totalWeight
}

// Return slice of `num` randomly picked platforms
func takeRandomPlatforms(platforms []*Platform, num int) []*Platform {
	if num > len(platforms) {
		num = len(platforms)
	}
	shuffled := make([]*Platform, num)

	indices := rand.Perm(len(platforms))
	for i, j := range indices[:num] {
		shuffled[i] = platforms[j]
	}

	return shuffled
}

// Calculate price adjusted by loyalty, if the offer is for the preferred hotel
func getAdjustedPrice(o *Offer, value float32, preferredHotel *Hotel, loyalty float32) float32 {
	price := o.Price
	if o.Hotel == preferredHotel {
		price /= loyalty
	}
	price *= (1.0 - valueRandomness) + rand.Float32()*valueRandomness*2.0 // spice it up a bit ;)
	return price
}

func (g *Guest) BuyOffer(platforms []*Platform) *Offer {
	// choose platforms
	consideredPlatforms := takeRandomPlatforms(platforms, g.numPlatforms)
	if preferredPlatform := g.getPreferredPlatform(); preferredPlatform != nil {
		consideredPlatforms = append(consideredPlatforms, preferredPlatform)
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

	// only consider offers below our acceptable price
	maxPrice := g.value
	if g.maxValue < maxPrice {
		maxPrice = g.maxValue
	}
	acceptableOffers := make([]*Offer, 0)
	for _, offer := range consideredOffers {
		if getAdjustedPrice(offer, g.value, preferredHotel, g.loyalty) <= maxPrice {
			acceptableOffers = append(acceptableOffers, offer)
		}
	}

	sort.Slice(acceptableOffers, func(i, j int) bool {
		v1 := getAdjustedPrice(acceptableOffers[i], g.value, preferredHotel, g.loyalty)
		v2 := getAdjustedPrice(acceptableOffers[j], g.value, preferredHotel, g.loyalty)
		return v1 < v2
	})

	if len(acceptableOffers) == 0 {
		return nil
	}
	g.boughtOffers = append(g.boughtOffers, acceptableOffers[0])
	return acceptableOffers[0]
}

func MakeGuest() *Guest {
	return &Guest{
		value:        100.0,
		maxValue:     100.0 + rand.Float32()*100.0,
		boughtOffers: make([]*Offer, 0),
		seenOffers:   make([]*Offer, 0),
		numPlatforms: rand.Intn(4) + 1,
		numOffers:    rand.Intn(5) + 5,
		loyalty:      1.0 + rand.Float32()*0.2,
	}
}

func (g *Guest) String() string {
	return fmt.Sprintf("Guest (numPlatforms: %d, numOffers: %d)", g.numPlatforms, g.numOffers)
}
