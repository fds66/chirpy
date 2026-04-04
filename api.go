package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"
)

type parameters struct {
	Body string `json:"body"`
}

type returnErrorVals struct {
	ErrorMessage string `json:"error"`
}

/*
	type returnJsonVals struct {
		Valid bool `json:"valid"`
	}
*/
type returnJsonVals struct {
	CleanedBody string `json:"cleaned_body"`
}

func handlerValidateChirp(w http.ResponseWriter, r *http.Request) {
	const maxChirpLength = 140
	//fmt.Println("validate handler")
	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		log.Printf("Error decoding request %v", err)
		respondWithError(w, http.StatusInternalServerError, "Something went wrong", err)
		return
	}

	if len(params.Body) > maxChirpLength {
		log.Printf("Chirp is too long")
		respondWithError(w, http.StatusBadRequest, "Chirp is too long", nil)
		//bad request 400
		return
	}
	cleanedChirp := cleanString(params.Body)
	respBody := returnJsonVals{
		CleanedBody: cleanedChirp,
	}

	respondWithJSON(w, 200, respBody)

}

func respondWithError(w http.ResponseWriter, code int, msg string, err error) {
	if err != nil {
		log.Println(err)
	}
	if code > 499 {
		log.Printf("Responding with 5XX error: %s", msg)
	}
	respBody := returnErrorVals{
		ErrorMessage: msg,
	}
	respondWithJSON(w, code, respBody)

}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	dat, err := json.Marshal(payload)
	if err != nil {
		log.Printf("error marshalling JSON %v", err)
		w.WriteHeader(500)
		return
	}

	w.WriteHeader(code)
	w.Write(dat)
}

func cleanString(inputString string) string {
	badWords := []string{"kerfuffle", "sharbert", "fornax"}
	wordList := strings.Split(inputString, " ")
	for i, word := range wordList {
		for _, badWord := range badWords {
			if badWord == strings.ToLower(word) {
				wordList[i] = "****"
			}
		}

	}
	return strings.Join(wordList, " ")
}
