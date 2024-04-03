package game

import (
    "time"
    "encoding/json"
)

type Game struct {
    ID        int64
    Player1ID int64
    Player2ID int64
    Board     [3][3]string
    CurrentPlayer string
    Winner    string
    Over      bool
    CreatedAt time.Time
    UpdatedAt time.Time
}

type Move struct {
    PlayerID int64 `json:"playerId"`
    X        int   `json:"x"`
    Y        int   `json:"y"`
}

func NewGame(player1ID int64) *Game {
    return &Game{
        Player1ID: player1ID,
        Board:     [3][3]string{},
        CurrentPlayer: "X",
        CreatedAt: time.Now(),
        UpdatedAt: time.Now(),
    }
}

func (g *Game) PlayMove(playerID int64, x, y int) bool {
    // Check if the game is over or the cell is already occupied
    if g.Board[x][y] != "" || g.Over {
        return false
    }

    // Check if it's the correct player's turn
    if (g.CurrentPlayer == "X" && playerID != g.Player1ID) || (g.CurrentPlayer == "O" && playerID != g.Player2ID) {
        return false
    }

    // Make the move
    g.Board[x][y] = g.CurrentPlayer
    g.checkWinner()
    g.switchPlayer()
    g.UpdatedAt = time.Now()
    return true
}

func (g *Game) switchPlayer() {
    if g.CurrentPlayer == "X" {
        g.CurrentPlayer = "O"
    } else {
        g.CurrentPlayer = "X"
    }
}

func (g *Game) checkWinner() {
    // Check horizontal and vertical lines
    for i := 0; i < 3; i++ {
        if g.Board[i][0] != "" && g.Board[i][0] == g.Board[i][1] && g.Board[i][0] == g.Board[i][2] {
            g.Winner = g.Board[i][0]
            g.Over = true
            return
        }
        if g.Board[0][i] != "" && g.Board[0][i] == g.Board[1][i] && g.Board[0][i] == g.Board[2][i] {
            g.Winner = g.Board[0][i]
            g.Over = true
            return
        }
    }

    // Check diagonal lines
    if g.Board[0][0] != "" && g.Board[0][0] == g.Board[1][1] && g.Board[0][0] == g.Board[2][2] {
        g.Winner = g.Board[0][0]
        g.Over = true
        return
    }
    if g.Board[0][2] != "" && g.Board[0][2] == g.Board[1][1] && g.Board[0][2] == g.Board[2][0] {
        g.Winner = g.Board[0][2]
        g.Over = true
        return
    }

    // Check for a draw
    isDraw := true
    for i := 0; i < 3; i++ {
        for j := 0; j < 3; j++ {
            if g.Board[i][j] == "" {
                isDraw = false
                break
            }
        }
        if !isDraw {
            break
        }
    }

    if isDraw {
        g.Winner = "Draw"
        g.Over = true
    }
}

func (g *Game) Serialize() (string, error) {
    data, err := json.Marshal(g)
    if err != nil {
        return "", err
    }
    return string(data), nil
}
