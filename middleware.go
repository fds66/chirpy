package main

import (
	"fmt"
	//"log"
	"net/http"
	//"sync/atomic"
)

// middleware method on *apiConfig to increment the fileserverHits counter every time it is called

func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		newCount := cfg.fileserverHits.Add(1)
		fmt.Printf("new count %v\n", newCount)
		next.ServeHTTP(w, r)
	})
}
