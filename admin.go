package main

import (
	"fmt"
	"log"
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
	if cfg.platform != "dev" {
		log.Printf("Error accessing command, restricted access ")
		respondWithError(w, http.StatusForbidden, "Reset only allowed in dev environment", nil)
		return
	}
	cfg.fileserverHits.Store(0)
	err := cfg.db.ResetUsers(r.Context())
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to reset database", nil)
		return

	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Hits reset to 0 and database reset to initial state."))
}
