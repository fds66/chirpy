package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"sync/atomic"

	"workspace/fds66/github.com/fds66/chirpy/internal/database"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

type apiConfig struct {
	fileserverHits atomic.Int32
	db             *database.Queries
	platform       string
	secret         string
}

func main() {
	godotenv.Load()
	dbURL := os.Getenv("DB_URL")
	if dbURL == "" {
		log.Fatal("DB_URL must be set")
	}
	platformString := os.Getenv("PLATFORM")
	if platformString == "" {
		log.Fatal("PLATFORM must be set")
	}
	fmt.Printf("dbURL %s\n", dbURL)
	dbConn, err := sql.Open("postgres", dbURL)
	if err != nil {
		fmt.Printf("error opening database connections")
		os.Exit(1)
	}
	dbQueries := database.New(dbConn)
	secret := os.Getenv("SECRET")

	const filepathRoot = "."
	const port = "8080"
	apiCfg := apiConfig{
		fileserverHits: atomic.Int32{},
		db:             dbQueries,
		platform:       platformString,
		secret:         secret,
	}

	//Initialise the Mux
	serveMux := http.NewServeMux()
	//fileserver with built-in handler
	handler := http.StripPrefix("/app", http.FileServer(http.Dir(filepathRoot)))
	serveMux.Handle("/app/", apiCfg.middlewareMetricsInc(handler))

	//register our handlers, this wraps them in HandleFunc which gives them the ServeHTTP method
	serveMux.HandleFunc("GET /api/healthz", handlerReady)
	serveMux.HandleFunc("GET /admin/metrics", apiCfg.handlerMetrics)
	serveMux.HandleFunc("POST /admin/reset", apiCfg.handlerReset)
	serveMux.HandleFunc("POST /api/users", apiCfg.handlerUsersCreate)
	serveMux.HandleFunc("PUT /api/users", apiCfg.handlerUsersUpdate)
	serveMux.HandleFunc("POST /api/login", apiCfg.handlerUsersLogin)
	serveMux.HandleFunc("POST /api/chirps", apiCfg.handlerCreateChirps)
	serveMux.HandleFunc("GET /api/chirps", apiCfg.handlerGetChirps)
	serveMux.HandleFunc("GET /api/chirps/{chirpID}", apiCfg.handlerGetChirpByID)
	serveMux.HandleFunc("DELETE /api/chirps/{chirpID}", apiCfg.handlerDeleteChirpByID)
	serveMux.HandleFunc("POST /api/refresh", apiCfg.handlerRefresh)
	serveMux.HandleFunc("POST /api/revoke", apiCfg.handlerRevoke)

	svr := http.Server{
		Addr:    ":" + port,
		Handler: serveMux,
	}

	fmt.Println("serving")
	log.Printf("Serving files from %s on port: %s\n", filepathRoot, port)
	log.Fatal(svr.ListenAndServe())
}
