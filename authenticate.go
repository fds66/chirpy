package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/google/uuid"

	"workspace/fds66/github.com/fds66/chirpy/internal/auth"
)

func (cfg *apiConfig) AuthenticateUserByJWT(header http.Header) (uuid.UUID, error) {

	token, err := auth.GetBearerToken(header)
	if err != nil {
		log.Printf("Error extracting token %v", err)
		errorMessage := fmt.Errorf("Incorrect Token, %v", err)
		return uuid.UUID{}, errorMessage
	}
	userJWT, err := auth.ValidateJWT(token, cfg.secret)
	fmt.Printf("returned userJWT and error from ValidateJWT %v, %v\n", userJWT, err)
	if err != nil {
		log.Printf("Error checking JWT %v", err)
		errorMessage := fmt.Errorf("Incorrect Token, %v", err)
		return uuid.UUID{}, errorMessage
	}
	return userJWT, nil

}
