package main

import (
	"flag"
	"log"
	"net/http"
	"time"

	"github.com/gabrielopesantos/carracing/pkg/game"
	"github.com/gorilla/websocket"
)

var addr = flag.String("addr", "127.0.0.1:8888", "Server address")

// Upgrades an HTTP connection to a WebSocket connection.
var upgrader = websocket.Upgrader{
	HandshakeTimeout: 3 * time.Second,
}

func main() {
	flag.Parse()
	conns := make(chan *websocket.Conn, 10)

	go func() {
		game.RunGame(conns)
	}()

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		c, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Println("Failed to upgrade connection!")
			return
		}
		conns <- c
	})

	log.Printf("Server listening on address %s", *addr)
	log.Fatal(http.ListenAndServe(*addr, nil))
}
