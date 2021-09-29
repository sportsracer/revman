package model

import (
	"fmt"
)

const (
	// cost of renting out a room
	VariableCost = 50
	// fix cost of operating for a day
	FixCost = 2000
)

type Hotel struct {
	Rooms   int
	Balance float32
}

func MakeHotel() *Hotel {
	return &Hotel{
		Rooms:   100,
		Balance: 0,
	}
}

func (h *Hotel) String() string {
	return fmt.Sprintf("Hotel (rooms: %d, balance: %f)", h.Rooms, h.Balance)
}
