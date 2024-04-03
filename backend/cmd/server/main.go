package main

import (
	// "fmt"
	"net/http"

	"github.com/reidelkins/kube-tic-tac-toe/internal/api"

	"github.com/go-chi/chi/v5"
	log "github.com/sirupsen/logrus"
)

func main() {
	log.SetReportCaller(true)
    var router *chi.Mux = chi.NewRouter()
	// handlers.Handler(router)

    // router.Post("/create-game", createGameHandler)
    router.Get("/list-active-games", api.ListGamesHandler)
	router.Get("/test-route", api.TestRouteHandler)
    // router.Post("/join-game/{gameId}", joinGameHandler)
    // router.Post("/play-move", playMoveHandler)
    router.Post("/create-game", api.CreateGameHandler)
    // router.Post("/play-move", api.PlayMoveHandler)

    log.Println("Starting server on :8080")
    err := http.ListenAndServe(":8080", router)
	if err != nil {
		log.Error(err)
	}
}

// func createGameHandler(w http.ResponseWriter, r *http.Request) {
//     fmt.Fprint(w, "Create game endpoint")
// }

// func listGamesHandler(w http.ResponseWriter, r *http.Request) {
//     fmt.Fprint(w, "List games endpoint")
// }

// func joinGameHandler(w http.ResponseWriter, r *http.Request) {
//     fmt.Fprint(w, "Join game endpoint")
// }

// func playMoveHandler(w http.ResponseWriter, r *http.Request) {
//     fmt.Fprint(w, "Play move endpoint")
// }
