package model

import (
	"fmt"
)

type Offer struct {
	Hotel    *Hotel
	Price    float32
	Platform *Platform
}

func MakeOffer(hotel *Hotel, price float32, platform *Platform) *Offer {
	return &Offer{
		Hotel:    hotel,
		Price:    price,
		Platform: platform,
	}
}

func (o *Offer) String() string {
	return fmt.Sprintf("Offer (hotel: %s, price: %f, platform: %s)", o.Hotel, o.Price, o.Platform)
}
