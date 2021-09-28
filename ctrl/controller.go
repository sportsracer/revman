package ctrl

import (
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/sportsracer/revman/server"
)

const (
	buffer_size = 32
	tick_s = 1
)

type Controller struct {
	server server.Server
	*server.Subscriber

	players map[int]*Player
	game *Game
	mutex *sync.RWMutex
	ticker *time.Ticker

	state struct {
		players []map[string]interface{}
		offers []map[string]interface{}
	}
}

func (c *Controller) Run() {
	for {
		select {
		case id := <-c.Join:
			player := &Player{ id: id }
			c.players[id] = player
			log.Printf("Controller: Connected: %s", player)
		case msg := <-c.Rcv:
			player := c.players[msg.Id]
			switch msg.Data.Msg {
			case "join":
				if player.status != connected {
					log.Printf("Controller: Player cannot join, already joined! %s", player)
					continue
				}

				player.name = msg.Data.Data["name"].(string)
				player.status = in_game

				log.Printf("Controller: Joining game: %s", player)
				c.mutex.Lock()
				c.game.AddPlayer(player)
				c.mutex.Unlock()

				protMsg := map[string]interface{}{"playerIndex": player.id}
				c.server.Send(player.id, "joined", protMsg)
			case "input":
				if player.status != in_game {
					log.Printf("Controller: Player cannot send input, hasn't joined yet! %s", player)
					continue
				}

				c.mutex.Lock()
				offers := msg.Data.Data["offers"]
				for _, _offer := range offers.([]interface{}) {
					if offer, ok := _offer.(map[string]interface{}); ok {
						if _platform, ok := offer["platform"]; ok {
							if platform, ok := _platform.(string); ok {
								if _price, ok := offer["price"]; ok {
									if price, ok := _price.(float64); ok {
										c.game.AddOffer(player, platform, float32(price))
									}
								}
							}
						}
					}
				}
				c.mutex.Unlock()
			default:
				log.Printf("Controller: Unknown message type %s", msg.Data.Msg)
			}
		case <-c.ticker.C:
			c.mutex.Lock()
			states, players, offers := c.game.Tick()
			c.state.players = players
			c.state.offers = offers
			c.mutex.Unlock()
			for id, playerState := range states {
				c.server.Send(id, "state", playerState)
			}
		}
	}
}

func (c *Controller) getState() map[string]interface{} {
	state := make(map[string]interface{})
	state["players"] = c.state.players
	state["offers"] = c.state.offers
	return state
}

const (
	connected = iota
	in_game
)

type Player struct {
	id int
	name string
	status int
}

func (p *Player) String() string {
	return fmt.Sprintf("Player (id: %d, name: %s, status: %d)", p.id, p.name, p.status)
}

func MakeController(s server.Server) *Controller {
	sub := server.MakeSubscriber(buffer_size)
	game := MakeGame()
	ctrl := &Controller{
		server: s,
		Subscriber: sub,
		players: make(map[int]*Player),
		game: game,
		mutex: &sync.RWMutex{},
		ticker: time.NewTicker(tick_s * time.Second),
	}
	s.Subscribe(ctrl.Subscriber)
	return ctrl
}
