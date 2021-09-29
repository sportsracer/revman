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

type Platform struct {
	Name   string
	offers []*Offer
}

func (p *Platform) countOffers(hotel *Hotel) int {
	c := 0
	for _, offer := range p.offers {
		if offer.Hotel == hotel {
			c++
		}
	}
	return c
}

func (p *Platform) AddOffer(offer *Offer) error {
	if p.countOffers(offer.Hotel) > offer.Hotel.Rooms {
		return errors.New("Too many offers!")
	}
	p.offers = append(p.offers, offer)
	return nil
}

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

func (p *Platform) Reset() {
	p.offers = make([]*Offer, 0, 128)
}

func (p *Platform) getOfferByHotel(hotel *Hotel) *Offer {
	for _, offer := range p.offers {
		if offer.Hotel == hotel {
			return offer
		}
	}
	return nil
}

// wrapper to make offers sortable by price

type offerByPrice struct {
	offers []*Offer
}

func (o *offerByPrice) Len() int {
	return len(o.offers)
}

func (o *offerByPrice) Less(i, j int) bool {
	p1 := o.offers[i].Price * (1.0 + rand.Float32()*priceRandomness)
	p2 := o.offers[j].Price * (1.0 + rand.Float32()*priceRandomness)
	return p1 < p2
}

func (o *offerByPrice) Swap(i, j int) {
	o.offers[i], o.offers[j] = o.offers[j], o.offers[i]
}

func (p *Platform) getCheapestOffers(num int) []*Offer {
	cp := make([]*Offer, len(p.offers))
	copy(cp, p.offers)
	offers := &offerByPrice{
		offers: cp,
	}
	sort.Sort(offers)
	if num >= len(offers.offers) {
		num = len(offers.offers)
	}
	return offers.offers[:num]
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
