package ctrl

import (
	"errors"
	"log"

	"github.com/sportsracer/revman/model"
)

type Game struct {
	players map[int]*Player
	hotels map[int]*model.Hotel
	platforms map[string]*model.Platform
	guests []*model.Guest
}

func (g *Game) AddPlayer(player *Player) {
	g.players[player.id] = player
	hotel := model.MakeHotel()
	g.hotels[player.id] = hotel

	// add some more guests to the pool
	numNewGuests := int(float32(hotel.Rooms) * 0.80)
	for i := 0; i < numNewGuests; i++ {
		guest := model.MakeGuest()
		g.guests = append(g.guests, guest)
	}
	log.Printf("Game: %s joined, created %s, added %d guests", player, hotel, numNewGuests)
}

func (g *Game) AddOffer(player *Player, platformName string, price float32) error {
	hotel := g.hotels[player.id]
	platform, ok := g.platforms[platformName]
	if !ok {
		return errors.New("Unknown platform")
	}
	offer := model.MakeOffer(hotel, price, platform)
	err := platform.AddOffer(offer)
	return err
}

// Perform an iteration, asking guests to choose offers. Return nested map
// containing all bought offers
func (g *Game) Tick() (statuses map[int]map[string]interface{}, players []map[string]interface{}, offers []map[string]interface{}) {
	log.Printf("Game: Starting round (hotels: %d, guests: %d)", len(g.hotels), len(g.guests))
	log.Println("Game: Offers:")

	platforms := make([]*model.Platform, 0, len(g.platforms))
	for _, platform := range g.platforms {
		log.Printf("\t%s", platform)
		platforms = append(platforms, platform)
	}

	// offers
	offers = make([]map[string]interface{}, 0)
	for _, platform := range g.platforms {
		for _, offer := range platform.GetOffers() {
			var id int = -1
			for _id, hotel := range g.hotels {
				if hotel == offer.Hotel {
					id = _id
					break
				}
			}
			if id == -1 {
				log.Printf("Hotel not found :(")
				continue
			}
			offerMap := map[string]interface{}{
				"platform": offer.Platform.Name,
				"price": offer.Price,
				"player": g.players[id].name,
			}
			offers = append(offers, offerMap)
		}
	}

	boughtPerPlayer := make(map[int][]*map[string]interface{})
	for _, guest := range g.guests {
		offer := guest.BuyOffer(platforms)
		if offer == nil {
			continue
		}
		if err := offer.Platform.BuyOffer(offer); err != nil {
			log.Panic("Game: Error in transaction for %s", offer)
			continue
		}

		// store state
		id := -1
		for _id, hotel := range g.hotels {
			if offer.Hotel == hotel {
				id = _id
				break
			}
		}
		if id == -1 {
			log.Panic("Hotel %s for offer not found", offer.Hotel)
		}
		offerStatus := map[string]interface{}{"platform": offer.Platform.Name, "price": offer.Price}
		if _, ok := boughtPerPlayer[id]; !ok {
			boughtPerPlayer[id] = []*map[string]interface{}{&offerStatus}
		} else {
			boughtPerPlayer[id] = append(boughtPerPlayer[id], &offerStatus)
		}
	}

	for _, hotel := range g.hotels {
		hotel.Balance -= model.FixCost
	}

	log.Println("Game: Ending round")
	log.Println("Game: Hotels:")
	for playerId, hotel := range g.hotels {
		log.Printf("\t%d: %s", playerId, hotel)
	}

	// reset offers on platforms
	for _, platform := range g.platforms {
		platform.Reset()
	}
	
	// return status
	statuses = make(map[int]map[string]interface{})
	for id, hotel := range g.hotels {
		bought, ok := boughtPerPlayer[id]
		if !ok {
			bought = []*map[string]interface{}{}
		}
		statuses[id] = map[string]interface{}{
			"balance": hotel.Balance,
			"offers": bought,
		}
	}

	players = make([]map[string]interface{}, 0, len(g.players))
	for _, player := range g.players {
		playerStatus := map[string]interface{}{
			"name": player.name,
			"playerIndex": player.id,
			"balance": g.hotels[player.id].Balance,
		}
		players = append(players, playerStatus)
	}

	return
}

func MakeGame() *Game {
	return &Game{
		players: make(map[int]*Player),
		hotels: make(map[int]*model.Hotel),
		platforms: map[string]*model.Platform{
			"ta": model.MakePlatform("ta"),
			"book": model.MakePlatform("book"),
			"hc": model.MakePlatform("hc"),
			"hrs": model.MakePlatform("hrs"),
		},
		guests: make([]*model.Guest, 0, 1024),
	}
}
