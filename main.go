package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	const filepathRoot = "."
	const port = "8080"
	//Initialise the Mux
	serveMux := http.NewServeMux()
	//fileserver with built-in handler
	serveMux.Handle("/app/", http.StripPrefix("/app", http.FileServer(http.Dir(filepathRoot))))
	server := http.Server{
		Addr:    ":" + port,
		Handler: serveMux,
	}
	//register our handler
	serveMux.HandleFunc("/healthz", handlerReady)

	fmt.Println("serving")
	log.Printf("Serving files from %s on port: %s\n", filepathRoot, port)
	log.Fatal(server.ListenAndServe())
}

// handler that gives the ready response
func handlerReady(w http.ResponseWriter, r *http.Request) {
	//Header
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}
