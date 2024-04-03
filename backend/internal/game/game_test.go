package game

import "testing"

func TestGamePlayMoveAndCheckWinner(t *testing.T) {
    scenarios := []struct {
        player1ID int64
        player2ID int64
        moves     []struct {
            playerID int64
            x        int
            y        int
        }
        winner string
        isOver bool
    }{
        { // Horizontal win
            player1ID: 1,
            player2ID: 2,
            moves: []struct{ playerID int64; x int; y int }{
                {1, 0, 0}, {2, 1, 0}, {1, 0, 1}, {2, 1, 1}, {1, 0, 2},
            },
            winner: "X",
            isOver: true,
        },
        // { // Vertical win
        //     moves:  []struct{ x, y int }{{0, 0}, {0, 1}, {1, 0}, {1, 1}, {2, 0}},
        //     winner: "X",
        //     isOver: true,
        // },
        // { // Diagonal win
        //     moves:  []struct{ x, y int }{{0, 0}, {0, 1}, {1, 1}, {1, 2}, {2, 2}},
        //     winner: "X",
        //     isOver: true,
        // },
        // { // Draw
        //     moves:  []struct{ x, y int }{{0, 0}, {0, 1}, {0, 2}, {1, 0}, {1, 1}, {1, 2}, {2, 0}, {2, 1}, {2, 2}},
        //     winner: "Draw",
        //     isOver: true,
        // },
    }

    for _, scenario := range scenarios {
        game := NewGame(scenario.player1ID)

        for _, move := range scenario.moves {
            game.PlayMove(move.playerID, move.x, move.y)
        }

        if game.Winner != scenario.winner {
            t.Errorf("Expected winner to be %s, got %s", scenario.winner, game.Winner)
        }
        if game.Over != scenario.isOver {
            t.Errorf("Expected game over to be %v, got %v", scenario.isOver, game.Over)
        }
    }
}
