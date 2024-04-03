package db

import (
	"database/sql"
	"tic-tac-toe/backend/internal/game"

	_ "github.com/lib/pq"
)

type DB struct {
    *sql.DB
}

func NewDB(dataSourceName string) (*sql.DB, error) {
    db, err := sql.Open("postgres", dataSourceName)
    if err != nil {
        return nil, err
    }
    if err = db.Ping(); err != nil {
        return nil, err
    }
    return db, nil
}

func (db *DB) CreateGame(g *game.Game) (int64, error) {
    var gameId int64
    err := db.QueryRow(
        `INSERT INTO games (player1_id, state, over) VALUES ($1, $2, $3) RETURNING id`,
        g.Player1ID, // Assuming Player1ID is set when creating the game
        g.State,     // State could be a JSON field or a serialized version of the game state
        g.Over,
    ).Scan(&gameId)
    if err != nil {
        return 0, err
    }
    return gameId, nil
}

func (db *DB) GetGame(gameId int64) (*game.Game, error) {
    var g game.Game
    err := db.QueryRow(`SELECT id, player1_id, player2_id, state, over FROM games WHERE id = $1`, gameId).Scan(
        &g.ID,
        &g.Player1ID,
        &g.Player2ID,
        &g.State, // Make sure to deserialize the state into the Game struct
        &g.Over,
    )
    if err != nil {
        return nil, err
    }
    return &g, nil
}

func (db *DB) UpdateGame(g *game.Game) error {
    _, err := db.Exec(
        `UPDATE games SET player1_id = $1, player2_id = $2, state = $3, over = $4 WHERE id = $5`,
        g.Player1ID,
        g.Player2ID,
        g.State, // State should be serialized or converted to a JSON string if using a JSON field
        g.Over,
        g.ID,
    )
    return err
}