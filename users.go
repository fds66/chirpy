package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
)

/*
request

	{
	  "email": "user@example.com"
	}

response HTTP 201 created

	{
	  "id": "50746277-23c6-4d85-a890-564c0044c2fb",
	  "created_at": "2021-07-07T00:00:00Z",
	  "updated_at": "2021-07-07T00:00:00Z",
	  "email": "user@example.com"
	}
*/
type userJsonStruct struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Email     string    `json:"email"`
}

func (cfg *apiConfig) handlerUsersCreate(w http.ResponseWriter, r *http.Request) {
	type emailParameters struct {
		Email string `json:"email"`
	}
	decoder := json.NewDecoder(r.Body)
	params := emailParameters{}
	err := decoder.Decode(&params)
	if err != nil {
		log.Printf("Error decoding request %v", err)
		respondWithError(w, http.StatusInternalServerError, "Something went wrong", err)
		return
	}
	//fmt.Printf("email %s", params.Email)
	if len(params.Email) == 0 {
		log.Printf("No email address found")
		respondWithError(w, http.StatusBadRequest, "No email address found", nil)
		//bad request 400
		return
	}
	newUser, err := cfg.db.CreateUser(context.Background(), params.Email)
	if err != nil {
		log.Printf("Error creating user %v", err)
		respondWithError(w, http.StatusInternalServerError, "Something went wrong", err)
		return
	}
	respBody := userJsonStruct{
		ID:        newUser.ID,
		CreatedAt: newUser.CreatedAt,
		UpdatedAt: newUser.UpdatedAt,
		Email:     newUser.Email,
	}
	respondWithJSON(w, 201, respBody)

}
