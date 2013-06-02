package model

import (
	"math/rand"
)

type Guest struct {
}

func random(min, max int) int {
	if max - min == 0 {
		return min
	}
	return rand.Intn(max - min) + min
}

func (g *Guest) BuyOffer(platforms []*Platform) (offer *Offer, platform *Platform) {
	if len(platforms) == 0 {
		return
	}
	platformI := random(0, len(platforms))
	platform = platforms[platformI]
	if len(platform.offers) == 0 {
		return
	}
	offerI := random(0, len(platform.offers))
	offer = platform.offers[offerI]
	return
}

func MakeGuest() *Guest {
	return &Guest {}
}
