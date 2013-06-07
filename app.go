package main

import (
	"flag"
	"log"
	"net/http"
)

import (
	"me/revman/ctrl"
	"me/revman/server"
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
	if err := http.ListenAndServe(*addr, nil); err != nil {
		log.Fatal("Server startup", err)
	}
}
