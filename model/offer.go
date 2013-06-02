package model

import (
	"fmt"
)

type Offer struct {
	Hotel *Hotel
	Price float32
}

func MakeOffer(hotel *Hotel, price float32) *Offer {
	return &Offer{
		Hotel: hotel,
		Price: price,
	}
}

func (o *Offer) String() string {
	return fmt.Sprintf("Offer (hotel: %s, price: %f)", o.Hotel, o.Price)
}
