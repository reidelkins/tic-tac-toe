package db

import (
	"database/sql"
	"encoding/json"
	"fmt"

	_ "github.com/lib/pq"

	"github.com/reidelkins/kube-tic-tac-toe/internal/game"
)

type DB struct {
    *sql.DB
}

func NewDB(dataSourceName string) (*DB, error) {
    db, err := sql.Open("postgres", dataSourceName)
    if err != nil {
        return nil, err
    }
    if err = db.Ping(); err != nil {
        return nil, err
    }
    return &DB{db}, nil
}

// CreateGetPlayer checks if a player exists in the users table by username; if not, it creates a new user
func (db *DB) CreateGetPlayer(username string) (int64, error) {
    var id int64

    // First, try to find the user by username
    err := db.QueryRow(`SELECT id FROM users WHERE username = $1`, username).Scan(&id)
    if err != nil {
        if err == sql.ErrNoRows {
            // The user does not exist, so create a new user
            err = db.QueryRow(`INSERT INTO users (username) VALUES ($1) RETURNING id`, username).Scan(&id)
            if err != nil {
                return 0, fmt.Errorf("error creating new user: %v", err)
            }
            // New user created, return the new user ID
            return id, nil
        }
        // Some other error occurred
        return 0, fmt.Errorf("error checking user existence: %v", err)
    }

    // Existing user found, return the existing user ID
    return id, nil
}

func (db *DB) CreateGame(g *game.Game) (int64, error) {
    var gameId int64    
    serializedState, err := g.Serialize()    
    if err != nil {
        return 0, err
    }    

    err = db.QueryRow(
        `INSERT INTO games (player1_id, state, over) VALUES ($1, $2, $3) RETURNING id`,
        g.Player1ID,        
        serializedState,
        g.Over,
    ).Scan(&gameId)
    
    if err != nil {
        fmt.Println("Error: ", err)
        return 0, err
    }

    g.ID = gameId
    fmt.Println("Game ID: ", gameId)
    db.UpdateGame(g)

    return gameId, nil
}

func (db *DB) GetGame(gameId int64) (*game.Game, error) {
    var g game.Game
    var serializedState string
    var player2ID sql.NullInt64

    err := db.QueryRow(`SELECT id, player1_id, player2_id, state, over FROM games WHERE id = $1`, gameId).Scan(
        &g.ID,
        &g.Player1ID,
        &player2ID, // Scan into sql.NullInt64
        &serializedState, // Deserialize this into the Game struct
        &g.Over,
    )
    
    if err != nil {        
        return nil, err
    }

    // Check if player2ID is valid (not NULL), then assign it to g.Player2ID
    if player2ID.Valid {
        g.Player2ID = player2ID.Int64
    } else {
        // handle the case where player2ID is NULL if necessary
        g.Player2ID = 0 // or any other default value or logic you wish to apply
    }        

    // Deserialize the state into the Game struct
    if err := json.Unmarshal([]byte(serializedState), &g); err != nil {        
        return nil, err
    }
    
    return &g, nil
}


func (db *DB) UpdateGame(g *game.Game) error {
    serializedState, err := g.Serialize()  
    if err != nil {
        return err
    }
    fmt.Println("Serialized State: ", serializedState)
    
    // Using sql.NullInt64 to handle nullable player2_id correctly
    player2IDValue := sql.NullInt64{Int64: g.Player2ID, Valid: g.Player2ID != 0}
    fmt.Println("Player2ID Value: ", player2IDValue)
    _, err = db.Exec(
        `UPDATE games SET player1_id = $1, player2_id = $2, state = $3, over = $4 WHERE id = $5`,
        g.Player1ID,
        player2IDValue, // Using NullInt64 for player2_id
        serializedState,
        g.Over,
        g.ID,
    )
    if err != nil {
        fmt.Println("Error: ", err)
        return err
    }

    return nil
}

func (db *DB) ListGames() ([]game.Game, error) {
    rows, err := db.Query(`
        SELECT 
            g.id, 
            g.player1_id, 
            u1.username AS player1_username, 
            g.player2_id, 
            u2.username AS player2_username, 
            g.state, 
            g.over
        FROM 
            games g
        LEFT JOIN 
            users u1 ON g.player1_id = u1.id
        LEFT JOIN 
            users u2 ON g.player2_id = u2.id
    `)    
    if err != nil {
        return nil, err
    }
    
    defer rows.Close()
    
    var games []game.Game    
    
    for rows.Next() {         
        var g game.Game
        var serializedState string
        var player1Username, player2Username sql.NullString  // Use NullString to handle potential NULL values
        var player2ID sql.NullInt64
        if err := rows.Scan(&g.ID, &g.Player1ID, &player1Username, &player2ID, &player2Username, &serializedState, &g.Over); err != nil {            
            return nil, err
        }
        
        if player2ID.Valid {
            g.Player2ID = player2ID.Int64
        } else {
            g.Player2ID = 0  // Or however you want to handle a NULL player2_id
        }

        // Check if the username is valid (not NULL) before assigning
        g.Player1Username = player1Username.String  // The String method of NullString returns the string value
        g.Player2Username = player2Username.String

        // Deserialize the state into the Game struct
        if err := json.Unmarshal([]byte(serializedState), &g); err != nil {
            return nil, err
        }

        games = append(games, g)
    }    

    return games, nil
}