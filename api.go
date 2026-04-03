package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
)

func handlerValidateChirp(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Body string `json:"body"`
	}

	//fmt.Println("validate handler")
	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		log.Printf("Error decoding request %v", err)
		respondWithError(w, 500, "Something went wrong")
		return
	}

	if len(params.Body) > 140 {
		log.Printf("Chirp is too long")
		respondWithError(w, 400, "Chirp is too long")
		return
	}
	type returnJsonVals struct {
		Valid bool `json:"valid"`
	}
	respBody := returnJsonVals{
		Valid: true,
	}

	respondWithJSON(w, 200, respBody)

}

func respondWithError(w http.ResponseWriter, code int, msg string) {
	type returnErrorVals struct {
		ErrorMessage string `json:"error"`
	}
	respBody := returnErrorVals{
		ErrorMessage: msg,
	}
	dat, err := json.Marshal(respBody)
	if err != nil {
		log.Printf("error marshalling JSON %v", err)
		os.Exit(1)
	}

	//Header
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(dat)
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {

	dat, err := json.Marshal(payload)
	if err != nil {
		log.Printf("error marshalling JSON %v", err)
		respondWithError(w, 500, "Something went wrong")
		return
	}

	//Header
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(dat)

}
