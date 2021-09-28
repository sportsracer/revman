package main

import (
	"flag"
	"log"
	"net/http"

	"github.com/sportsracer/revman/ctrl"
	"github.com/sportsracer/revman/server"
)

func main() {
	var addr = flag.String("addr", ":8090", "Port")
	flag.Parse()
	var server = server.MakeWsServer()
	var ctrl = ctrl.MakeController(server)

	go server.Run()
	go ctrl.Run()

	log.Printf("RevMan listening on port %s", *addr)
	http.Handle("/ws", server.MakeWsHandler())
	http.Handle("/", http.FileServer(http.Dir("static")))
	http.HandleFunc("/state", ctrl.MakeStateHandler())
	if err := http.ListenAndServe(*addr, nil); err != nil {
		log.Fatal("Server startup", err)
	}
}
