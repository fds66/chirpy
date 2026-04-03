package main

import (
	"fmt"
	//"log"
	"net/http"
	//"sync/atomic"
)

func (cfg *apiConfig) handlerMetrics(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(http.StatusOK)
	message := fmt.Sprintf("<html>  <body>    <h1>Welcome, Chirpy Admin</h1>    <p>Chirpy has been visited %d times!</p>  </body></html>", cfg.fileserverHits.Load())
	w.Write([]byte(message))
}

func (cfg *apiConfig) handlerReset(w http.ResponseWriter, r *http.Request) {
	cfg.fileserverHits.Store(0)
}
