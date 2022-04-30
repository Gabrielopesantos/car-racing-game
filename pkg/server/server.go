package server

import (
	"fmt"
	"log"
	"strconv"
	"sync"
	"time"

	"github.com/gabrielopesantos/carracing/pkg/game"
	"github.com/gorilla/websocket"
)

func waitForReady(conn *websocket.Conn, wg *sync.WaitGroup) {
	defer func() {
		defer wg.Done()
	}()

	err := conn.WriteMessage(websocket.TextMessage, []byte("ready"))
	if err != nil {
		log.Fatal("Failed to send `ready` message")
	}

	_, msg, err := conn.ReadMessage()
	if err != nil {
		log.Fatal("Failed to recieve `ready` message")
	}

	if string(msg) != "ready" {
		log.Fatal("Invalid message recieved")
	}
}

func Run(conns <-chan *websocket.Conn, wg *sync.WaitGroup) {
	player1 := <-conns
	player2 := <-conns
	players := []*websocket.Conn{player1, player2}

	go func() {
		// Ready players
		wg := &sync.WaitGroup{}
		for _, p := range players {
			wg.Add(1)
			waitForReady(p, wg)
		}
		wg.Wait()

		// Create game and start
		game := game.Game{
			Distance:     1000,
			GameMessages: make(chan game.GameMessage),
		}
		for _, p := range players {
			wg.Add(1)
			p.WriteMessage(websocket.TextMessage, []byte("play"))
		}

		listenForMessages(player1, game.GameMessages, '1')
		listenForMessages(player2, game.GameMessages, '2')

		// podium := map[string]*websocket.Conn{"winner": nil, "loser": nil}
		ticker := time.NewTicker(5 * time.Second)

	gameMainLoop:
		for {
			select {
			case gMsg := <-game.GameMessages:
				fmt.Println(gMsg)
			case <-ticker.C:
				fmt.Println("ticker")
				break gameMainLoop
			}
		}
	}()
}

func listenForMessages(pConn *websocket.Conn, gameMessages chan<- game.GameMessage, playerId byte) {
	go func() {
		_, msg, err := pConn.ReadMessage()
		if err != nil {
			log.Fatalf("Failed to read message")
		}

		d, _ := strconv.Atoi(string(msg))
		gMsg := game.GameMessage{
			PlayerId: playerId,
			Distance: d,
		}
		gameMessages <- gMsg
	}()
}
