package main

import (
	"context"
	"encoding/json"

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
	userJWT, err := cfg.AuthenticateUserByJWT(r.Header)
	if err != nil {
		log.Printf("Error extracting token %v\n", err)
		respondWithError(w, http.StatusUnauthorized, "Unauthorised", err)
		return
	}
	// authenticate user via JWT
	/*token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		log.Printf("Error extracting token %v", err)
		respondWithError(w, http.StatusUnauthorized, "Incorrect Token", err)
		return
	}
	userJWT, err := auth.ValidateJWT(token, cfg.secret)
	fmt.Printf("returned userJWT and error from ValidateJWT %v, %v\n", userJWT, err)
	if err != nil {
		log.Printf("Error checking JWT %v", err)
		respondWithError(w, http.StatusUnauthorized, "Incorrect Token", err)
		return
	}
	*/

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
	chirps, err := cfg.db.GetAllChirpsSorted(context.Background())
	if err != nil {
		log.Printf("Error creating chirp record %v", err)
		respondWithError(w, http.StatusInternalServerError, "Something went wrong", err)
		return
	}
	allChirps := make([]Chirp, len(chirps))
	for i, chirp := range chirps {
		allChirps[i] = convertDatabaseToLocalChirp(&chirp)

	}
	respondWithJSON(w, 200, allChirps)
}

func (cfg *apiConfig) handlerGetChirpByID(w http.ResponseWriter, r *http.Request) {
	//func (r *Request) PathValue(name string) string
	chirpIDString := r.PathValue("chirpID")
	if chirpIDString == "" {
		log.Printf("No chirp ID found")
		respondWithError(w, http.StatusBadRequest, "No chirp ID found", nil)
		//bad request 400
		return
	}
	// convert incoming string ID to a uuid.UUID value
	chirpID, err := uuid.Parse(chirpIDString)
	if err != nil {
		log.Printf("Error converting chirp ID %v", err)
		respondWithError(w, http.StatusInternalServerError, "Something went wrong", err)
		return
	}

	foundChirp, err := cfg.db.GetChirpByID(context.Background(), chirpID)
	if err != nil {
		log.Printf("Error retrieving chirp from database %v", err)
		respondWithError(w, 404, "Something went wrong", err)
		return
	}
	respBody := convertDatabaseToLocalChirp(&foundChirp)
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
