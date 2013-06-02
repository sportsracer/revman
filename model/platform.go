package model

import (
	"errors"
	"fmt"
)

type Platform struct {
	Name string
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
			p.offers = p.offers[0:len(p.offers)-1]
			// wire the money!
			offer.Hotel.Balance += offer.Price
			return nil
		}
	}
	return errors.New("Offer not found!")
}

func (p *Platform) Reset() {
	p.offers = make([]*Offer, 0, 128)
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
