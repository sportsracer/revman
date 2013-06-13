package ctrl

import (
	"encoding/json"
	"log"
	"net/http"
)

func (c *Controller) MakeStateHandler() func(http.ResponseWriter, *http.Request) {

	handle := func(w http.ResponseWriter, r *http.Request) {
		state := c.getState()
		stateJson, err := json.Marshal(state)
		if err != nil {
			log.Printf("Error encoding state!")
		}
		w.Write([]byte(stateJson))
	}
	return handle
}
