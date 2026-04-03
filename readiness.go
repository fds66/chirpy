package main

import (
	"net/http"
)

// handler that gives the ready response
func handlerReady(w http.ResponseWriter, r *http.Request) {
	//Header
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(http.StatusText(http.StatusOK)))
}
