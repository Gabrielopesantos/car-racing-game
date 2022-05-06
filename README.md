# "Game"

Small project to get familiar with golang concurrency primitives and websockets.

## Idea
As the game suggests, each "game" is a car race with a defined distance. For a race to start, two players have to connect to a game, after that `ready` messages and exchanged between server and clients. 

If this exchange succeeds, players continuously send the traveled distance (a random value) until the one of them reaches the defined distance.

When reached, servers sends custom message to winner and loser and closes connection.

## Game state steps (WIP)
- When two players connect
- Server sends message `Ready`
- After both send `Ready` to server game starts
- Players continuously send distance travelled
- When a player hits track distance, server sends `GameOver`

### License

MIT