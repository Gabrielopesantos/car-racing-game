package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/url"

	"github.com/gabrielopesantos/carracing/pkg/game"
	"github.com/gorilla/websocket"
)

var host = flag.String("host", "127.0.0.1", "Server host")
var port = flag.String("port", "8888", "Server port")
var path = flag.String("path", "/", "URL path")

func readMessages(conn *websocket.Conn, messages chan game.GameMessage) {
	go func() {
		for {
			msgType, msg, err := conn.ReadMessage()
			log.Printf("Message read: %v, %v", msgType, msg)
			if err != nil {
				log.Println("Failed to read msg")
				return
			}

			if msgType != websocket.TextMessage {
				log.Print("?")
				continue
			}

			var gameMsg game.GameMessage
			err = json.Unmarshal(msg, &gameMsg)
			if err != nil {
				log.Print("Failed to unmarshal msg")
			}

			messages <- gameMsg
		}
	}()
}

func run(url url.URL) {
	messages := make(chan game.GameMessage)
	conn, _, err := websocket.DefaultDialer.Dial(url.String(), nil)
	if err != nil {
		log.Printf("Failed to connect to server")
		return
	}
	defer conn.Close()

	readMessages(conn, messages)

	go func() {
		for {
			msg := <-messages
			fmt.Println(msg)
		}
	}()

	for i := 0; i < 100000000; i++ {
		msg := game.CreateMessage(1)
		conn.WriteJSON(msg)
	}
}

func main() {
	flag.Parse()
	srvURL := url.URL{
		Scheme: "ws",
		Host:   *host + ":" + *port,
		Path:   *path,
	}

	run(srvURL)

}
