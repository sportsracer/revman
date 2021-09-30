package model

import (
	"errors"
	"fmt"
	"math/rand"
	"sort"
)

const (
	// hotels are sorted by criteria other than price, too ...
	priceRandomness = 1
)

// Platform is a marketplace where hotels can offer their rooms for a price
type Platform struct {
	Name   string
	offers []*Offer
}

// Add an offer. Returns an error if the hotel attempts to market more rooms than it has
func (p *Platform) AddOffer(offer *Offer) error {
	if p.countOffers(offer.Hotel) > offer.Hotel.Rooms {
		return errors.New("Too many offers!")
	}
	p.offers = append(p.offers, offer)
	return nil
}

// Buy an offer. The offer gets removed, and the offer price (minus variable cost) granted to the hotel balance
func (p *Platform) BuyOffer(offer *Offer) error {
	for i, _offer := range p.offers {
		if offer == _offer {
			// remove offer
			p.offers[i] = p.offers[len(p.offers)-1]
			p.offers = p.offers[0 : len(p.offers)-1]
			// wire the money!
			offer.Hotel.Balance += offer.Price - VariableCost
			return nil
		}
	}
	return errors.New("Offer not found!")
}

// Remove all offers
func (p *Platform) Reset() {
	p.offers = make([]*Offer, 0, 128)
}

// Iterable that filters offers by hotel
type offersByHotel struct {
	offers []*Offer
	hotel  *Hotel
}

func (o *offersByHotel) Iter() <-chan interface{} {
	ch := make(chan interface{})

	go (func() {
		for _, offer := range o.offers {
			if offer.Hotel == o.hotel {
				ch <- offer
			}
		}
		close(ch)
	})()

	return ch
}

// How many offers does this hotel advertise on this platform?
func (p *Platform) countOffers(hotel *Hotel) int {
	offers := offersByHotel{p.offers, hotel}
	c := 0
	for range offers.Iter() {
		c++
	}
	return c
}

// Get the first offer by this hotel, or nil
func (p *Platform) getOfferByHotel(hotel *Hotel) *Offer {
	offers := offersByHotel{p.offers, hotel}
	for offer := range offers.Iter() {
		return offer.(*Offer)
	}
	return nil
}

// Get the (up to) `num` cheapest offers. There's some randomness in the sorting to simulate human irrationality
func (p *Platform) getCheapestOffers(num int) []*Offer {
	offers := make([]*Offer, len(p.offers))
	copy(offers, p.offers)
	sort.Slice(offers, func(i, j int) bool {
		p1 := offers[i].Price * (1.0 + rand.Float32()*priceRandomness)
		p2 := offers[j].Price * (1.0 + rand.Float32()*priceRandomness)
		return p1 < p2
	})
	if num >= len(offers) {
		return offers
	}
	return offers[:num]
}

func (p *Platform) GetOffers() []*Offer {
	return p.offers
}

func MakePlatform(name string) *Platform {
	p := &Platform{
		Name: name,
	}
	p.Reset()
	return p
}

func (p *Platform) String() string {
	return fmt.Sprintf("Platform (name: %s, offers: %d)", p.Name, len(p.offers))
}
