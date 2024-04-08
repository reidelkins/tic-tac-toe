package api

import (
	"encoding/json"
	"fmt"
	"net/http"

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
	json.NewEncoder(w).Encode(currentGame)
}


// func PlayMoveHandler(w http.ResponseWriter, r *http.Request) {
// 	// Extracting a game ID from the URL path
// 	gameIDStr := r.URL.Query().Get("gameId")
// 	gameID, err := strconv.ParseInt(gameIDStr, 10, 64)
// 	if err != nil {
// 		http.Error(w, "Invalid game ID", http.StatusBadRequest)
// 		return
// 	}

// 	// Fetch the game from the database
// 	currentGame, err := dbConn.GetGame(gameID)
// 	if err != nil {
// 		http.Error(w, "Game not found", http.StatusNotFound)
// 		return
// 	}

// 	var move game.Move
// 	if err := json.NewDecoder(r.Body).Decode(&move); err != nil {
// 		http.Error(w, "Invalid move", http.StatusBadRequest)
// 		return
// 	}

// 	playerID := move.PlayerID // Assume Move struct includes PlayerID
// 	if err := currentGame.PlayMove(playerID, move.X, move.Y); err != nil {
// 		http.Error(w, "Failed to play move", http.StatusBadRequest)
// 		return
// 	}

// 	// Update the game in the database
// 	if err := dbConn.UpdateGame(currentGame); err != nil {
// 		http.Error(w, "Failed to update game", http.StatusInternalServerError)
// 		return
// 	}

// 	if currentGame.IsGameOver() {
// 		json.NewEncoder(w).Encode(map[string]interface{}{
// 			"result": currentGame.GetResult(),
// 			"state":  currentGame,
// 		})
// 		return
// 	}

// 	json.NewEncoder(w).Encode(currentGame)
// }
