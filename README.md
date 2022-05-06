# "Game"

Small project to get familiar with golang concurrency primitives and websockets.

## Idea
As the game suggests, each "game" is a car race with a defined distance. For a race to start, two players have to connect to a game, after that `ready` messages and exchanged between server and clients. 

If this exchange succeeds, players continuously send the traveled distance (a random value) until the one of them reaches the defined distance.

When reached, servers sends custom message to winner and loser and closes connection.

## Game state steps
- server waits for two players to connect
- `ready` message is exchanged between server and clients
- server sends `play`
- client send distance traveled (random number)
- first cient to reach game distance wins
- server sends `over` message with winner or loser

## How to run
- start by running the server (eg. `go run ./cmd/server/main.go`)
- run client (eg. `go run ./cmd/client/main.go`)

### License

MIT

