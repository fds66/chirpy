package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"

	"workspace/fds66/github.com/fds66/chirpy/internal/auth"
	"workspace/fds66/github.com/fds66/chirpy/internal/database"
)

type userJsonStruct struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Email     string    `json:"email"`
}

func (cfg *apiConfig) handlerUsersCreate(w http.ResponseWriter, r *http.Request) {
	type InputParameters struct {
		Password string `json:"password"`
		Email    string `json:"email"`
	}
	decoder := json.NewDecoder(r.Body)
	inputParams := InputParameters{}
	err := decoder.Decode(&inputParams)
	if err != nil {
		log.Printf("Error decoding request %v", err)
		respondWithError(w, http.StatusInternalServerError, "Something went wrong", err)
		return
	}
	//fmt.Printf("email %s", params.Email)
	if len(inputParams.Email) == 0 {
		log.Printf("No email address found")
		respondWithError(w, http.StatusBadRequest, "No email address found", nil)
		//bad request 400
		return
	}
	if len(inputParams.Password) == 0 {
		log.Printf("No password found")
		respondWithError(w, http.StatusBadRequest, "No password found", nil)
		//bad request 400
		return
	}
	hashedPassword, err := auth.HashPassword(inputParams.Password)
	if err != nil {
		log.Printf("Error hashing password %v", err)
		respondWithError(w, http.StatusInternalServerError, "Something went wrong", err)
		return
	}

	createParams := database.CreateUserParams{
		Email:          inputParams.Email,
		HashedPassword: hashedPassword,
	}
	newUser, err := cfg.db.CreateUser(context.Background(), createParams)
	if err != nil {
		log.Printf("Error creating user %v", err)
		if strings.Contains(err.Error(), "duplicate key") {
			respondWithError(w, http.StatusInternalServerError, "Duplicate user", err)
		} else {
			respondWithError(w, http.StatusInternalServerError, "Something went wrong", err)
		}

		return
	}
	respBody := userJsonStruct{
		ID:        newUser.ID,
		CreatedAt: newUser.CreatedAt,
		UpdatedAt: newUser.UpdatedAt,
		Email:     newUser.Email,
	}
	respondWithJSON(w, 200, respBody)

}
