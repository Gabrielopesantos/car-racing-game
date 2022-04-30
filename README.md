# Game

Small project to get familiar with golang threading primitives and websockets. (WIP)

## Game state steps
- When two players connect
- Server sends message `Ready`
- After both send `Ready` to server game starts
- Players continuously send distance travelled
- When a player hits track distance, server sends `GameOver`