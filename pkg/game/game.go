package game

const (
	Ready = iota
	Play
	Over
)

type Game struct {
	Distance             int
	CurrDistanceTraveled int
	GameMessages         chan GameMessage
}

type GameMessage struct {
	PlayerId byte
	Distance int
}
