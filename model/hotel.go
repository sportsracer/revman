package model

import (
	"fmt"
)

type Hotel struct {
	Rooms int
	Balance float32
	TrustScore int
}

func (h *Hotel) getReputationValue() float32 {
	return (float32(h.TrustScore) - 70.0) * 0.25
}

func MakeHotel() *Hotel {
	return &Hotel{
		Rooms: 100,
		Balance: 0,
		TrustScore: 80,
	}
}

func (h *Hotel) String() string {
	return fmt.Sprintf("Hotel (rooms: %d, balance: %f, TrustScore: %d)", h.Rooms, h.Balance, h.TrustScore)
}
