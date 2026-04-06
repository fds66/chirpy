package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"workspace/fds66/github.com/fds66/chirpy/internal/database"

	"github.com/google/uuid"
)

//type parameters struct {
//Body string `json:"body"`
//}

type returnErrorVals struct {
	ErrorMessage string `json:"error"`
}

type returnJsonVals struct {
	CleanedBody string `json:"cleaned_body"`
}

/*
accepts

		{
		  "body": "Hello, world!",
		  "user_id": "123e4567-e89b-12d3-a456-426614174000"
		}

		  respond 201 and
		  {
		  "id": "94b7e44c-3604-42e3-bef7-ebfcc3efff8f",
		  "created_at": "2021-01-01T00:00:00Z",
		  "updated_at": "2021-01-01T00:00:00Z",
		  "body": "Hello, world!",
		  "user_id": "123e4567-e89b-12d3-a456-426614174000"
		}
		database type
		type CreateChirpParams struct {
			Body   string
			UserID uuid.UUID
		}
		type apiConfig struct {
		fileserverHits atomic.Int32
		db             *database.Queries
		platform       string
	}
*/
type CreateChirpParams struct {
	Body   string    `json:"body"`
	UserID uuid.UUID `json:"user_id"`
}

type Chirp struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Body      string    `json:"body"`
	UserID    uuid.UUID `json:"user_id"`
}

func (cfg *apiConfig) handlerChirps(w http.ResponseWriter, r *http.Request) {
	const maxChirpLength = 140
	//fmt.Println("validate handler")
	decoder := json.NewDecoder(r.Body)
	params := CreateChirpParams{}
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
	cleanedBody := cleanString(params.Body)

	createParams := database.CreateChirpParams{
		Body:   cleanedBody,
		UserID: params.UserID,
	}
	fmt.Printf("params %s\n%v\n, createParams %s\n%v\n", params.Body, params.UserID, createParams.Body, createParams.UserID)

	createdChirp, err := cfg.db.CreateChirp(context.Background(), createParams)
	if err != nil {
		log.Printf("Error creating chirp record %v", err)
		respondWithError(w, http.StatusInternalServerError, "Something went wrong", err)
		return
	}
	fmt.Printf("createdChirp %+v\n", createdChirp)
	respBody := Chirp{
		ID:        createdChirp.ID,
		CreatedAt: createdChirp.CreatedAt,
		UpdatedAt: createdChirp.UpdatedAt,
		Body:      createdChirp.Body,
		UserID:    createdChirp.UserID,
	}
	fmt.Printf("respBody %+v\n", respBody)

	respondWithJSON(w, 201, respBody)

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
