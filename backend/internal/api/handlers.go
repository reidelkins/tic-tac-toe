package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/gorilla/websocket"
	log "github.com/sirupsen/logrus"

	"github.com/reidelkins/kube-tic-tac-toe/internal/db"
	"github.com/reidelkins/kube-tic-tac-toe/internal/game"
)

type Handler struct {
    DBConn *db.DB
}

func (h *Handler) CreateGameHandler(w http.ResponseWriter, r *http.Request) {
    // Example of extracting player ID from the request, adjust as necessary
    var playerInfo struct {
        Player1Username string `json:"player1Username"`
    }	
    if err := json.NewDecoder(r.Body).Decode(&playerInfo); err != nil {		
        http.Error(w, "Invalid request", http.StatusBadRequest)
        return
    }

    	
    // Check if the database connection is nil
	if h.DBConn == nil {
		fmt.Println("DB Conn is nil")
		return
	}	

	playerID, err := h.DBConn.CreateGetPlayer(playerInfo.Player1Username)
	if err != nil {
		http.Error(w, "Failed to get or create player", http.StatusInternalServerError)
		return
	}

	// Initialize a new game with the player ID
    newGame := game.NewGame(playerID, playerInfo.Player1Username)	
	

    gameID, err := h.DBConn.CreateGame(newGame)
	
    if err != nil {		
        http.Error(w, "Failed to create game", http.StatusInternalServerError)
        return
    }
	    
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(gameID)
}

func (h *Handler) ListGamesHandler(w http.ResponseWriter, r *http.Request) {	
	// Fetch all active games from the database
	games, err := h.DBConn.ListGames()
	if err != nil {
		http.Error(w, "Failed to list games", http.StatusInternalServerError)
		return
	}	

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(games); err != nil {
        http.Error(w, "Error encoding games", http.StatusInternalServerError)
        return
    }
}

func (h *Handler) GetGameHandler(w http.ResponseWriter, r *http.Request) {
    // Extracting a game ID from the URL path    
    gameIDStr := chi.URLParam(r, "gameId")
    gameID, err := strconv.ParseInt(gameIDStr, 10, 64)
    if err != nil {
        http.Error(w, "Invalid game ID", http.StatusBadRequest)
        return
    }

    // Fetch the game from the database
    currentGame, err := h.DBConn.GetGame(gameID)
    if err != nil {
        http.Error(w, "Game not found", http.StatusNotFound)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(currentGame)
}

func (h *Handler) JoinGameHandler(w http.ResponseWriter, r *http.Request) {
	// Extracting player ID from the request, adjust as necessary
	var playerInfo struct {
		Player2Username string `json:"player2Username"`
		GameID          int64  `json:"gameId"`
	}
	if err := json.NewDecoder(r.Body).Decode(&playerInfo); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	// Check if the database connection is nil
	if h.DBConn == nil {
		fmt.Println("DB Conn is nil")
		return
	}

	playerID, err := h.DBConn.CreateGetPlayer(playerInfo.Player2Username)
	if err != nil {
		http.Error(w, "Failed to get or create player", http.StatusInternalServerError)
		return
	}

	// Fetch the game from the database
	currentGame, err := h.DBConn.GetGame(playerInfo.GameID)
	if err != nil {
		http.Error(w, "Game not found", http.StatusNotFound)
		return
	}

	// Check if the game already has two players
	if currentGame.Player2ID != 0 {
		http.Error(w, "Game already has two players", http.StatusBadRequest)
		return
	}
	
	// Check if new player is trying to join their own game
	if currentGame.Player1ID == playerID {
		http.Error(w, "Player cannot join their own game", http.StatusBadRequest)
		return
	}

	// Assign the player ID to the game
	currentGame.Player2ID = playerID
	currentGame.Player2Username = playerInfo.Player2Username

	// Update the game in the database
	if err := h.DBConn.UpdateGame(currentGame); err != nil {
		http.Error(w, "Failed to update game", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(currentGame.ID)
}

var connections = make(map[int64][]*websocket.Conn)

var upgrader = websocket.Upgrader{
    ReadBufferSize:  1024,
    WriteBufferSize: 1024,
    CheckOrigin: func(r *http.Request) bool {
        // Get the origin of the request
        origin := r.Header.Get("Origin")
		
        // Check if the origin is allowed
        allowedOrigins := []string{
            "http://localhost:4200",  // Example: Allow connections from localhost:4200
			"http://127.0.0.1:4200",
            "https://your-angular-app-url.com",  // Example: Allow connections from your Angular app URL
        }

        for _, allowedOrigin := range allowedOrigins {
            if origin == allowedOrigin {
                return true
            }
        }

        // If the origin is not allowed, return false
        return false
    },
}

func (h *Handler) WebSocketHandler(w http.ResponseWriter, r *http.Request) {    
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
        currentGame, err := h.DBConn.GetGame(gameID)
        if err != nil {
            log.Println("Get Game:", err)
            continue
        }        
        if currentGame.PlayMove(msg.PlayerID, msg.X, msg.Y) {
            // Save the updated game state to the database
            err = h.DBConn.UpdateGame(currentGame)
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
        } else {
            log.Println("Invalid move")            
        }
    }
}