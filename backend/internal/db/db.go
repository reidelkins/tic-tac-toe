package db

import (
	"database/sql"
	"encoding/json"

	_ "github.com/lib/pq"

	"github.com/reidelkins/kube-tic-tac-toe/internal/game"
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
    serializedState, err := g.Serialize()
    if err != nil {
        return 0, err
    }

    err = db.QueryRow(
        `INSERT INTO games (player1_id, player2_id, state, over) VALUES ($1, $2, $3, $4) RETURNING id`,
        g.Player1ID,
        g.Player2ID,
        serializedState,
        g.Over,
    ).Scan(&gameId)
    if err != nil {
        return 0, err
    }

    return gameId, nil
}

func (db *DB) GetGame(gameId int64) (*game.Game, error) {
    var g game.Game
    var serializedState string

    err := db.QueryRow(`SELECT id, player1_id, player2_id, state, over FROM games WHERE id = $1`, gameId).Scan(
        &g.ID,
        &g.Player1ID,
        &g.Player2ID,
        &serializedState, // Deserialize this into the Game struct
        &g.Over,
    )
    if err != nil {
        return nil, err
    }

    // Deserialize the state into the Game struct
    if err := json.Unmarshal([]byte(serializedState), &g); err != nil {
        return nil, err
    }

    return &g, nil
}


func (db *DB) UpdateGame(g *game.Game) error {
    serializedState, err := json.Marshal(g)
    if err != nil {
        return err
    }

    _, err = db.Exec(
        `UPDATE games SET player1_id = $1, player2_id = $2, state = $3, over = $4 WHERE id = $5`,
        g.Player1ID,
        g.Player2ID,
        serializedState, // Now a serialized JSON string
        g.Over,
        g.ID,
    )
    return err
}

func (db *DB) ListGames() ([]game.Game, error) {
    rows, err := db.Query(`SELECT id, player1_id, player2_id, state, over FROM games`)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var games []game.Game
    for rows.Next() {
        var g game.Game
        var serializedState string
        if err := rows.Scan(&g.ID, &g.Player1ID, &g.Player2ID, &serializedState, &g.Over); err != nil {
            return nil, err
        }

        // Deserialize the state into the Game struct
        if err := json.Unmarshal([]byte(serializedState), &g); err != nil {
            return nil, err
        }

        games = append(games, g)
    }

    return games, nil
}