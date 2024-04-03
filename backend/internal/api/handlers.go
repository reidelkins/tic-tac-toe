package api

import (
	"encoding/json"
	"net/http"
	"tic-tac-toe/backend/internal/db"
	"tic-tac-toe/backend/internal/game"
)

// dbConn represents the database connection
var dbConn *db.DB

func CreateGameHandler(w http.ResponseWriter, r *http.Request) {
	// Initialize a new game
	newGame := game.NewGame()

	// Save the new game to the database
	gameID, err := dbConn.CreateGame(newGame)
	if err != nil {
		http.Error(w, "Failed to create game", http.StatusInternalServerError)
		return
	}

	// Fetch the newly created game from the database
	savedGame, err := dbConn.GetGame(gameID)
	if err != nil {
		http.Error(w, "Failed to retrieve game", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(savedGame)
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
