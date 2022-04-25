package game

import (
	"encoding/json"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

type Game struct {
	Players [2]*Player
	Bullets []*Bullet
}

type Debug struct {
	Time  int64
	Place string
}

func readyUp(wg *sync.WaitGroup, conn *websocket.Conn) {
	go func() {
		// Why defer a closure and not simply `wg.Done()`?
		defer func() {
			wg.Done()
		}()

		conn.WriteJSON(CreateMessage(Ready))
		var resp GameMessage
		_, msg, err := conn.ReadMessage()

		if err != nil {
			log.Fatalf("Error reading message")
		}

		err = json.Unmarshal(msg, &resp)
		if err != nil {
			log.Fatalf("Error decoding message")
		}
	}()
}

type NamedGameMessage struct {
	msg  GameMessage
	name byte
}

func listenForFires(channel chan<- NamedGameMessage, conn *websocket.Conn, name byte) {
	go func() {
		for {
			var cmd GameMessage
			msgType, msg, err := conn.ReadMessage()
			if err != nil {
				log.Printf("Error reading `fire` message. Error: %v", err)
				return
			}

			if msgType != websocket.TextMessage {
				log.Print("?")
				continue
			}

			json.Unmarshal(msg, &cmd)
			if cmd.Type == Fire {
				channel <- NamedGameMessage{msg: cmd, name: name}
			} else {
				log.Fatalf("Message type recieved `%d`. Sshould have been `Fire`", cmd.Type)
			}
		}
	}()
}

func checkBulletCollisions(g *Game) {
bulletsLoop:
	for i := 0; i < len(g.Bullets); {
		bullet := g.Bullets[i]
		for j := i + 1; j < len(g.Bullets); j++ {
			bullet2 := g.Bullets[j]
			if bullet.Geo.HasCollision(&bullet2.Geo) {
				// that is also very crappy code.  Why would I ever do this...
				g.Bullets = append(g.Bullets[:j], g.Bullets[(j+1):]...)
				g.Bullets = append(g.Bullets[:i], g.Bullets[(i+1):]...)
				break bulletsLoop
			}
		}
		i += 1
	}
}

func RunGame(conns chan *websocket.Conn) {
	for {
		pA := <-conns
		pB := <-conns

		go func() {
			defer pA.Close()
			defer pB.Close()

			gameStartTime := time.Now()
			hack := []Debug{}

			// Wait for both player to be ready?
			wg := sync.WaitGroup{}
			wg.Add(2)
			readyUp(&wg, pA)
			readyUp(&wg, pB)
			wg.Wait()

			// Game state
			game := Game{
				Players: [2]*Player{
					NewPlayer(Vector2D{2500.0, 0.0}, Vector2D{-1.0, 0.0}, 180),
					NewPlayer(Vector2D{-2500.0, 0.0}, Vector2D{1.0, 0.0}, 300),
				},
				Bullets: []*Bullet{},
			}

			// Play game
			stats := NewGameStat()
			playMsg, err := json.Marshal(CreateMessage(Play))
			if err != nil {
				log.Fatalf("Failed to marshal `play` game message")
			}
			pA.WriteMessage(websocket.TextMessage, playMsg)
			pB.WriteMessage(websocket.TextMessage, playMsg)

			AddActiveGame()
			// Listen for fires
			fires := make(chan NamedGameMessage, 10)
			listenForFires(fires, pA, 'a')
			listenForFires(fires, pB, 'b')

			// Steps 5. The rust version has a tokio::select
			ticker := time.NewTicker(time.Millisecond * 16)
			last_start := time.Now()

			var winner *websocket.Conn
			var loser *websocket.Conn

		gameMainLoop:
			for {
				select {
				case fire := <-fires:
					player := game.Players[0]
					if fire.name == 'b' {
						player = game.Players[1]
					}

					if PlayerFire(player) {
						game.Bullets = append(game.Bullets, CreateBulletFromPlayer(player, 1.0))
						hack = append(hack, Debug{
							int64(time.Since(gameStartTime).Milliseconds()),
							fmt.Sprintf("fire%v", fire.name),
						})
					}

				case <-ticker.C:
					// 6. part 1 : calculate the time difference between each loop.
					diff := time.Since(last_start).Microseconds()
					last_start = time.Now()

					// 6. do all the collision / updating
					for i := 0; i < len(game.Bullets); i += 1 {
						UpdateBullet(game.Bullets[i], diff)
					}

					checkBulletCollisions(&game)

					for i := 0; i < len(game.Bullets); i += 1 {
						if game.Players[0].Geo.HasCollision(&game.Bullets[i].Geo) {
							winner = pA
							loser = pB
							break gameMainLoop
						}
						if game.Players[1].Geo.HasCollision(&game.Bullets[i].Geo) {
							winner = pA
							loser = pB
							break gameMainLoop
						}
					}

					stats.AddDelta(diff)
				}
			}

			// Part 7. Send out the winner / loser message and close down the
			// suckets
			winnerMsg := CreateWinnerMessage(stats)
			loserMsg := CreateLoserMessage()

			winner.WriteJSON(winnerMsg)
			loser.WriteJSON(loserMsg)

			if stats.FrameBuckets[0] > 1000 {
				log.Printf("COLLISIONS\n")
				for _, debug := range hack {
					log.Printf("%+v\n", debug)
				}
				log.Println()
			}

			RemoveActiveGame()
		}()
	}
}
