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
	"github.com/reidelkins/kube-tic-tac-toe/internal/utils"
)

type Handler struct {
    DBConn *db.DB
}

func (h *Handler) CreateGameHandler(w http.ResponseWriter, r *http.Request) {    
    tokenCookie, err := r.Cookie("csrftoken")    
    if err != nil {
        // Handle the case when the cookie is not found
        http.Error(w, "Token cookie not found", http.StatusUnauthorized)
        return
    }
        
    // Verify the Google ID token
    payload, err := utils.VerifyGoogleIDToken(r.Context(), tokenCookie.Value)
    if err != nil {        
        http.Error(w, "Invalid token", http.StatusUnauthorized)
        return
    }
    fmt.Println(payload)
    // Extract relevant information from the token payload
    email := payload.Claims["email"].(string)    
    
    	
    // Check if the database connection is nil
	if h.DBConn == nil {
		fmt.Println("DB Conn is nil")
		return
	}	

	playerID, err := h.DBConn.CreateGetPlayer(email)
	if err != nil {
		http.Error(w, "Failed to get or create player", http.StatusInternalServerError)
		return
	}
    
	// Initialize a new game with the player ID
    newGame := game.NewGame(playerID, email)	
	

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
            "http://localhost",
            "https://c15f-70-116-143-2.ngrok-free.app",
             
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
    // Extract the token from the request (e.g., from the query parameters or headers)    
    tokenCookie, err := r.Cookie("token")
    if err != nil {
        // Handle the case when the cookie is not found
        http.Error(w, "Token cookie not found", http.StatusUnauthorized)
        return
    }

    // Verify the Google ID token
    _, err = utils.VerifyGoogleIDToken(r.Context(), tokenCookie.Value)
    if err != nil {
        http.Error(w, "Invalid token", http.StatusUnauthorized)
        return
    }


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

func (h *Handler) GoogleLoginHandler(w http.ResponseWriter, r *http.Request) {
    // Example: Extracting the Google ID token from the request, adjust as necessary
    var token struct {
        Token string `json:"token"`
    }

    if err := json.NewDecoder(r.Body).Decode(&token); err != nil {
        http.Error(w, "Invalid request", http.StatusBadRequest)
        return
    }    
    
    // Verify the Google ID token
    payload, err := utils.VerifyGoogleIDToken(r.Context(), token.Token)
    if err != nil {
        http.Error(w, "Invalid token", http.StatusUnauthorized)
        return
    }

    // Extract relevant information from the token payload    
    email := payload.Claims["email"].(string)

    // Example: Return the player ID or any other relevant data to the client
    response := map[string]string{        
        "email":  email,
        "token": token.Token,
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(response)     
}