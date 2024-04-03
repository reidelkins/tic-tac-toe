package api

import (
	"net/http"
	"strconv"

	"github.com/gorilla/websocket"
	log "github.com/sirupsen/logrus"

	"github.com/reidelkins/kube-tic-tac-toe/internal/db"
	"github.com/reidelkins/kube-tic-tac-toe/internal/game"
)

var connections = make(map[int64][]*websocket.Conn)
var dbConn *db.DB

var upgrader = websocket.Upgrader{
    ReadBufferSize:  1024,
    WriteBufferSize: 1024,
    CheckOrigin: func(r *http.Request) bool {
        // Check the origin of the request (modify according to your needs)
        return true
    },
}

func WebSocketHandler(w http.ResponseWriter, r *http.Request) {
    conn, err := upgrader.Upgrade(w, r, nil)
    if err != nil {
        log.Println("Upgrade:", err)
        return
    }
    defer conn.Close()

    // Example: Assume we get the game ID from the request (adjust as needed)
    gameID, err := strconv.ParseInt(r.URL.Query().Get("gameId"), 10, 64)
    if err != nil {
        log.Println("Invalid Game ID:", err)
        return
    }

    // Add the connection to the list for the game
    connections[gameID] = append(connections[gameID], conn)

    for {
        var msg game.Move
        err := conn.ReadJSON(&msg)
        if err != nil {
            log.Println("Read:", err)
            break
        }

        // Process the move (e.g., update the game state)
        currentGame, err := dbConn.GetGame(gameID)
        if err != nil {
            log.Println("Get Game:", err)
            continue
        }

        if currentGame.PlayMove(msg.PlayerID, msg.X, msg.Y) {
            // Save the updated game state to the database
            err = dbConn.UpdateGame(currentGame)
            if err != nil {
                log.Println("Update Game:", err)
                continue
            }

            // Broadcast the updated game state to all connected clients for this game
            for _, c := range connections[gameID] {
                err = c.WriteJSON(currentGame)
                if err != nil {
                    log.Println("Write:", err)
                    continue
                }
            }
        }
    }
}
