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

	msg := game.GameStateMessage{State: game.Ready}
	err := conn.WriteJSON(msg)
	if err != nil {
		log.Fatal("Failed to send `ready` message")
	}

	err = conn.ReadJSON(&msg)
	if err != nil {
		log.Fatal("Failed to recieve `ready` message")
	}

	if msg.State != game.Ready {
		log.Fatal("Invalid message recieved")
	}
}

func listenForMessages(pConn *websocket.Conn, gameMessages chan<- game.GamePlayMessage, playerId byte) {
	go func() {
		for {
			_, msg, err := pConn.ReadMessage()
			if err != nil {
				log.Fatalf("Failed to read message")
			}

			d, _ := strconv.Atoi(string(msg))
			gMsg := game.GamePlayMessage{
				PlayerId: playerId,
				Distance: d,
			}
			gameMessages <- gMsg
		}
	}()
}

func Run(conns <-chan *websocket.Conn) {
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
		g := game.Game{
			Distance:     1000,
			GameMessages: make(chan game.GamePlayMessage),
		}
		playMsg := game.GameStateMessage{State: game.Play}
		for _, p := range players {
			wg.Add(1)
			err := p.WriteJSON(playMsg)
			if err != nil {
				log.Fatalf("Failed to send `play` message")
			}
		}

		listenForMessages(player1, g.GameMessages, '1')
		listenForMessages(player2, g.GameMessages, '2')

		// podium := map[string]*websocket.Conn{"winner": nil, "loser": nil}
		ticker := time.NewTicker(5 * time.Second)

	gameMainLoop:
		for {
			select {
			case gMsg := <-g.GameMessages:
				fmt.Printf("%+v\n", gMsg)
			case <-ticker.C:
				fmt.Println("ticker")
				break gameMainLoop
			}
		}
	}()
}
