package main

import (
	"flag"
	"log"
	"net/url"
	"sync"

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

func readMessages(conn *websocket.Conn, messages chan<- game.GameMessage) {
	go func() {
		for {
			var gameMsg game.GameMessage
			err := conn.ReadJSON(&gameMsg)
			if err != nil {
				log.Printf("Failed to read msg. Error: %v", err)
				return
			}

			messages <- gameMsg
		}
	}()
}

func connect(wg *sync.WaitGroup) {
	defer func() {
		log.Print("Entering done")
		wg.Done()
	}()
	messages := make(chan game.GameMessage)

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
		switch msg.Type {
		case game.Ready:
			log.Print("Enterning in `game.Ready`")
			conn.WriteJSON(game.CreateMessage(game.Ready))
		case game.Play:
			conn.WriteJSON(game.CreateMessage(game.Fire))
		case game.GameOver:
			break sendMessagesLoop
		default:
			continue
		}
	}

	close(messages)
}

func run() {
	wg := sync.WaitGroup{}

	for i := 0; i < 1; i++ {
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
