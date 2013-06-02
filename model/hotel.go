package model

import (
	"fmt"
)

type Hotel struct {
	Rooms int
	Balance float32
}

func MakeHotel() *Hotel {
	return &Hotel{
		Rooms: 100,
		Balance: 0,
	}
}

func (h *Hotel) String() string {
	return fmt.Sprintf("Hotel (rooms: %d, balance: %f)", h.Rooms, h.Balance)
}
