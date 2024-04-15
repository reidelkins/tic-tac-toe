package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/reidelkins/kube-tic-tac-toe/internal/api"
	"github.com/reidelkins/kube-tic-tac-toe/internal/db"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"github.com/joho/godotenv"
	log "github.com/sirupsen/logrus"
)

// dbConn represents the database connection
var dbConn *db.DB

func main() {
	log.SetReportCaller(true)
    var router *chi.Mux = chi.NewRouter()
	// handlers.Handler(router)

	var err error

	// Load .env file
    err = godotenv.Load()
    if err != nil {
        log.Fatal("Error loading .env file")
    }

	dbHost := os.Getenv("DB_HOST")
    dbPort := os.Getenv("DB_PORT")
    dbName := os.Getenv("DB_NAME")
    dbUser := os.Getenv("DB_USER")
    dbPassword := os.Getenv("DB_PASSWORD")
    dbSslmode := os.Getenv("DB_SSLMODE")

	dsn := fmt.Sprintf("host=%s port=%s dbname=%s user=%s password=%s sslmode=%s",
        dbHost, dbPort, dbName, dbUser, dbPassword, dbSslmode)

	
	
    dbConn, err = db.NewDB(dsn)
	
    if err != nil {
        log.Fatalf("Could not connect to database: %v", err)
    }

	handler := &api.Handler{DBConn: dbConn}

    // Setting up CORS
	cors := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:4200", "http://localhost", "https://c15f-70-116-143-2.ngrok-free.app", }, // or your frontend host
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	})
	router.Use(cors.Handler)
    // Print the url of the incoming request
	router.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			log.Infof("Request URL: %s", r.URL)
			next.ServeHTTP(w, r)
		})
	})
	router.Post("/login/google", handler.GoogleLoginHandler)
	router.Post("/create-game", handler.CreateGameHandler)
	router.Post("/join-game", handler.JoinGameHandler)
    router.Get("/list-active-games", handler.ListGamesHandler)
	router.Get("/ws", handler.WebSocketHandler)
	router.Get("/get-game/{gameId}", handler.GetGameHandler)
	
    // router.Post("/join-game/{gameId}", joinGameHandler)
    
	
    // router.Post("/play-move", handler.PlayMoveHandler)

    log.Println("Starting server on :8080")
    err = http.ListenAndServe(":8080", router)
	if err != nil {
		log.Error(err)
	}
}