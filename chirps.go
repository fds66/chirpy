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

func (cfg *apiConfig) handlerCreateChirps(w http.ResponseWriter, r *http.Request) {
	const maxChirpLength = 140

	decoder := json.NewDecoder(r.Body)
	params := CreateChirpParams{}
	err := decoder.Decode(&params)
	if err != nil {
		log.Printf("Error decoding request %v", err)
		respondWithError(w, http.StatusInternalServerError, "Something went wrong", err)
		return
	}
	// authenticate user via JWT
	userJWT, err := cfg.AuthenticateUserByJWT(r.Header)
	if err != nil {
		log.Printf("Error extracting token %v\n", err)
		respondWithError(w, http.StatusUnauthorized, "Unauthorised", err)
		return
	}
	// check if chirp follows rules
	if len(params.Body) > maxChirpLength {
		log.Printf("Chirp is too long")
		respondWithError(w, http.StatusBadRequest, "Chirp is too long", nil)
		//bad request 400
		return
	}
	cleanedBody := cleanString(params.Body)

	createParams := database.CreateChirpParams{
		Body:   cleanedBody,
		UserID: userJWT,
	}

	createdChirp, err := cfg.db.CreateChirp(context.Background(), createParams)
	if err != nil {
		log.Printf("Error creating chirp record %v", err)
		respondWithError(w, http.StatusInternalServerError, "Something went wrong", err)
		return
	}

	respBody := convertDatabaseToLocalChirp(&createdChirp)
	respondWithJSON(w, 201, respBody)

}
func convertDatabaseToLocalChirp(in *database.Chirp) Chirp {
	return Chirp{
		ID:        in.ID,
		CreatedAt: in.CreatedAt,
		UpdatedAt: in.UpdatedAt,
		Body:      in.Body,
		UserID:    in.UserID,
	}
}

func (cfg *apiConfig) handlerGetChirps(w http.ResponseWriter, r *http.Request) {

	/*
			s := r.URL.Query().Get("author_id")
		// s is a string that contains the value of the author_id query parameter
		// if it exists, or an empty string if it doesn't
	*/
	var chirps []database.Chirp
	author := r.URL.Query().Get("author_id")
	var err error

	if author != "" {
		authorID, err := uuid.Parse(author)
		if err != nil {
			log.Printf("Error converting author ID %v", err)
			respondWithError(w, http.StatusInternalServerError, "Something went wrong", err)
			return
		}
		chirps, err = cfg.db.GetChirpsbyUserID(context.Background(), authorID)
	} else {
		chirps, err = cfg.db.GetAllChirpsSorted(context.Background())
		if err != nil {
			log.Printf("Error creating chirp record %v", err)
			respondWithError(w, http.StatusInternalServerError, "Something went wrong", err)
			return
		}
	}

	allChirps := make([]Chirp, len(chirps))
	for i, chirp := range chirps {
		allChirps[i] = convertDatabaseToLocalChirp(&chirp)

	}
	respondWithJSON(w, 200, allChirps)
}

func (cfg *apiConfig) handlerGetChirpByID(w http.ResponseWriter, r *http.Request) {

	foundChirp, err, code := cfg.chirpFromIDInRequest(r)
	if err != nil {
		log.Printf("Error retrieving chirp from database %v", err)
		respondWithError(w, code, err.Error(), err)
		return
	}

	respBody := convertDatabaseToLocalChirp(&foundChirp)
	respondWithJSON(w, 200, respBody)

}

func (cfg *apiConfig) handlerDeleteChirpByID(w http.ResponseWriter, r *http.Request) {
	//func (r *Request) PathValue(name string) string
	// authenticate user by JWT token
	userJWT, err := cfg.AuthenticateUserByJWT(r.Header)
	if err != nil {
		log.Printf("Error extracting token %v\n", err)
		respondWithError(w, http.StatusUnauthorized, "Unauthorised", err)
		return
	}
	foundChirp, err, code := cfg.chirpFromIDInRequest(r)
	if err != nil {
		log.Printf("Error retrieving chirp from database %v", err)
		respondWithError(w, code, err.Error(), err)
		return
	}
	if userJWT != foundChirp.UserID {
		log.Printf("UserID does not match chirp %v", err)
		respondWithError(w, http.StatusForbidden, "Unauthorised", err)
		return
	}
	err = cfg.db.DeleteChirpByID(context.Background(), foundChirp.ID)
	if err != nil {
		log.Printf("Error deleting chirp from database %v", err)
		respondWithError(w, http.StatusInternalServerError, "something went wrong", err)
		return
	}
	// if successful
	w.WriteHeader(http.StatusNoContent)

}

func (cfg *apiConfig) chirpFromIDInRequest(r *http.Request) (database.Chirp, error, int) {
	errorMessage := fmt.Errorf("Error")
	chirpIDString := r.PathValue("chirpID")
	if chirpIDString == "" {
		log.Printf("No chirp ID found")
		errorMessage = fmt.Errorf("No chirp ID found")
		return database.Chirp{}, errorMessage, http.StatusBadRequest
	}
	// convert incoming string ID to a uuid.UUID value
	chirpID, err := uuid.Parse(chirpIDString)
	if err != nil {
		log.Printf("Error converting chirp ID %v", err)
		errorMessage = fmt.Errorf("Something went wrong")
		return database.Chirp{}, errorMessage, http.StatusInternalServerError
	}
	var foundChirp database.Chirp
	foundChirp, err = cfg.db.GetChirpByID(context.Background(), chirpID)
	if err != nil {
		log.Printf("Error retrieving chirp from database %v", err)
		errorMessage = fmt.Errorf("Something went wrong")
		return database.Chirp{}, errorMessage, http.StatusNotFound
	}
	return foundChirp, nil, 0
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
