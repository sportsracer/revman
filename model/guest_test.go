package model

import "testing"

func makeHotelPlatformGuest() (*Hotel, *Platform, *Guest) {
	hotel := MakeHotel()
	platform := MakePlatform("ta")
	guest := MakeGuest()
	return hotel, platform, guest
}

func TestGetValue(t *testing.T) {
	hotel, platform, guest := makeHotelPlatformGuest()
	initialValue := guest.value
	for i := 0; i < 10; i++ {
		platform.AddOffer(MakeOffer(hotel, initialValue-10, platform))
	}

	guest.BuyOffer([]*Platform{platform})

	if guest.value >= initialValue {
		t.Errorf("Expected guest value to drop below %f, but was %f", initialValue, guest.value)
	}
}

func TestPlatformBuyOffer(t *testing.T) {
	t.Run("guest chooses cheapest among acceptable offers", func(t *testing.T) {
		hotel, platform, guest := makeHotelPlatformGuest()
		offer := MakeOffer(hotel, guest.value-10, platform)
		platform.AddOffer(offer)
		cheaperOffer := MakeOffer(hotel, offer.Price-50, platform)
		platform.AddOffer(cheaperOffer)

		boughtOffer := guest.BuyOffer([]*Platform{platform})

		if cheaperOffer != boughtOffer {
			t.Errorf("Expected guest to buy %s, but they bought %s instead", cheaperOffer, boughtOffer)
		}
	})

	t.Run("guest does not buy rooms at unacceptable prices", func(t *testing.T) {
		hotel, platform, guest := makeHotelPlatformGuest()
		expensiveOffer := MakeOffer(hotel, guest.maxValue+guest.value, platform)
		platform.AddOffer(expensiveOffer)

		boughtOffer := guest.BuyOffer([]*Platform{platform})

		if boughtOffer != nil {
			t.Errorf("Expected guest not to buy an offer, but bought %s", boughtOffer)
		}
	})
}
