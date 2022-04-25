package game

import (
	"fmt"
	"sync/atomic"
)

const (
	Ready = iota
	Play
	Fire
	GameOver
)

type GameMessage struct {
	Type int    `json:"type"`
	Msg  string `json:"msg,omitempty"`
}

func CreateMessage(msgType int) GameMessage {
	return GameMessage{
		Type: msgType,
		Msg:  "",
	}
}

func CreateWinnerMessage(gameStats *GameStats) GameMessage {
	return GameMessage{
		Type: GameOver,
		Msg:  fmt.Sprintf("winner(%v)___%v", atomic.LoadInt64(&ActiveGames), gameStats),
	}
}

func CreateLoserMessage() GameMessage {
	return GameMessage{
		Type: GameOver,
		Msg:  "Loser",
	}
}

func ErrorGameOver(msg string) GameMessage {
	return GameMessage{
		Type: GameOver,
		Msg:  msg,
	}
}
