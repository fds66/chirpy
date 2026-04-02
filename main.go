package main

import (
	"fmt"
	"log"
	"net/http"
	"sync/atomic"
)

type apiConfig struct {
	fileserverHits atomic.Int32
}

func main() {
	const filepathRoot = "."
	const port = "8080"
	apiCfg := apiConfig{
		fileserverHits: atomic.Int32{},
	}

	//Initialise the Mux
	serveMux := http.NewServeMux()
	//fileserver with built-in handler
	handler := http.StripPrefix("/app", http.FileServer(http.Dir(filepathRoot)))
	serveMux.Handle("/app/", apiCfg.middlewareMetricsInc(handler))

	//register our handlers, this wraps them in HandleFunc which gives them the ServeHTTP method
	serveMux.HandleFunc("GET /api/healthz", handlerReady)
	serveMux.HandleFunc("GET /api/metrics", apiCfg.handlerMetrics)
	serveMux.HandleFunc("POST /api/reset", apiCfg.handlerReset)

	svr := http.Server{
		Addr:    ":" + port,
		Handler: serveMux,
	}

	fmt.Println("serving")
	log.Printf("Serving files from %s on port: %s\n", filepathRoot, port)
	log.Fatal(svr.ListenAndServe())
}

// handler that gives the ready response
func handlerReady(w http.ResponseWriter, r *http.Request) {
	//Header
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(http.StatusText(http.StatusOK)))
}

// middleware method on *apiConfig to increment the fileserverHits counter every time it is called

func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		newCount := cfg.fileserverHits.Add(1)
		fmt.Printf("new count %v\n", newCount)
		next.ServeHTTP(w, r)
	})
}
func (cfg *apiConfig) handlerMetrics(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	message := fmt.Sprintf("Hits: %v", cfg.fileserverHits.Load())
	w.Write([]byte(message))
}

func (cfg *apiConfig) handlerReset(w http.ResponseWriter, r *http.Request) {
	cfg.fileserverHits.Store(0)
}
