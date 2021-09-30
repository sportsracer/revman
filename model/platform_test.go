package model

import "testing"

const roomPrice = 80

func makeHotelPlatformOffer() (*Hotel, *Platform, *Offer) {
	hotel := MakeHotel()
	platform := MakePlatform("ta")
	offer := MakeOffer(hotel, roomPrice, platform)
	return hotel, platform, offer
}

func TestAddOffer(t *testing.T) {
	_, platform, offer := makeHotelPlatformOffer()

	platform.AddOffer(offer)
	count := len(platform.GetOffers())

	if expectedCount := 1; expectedCount != count {
		t.Errorf("Expected %d offers, got %d", expectedCount, count)
	}
}

func TestBuyOffer(t *testing.T) {
	t.Run("wires money to the hotel", func(t *testing.T) {
		hotel, platform, offer := makeHotelPlatformOffer()
		initialBalance := hotel.Balance

		platform.AddOffer(offer)
		platform.BuyOffer(offer)

		if expectedBalance := initialBalance + roomPrice - VariableCost; expectedBalance != hotel.Balance {
			t.Errorf("Expected balance of %f, got %f", expectedBalance, hotel.Balance)
		}
	})

	t.Run("returns an error if offer not found", func(t *testing.T) {
		_, platform, offer := makeHotelPlatformOffer()

		err := platform.BuyOffer(offer)

		if err == nil {
			t.Error("Expected error here")
		}
	})
}

func TestGetCheapestOffers(t *testing.T) {
	hotel, platform, offer := makeHotelPlatformOffer()
	cheaperPrice := float32(roomPrice - priceRandomness*10)
	cheaperOffer := MakeOffer(hotel, cheaperPrice, platform)

	platform.AddOffer(offer)
	platform.AddOffer(cheaperOffer)
	num := 1
	offers := platform.getCheapestOffers(num)

	if len(offers) != num {
		t.Errorf("Expected to get %d offer", num)
	}

	if offers[0] != cheaperOffer {
		t.Errorf("Expected %s, got %s", cheaperOffer, offers[0])
	}
}

func TestCountOffers(t *testing.T) {
	hotel, platform, offer := makeHotelPlatformOffer()
	anotherHotel := MakeHotel()
	anotherOffer := MakeOffer(anotherHotel, 70, platform)

	platform.AddOffer(offer)
	platform.AddOffer(anotherOffer)
	count := platform.countOffers(hotel)

	if expectedCount := 1; count != expectedCount {
		t.Errorf("Expected to get %d offer", expectedCount)
	}
}
