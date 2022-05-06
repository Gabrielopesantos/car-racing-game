package main

import (
	"flag"
	"log"
	"math/rand"
	"net/url"
	"sync"
	"time"

	"github.com/gabrielopesantos/carracing/pkg/game"
	"github.com/gorilla/websocket"
)

var host = flag.String("host", "127.0.0.1", "Server host")
var port = flag.String("port", "8888", "Server port")
var path = flag.String("path", "/", "URL path")

var srvURL = url.URL{
	Scheme: "ws",
	Host:   *host + ":" + *port,
	Path:   *path,
}

func readMessages(conn *websocket.Conn, messages chan<- game.GameStateMessage) {
	go func() {
		for {
			var gameMsg game.GameStateMessage
			err := conn.ReadJSON(&gameMsg)
			if err != nil {
				log.Printf("Failed to read msg. Error: %v", err)
				continue
			}
			messages <- gameMsg
			if gameMsg.State == game.Over {
				return
			}
		}
	}()
}

func connect(wg *sync.WaitGroup) {
	defer func() {
		log.Print("Entering done")
		wg.Done()
	}()
	messages := make(chan game.GameStateMessage)

	conn, _, err := websocket.DefaultDialer.Dial(srvURL.String(), nil)
	if err != nil {
		log.Printf("Failed to connect to server")
		return
	}
	defer conn.Close()

	readMessages(conn, messages)

sendMessagesLoop:
	for {
		msg := <-messages
		switch msg.State {
		case game.Ready:
			log.Print("Entering in `game.Ready`")
			rMsg := game.GameStateMessage{State: game.Ready}
			conn.WriteJSON(rMsg)
		case game.Play:
			go func() {
				for {
					pMsg := game.GamePlayMessage{Distance: rand.Intn(500)}
					_ = conn.WriteJSON(pMsg)
					time.Sleep(300 * time.Millisecond)
				}
			}()
		case game.Over:
			log.Print("Game over")
			conn.Close()
			break sendMessagesLoop
		default:
			continue
		}
	}

	close(messages)
}

func run() {
	wg := sync.WaitGroup{}

	for i := 0; i < 2; i++ {
		wg.Add(1)
		go func() {
			connect(&wg)
		}()
	}

	wg.Wait()
}

func main() {
	flag.Parse() // Is this still necessary?

	run()
}
