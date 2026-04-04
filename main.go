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
	dbQ            *database.Queries
}

func main() {
	godotenv.Load()
	dbURL := os.Getenv("DB_URL")
	fmt.Printf("dbURL %s\n", dbURL)
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		fmt.Printf("error opening database connections")
		os.Exit(1)
	}
	dbQueries := database.New(db)

	const filepathRoot = "."
	const port = "8080"
	apiCfg := apiConfig{
		fileserverHits: atomic.Int32{},
		dbQ:            dbQueries,
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
	serveMux.HandleFunc("POST /api/validate_chirp", handlerValidateChirp)

	svr := http.Server{
		Addr:    ":" + port,
		Handler: serveMux,
	}

	fmt.Println("serving")
	log.Printf("Serving files from %s on port: %s\n", filepathRoot, port)
	log.Fatal(svr.ListenAndServe())
}
