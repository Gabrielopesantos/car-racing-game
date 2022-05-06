package game

import "github.com/gorilla/websocket"

const (
	Ready = iota
	Play
	Over
)

type Player struct {
	Identifier       string
	Conn             *websocket.Conn
	DistanceTraveled int
}

type Game struct {
	Distance     int
	Players      []Player
	GameMessages chan PlayMessage
}

func NewGame(distance int) *Game {
	return &Game{
		Distance: distance,
		Players: []Player{
			{
				Identifier:       "Player 1",
				DistanceTraveled: 0,
			},
			{
				Identifier:       "Player 2",
				DistanceTraveled: 0,
			},
		},
		GameMessages: make(chan PlayMessage),
	}

}

type StateMessage struct {
	State int
}

type PlayMessage struct {
	PlayerId string
	Distance int
}

func CreateStateMessage(state int) *StateMessage {
	return &StateMessage{
		State: state,
	}
}
