package server

import (
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/gabrielopesantos/carracing/pkg/game"
	"github.com/gorilla/websocket"
)

func waitForReady(conn *websocket.Conn, wg *sync.WaitGroup) {
	defer func() {
		defer wg.Done()
	}()

	msg := game.StateMessage{State: game.Ready}
	err := conn.WriteJSON(msg)
	if err != nil {
		log.Fatal("Failed to send `ready` message")
	}

	err = conn.ReadJSON(&msg)
	if err != nil {
		log.Fatal("Failed to receive `ready` message")
	}

	if msg.State != game.Ready {
		log.Fatal("Invalid message received")
	}
}

func listenForMessages(p *game.Player, gameMessages chan<- game.PlayMessage) {
	go func() {
		for {
			gMsg := game.PlayMessage{}
			err := p.Conn.ReadJSON(&gMsg)
			if err != nil {
				return
			}

			gMsg.PlayerId = p.Identifier
			gameMessages <- gMsg
		}
	}()
}

func Run(conns <-chan *websocket.Conn) {
	for {
		p1_conn := <-conns
		p2_conn := <-conns

		go func() {
			// Create game
			gameInstance := game.NewGame(1000)
			gameInstance.Players[0].Conn = p1_conn
			gameInstance.Players[1].Conn = p2_conn

			// Ready players
			wg := &sync.WaitGroup{}
			for _, p := range gameInstance.Players {
				wg.Add(1)
				waitForReady(p.Conn, wg)
			}
			wg.Wait()

			// Setup and start game
			playMsg := game.StateMessage{State: game.Play}
			for _, p := range gameInstance.Players {
				err := p.Conn.WriteJSON(playMsg)
				if err != nil {
					log.Fatalf("Failed to send `play` message")
				}
			}

			listenForMessages(&gameInstance.Players[0], gameInstance.GameMessages)
			listenForMessages(&gameInstance.Players[1], gameInstance.GameMessages)

			// podium := map[string]*websocket.Conn{"winner": nil, "loser": nil}
			ticker := time.NewTicker(30 * time.Second)

		gameMainLoop:
			for {
				select {
				case gMsg := <-gameInstance.GameMessages:
					log.Printf("P1\n%+v", gameInstance.Players)
					log.Printf("P2\n%+v", gameInstance.Players)
					for playerIndex := range gameInstance.Players {
						if gMsg.PlayerId == gameInstance.Players[playerIndex].Identifier {
							gameInstance.Players[playerIndex].DistanceTraveled += gMsg.Distance
						}
						if gameInstance.Players[playerIndex].DistanceTraveled >= gameInstance.Distance {
							break gameMainLoop
						}
					}

				case <-ticker.C:
					break gameMainLoop
				}
			}

			overMsg := game.StateMessage{State: game.Over}
			for _, p := range gameInstance.Players {
				// wg.Add(1)
				err := p.Conn.WriteJSON(overMsg)
				if err != nil {
					log.Fatalf("Failed to send `over` message")
				}
				p.Conn.Close()
			}
			fmt.Println("Done")
		}()
	}

}
