package game

const (
	Ready = iota
	Play
	Over
)

type Game struct {
	Distance             int
	CurrDistanceTraveled int
	GameMessages         chan GamePlayMessage
}

type GameStateMessage struct {
	State int
}

type GamePlayMessage struct {
	PlayerId byte
	Distance int
}
