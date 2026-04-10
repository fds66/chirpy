package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"

	"workspace/fds66/github.com/fds66/chirpy/internal/auth"
	"workspace/fds66/github.com/fds66/chirpy/internal/database"
)

type userJsonStruct struct {
	ID           uuid.UUID `json:"id"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
	Email        string    `json:"email"`
	Token        string    `json:"token"`
	RefreshToken string    `json:"refresh_token"`
	IsChirpyRed  bool      `json:"is_chirpy_red"`
}

type UserInputParameters struct {
	Password string `json:"password"`
	Email    string `json:"email"`
}

func (cfg *apiConfig) handlerUsersCreate(w http.ResponseWriter, r *http.Request) {

	inputParams, err := userInput(w, r)
	if err != nil {
		log.Printf("Error validating user input parameters %v", err)
		respondWithError(w, http.StatusInternalServerError, "Something went wrong", err)
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
		ID:          newUser.ID,
		CreatedAt:   newUser.CreatedAt,
		UpdatedAt:   newUser.UpdatedAt,
		Email:       newUser.Email,
		IsChirpyRed: newUser.IsChirpyRed,
	}
	respondWithJSON(w, 201, respBody)

}

func userInput(w http.ResponseWriter, r *http.Request) (UserInputParameters, error) {
	decoder := json.NewDecoder(r.Body)
	inputParams := UserInputParameters{}
	err := decoder.Decode(&inputParams)
	if err != nil {
		log.Printf("Error decoding request %v", err)
		respondWithError(w, http.StatusInternalServerError, "Something went wrong", err)
		return UserInputParameters{}, err
	}
	//fmt.Printf("email %s", params.Email)
	if len(inputParams.Email) == 0 {
		log.Printf("No email address found")
		respondWithError(w, http.StatusBadRequest, "No email address found", nil)
		//bad request 400, err
	}
	if len(inputParams.Password) == 0 {
		log.Printf("No password found")
		respondWithError(w, http.StatusBadRequest, "No password found", nil)
		//bad request 400
		return UserInputParameters{}, err
	}

	return inputParams, nil

}

func (cfg *apiConfig) handlerUsersLogin(w http.ResponseWriter, r *http.Request) {
	// default expiration of the JWT token in seconds, 1hour = 3600 secs

	inputParams, err := userInput(w, r)
	if err != nil {
		log.Printf("Error validating user input parameters %v", err)
		respondWithError(w, http.StatusInternalServerError, "Something went wrong", err)
		return
	}

	//lookup user in database by email
	user, err := cfg.db.GetUserByEmail(context.Background(), inputParams.Email)
	if err != nil {
		log.Printf("Error looking up user by email %v", err)
		respondWithError(w, http.StatusUnauthorized, "Incorrect email or password", err)
		return
	}

	passwordCheck, err := auth.CheckPasswordHash(inputParams.Password, user.HashedPassword)
	if err != nil {
		log.Printf("Error checking password %v", err)
		respondWithError(w, http.StatusUnauthorized, "Incorrect email or password", err)
		return
	}
	if passwordCheck == false {
		log.Printf("Password check failed ")
		respondWithError(w, http.StatusUnauthorized, "Incorrect email or password", err)
		return
	}
	// user password correct so generate JWT token
	//JWT token expires after an hour
	defaultExpire := 3600
	expireTime := time.Duration(defaultExpire) * time.Second
	userToken, err := auth.MakeJWT(user.ID, cfg.secret, expireTime)
	if err != nil {
		log.Printf("Error making JWT %v", err)
		respondWithError(w, http.StatusUnauthorized, "Something went wrong", err)
		return
	}
	refreshT := auth.MakeRefreshToken()
	/*
			type CreateRefreshTokenParams struct {
			Token     string
			UserID    uuid.UUID
			ExpiresAt time.Time
			RevokedAt sql.NullTime
		}
	*/
	//expires in 60 days
	expiration_date := time.Now()
	expiration_date.AddDate(0, 0, 60)
	refreshParams := database.CreateRefreshTokenParams{
		Token:     refreshT,
		UserID:    user.ID,
		ExpiresAt: expiration_date,
	}
	newRefreshToken, err := cfg.db.CreateRefreshToken(context.Background(), refreshParams)

	respBody := userJsonStruct{
		ID:           user.ID,
		CreatedAt:    user.CreatedAt,
		UpdatedAt:    user.UpdatedAt,
		Email:        user.Email,
		Token:        userToken,
		RefreshToken: newRefreshToken.Token,
		IsChirpyRed:  user.IsChirpyRed,
	}
	respondWithJSON(w, 200, respBody)
}
func (cfg *apiConfig) handlerRefresh(w http.ResponseWriter, r *http.Request) {
	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		log.Printf("Error extracting token %v", err)
		respondWithError(w, http.StatusUnauthorized, "Incorrect Token", err)
		return
	}
	tokenRecord, err := cfg.db.GetRefreshTokenByToken(context.Background(), token)
	if err != nil {
		log.Printf("Error finding token in database %v", err)
		respondWithError(w, http.StatusUnauthorized, "Not authorised", err)
		return
	}
	currentTime := time.Now()
	if currentTime.After(tokenRecord.ExpiresAt) {
		log.Printf("Refresh token expired %v", err)
		respondWithError(w, http.StatusUnauthorized, "Not authorised", err)
		return
	}
	if tokenRecord.RevokedAt.Valid == true {
		log.Printf("Refresh token revoked %v", err)
		respondWithError(w, http.StatusUnauthorized, "Not authorised", err)
		return
	}
	defaultExpire := 3600
	expireTime := time.Duration(defaultExpire) * time.Second
	newToken, err := auth.MakeJWT(tokenRecord.UserID, cfg.secret, expireTime)
	if err != nil {
		log.Printf("Error making JWT %v", err)
		respondWithError(w, http.StatusInternalServerError, "Something went wrong", err)
		return
	}
	type outputToken struct {
		Token string `json:"token"`
	}
	respBody := outputToken{
		Token: newToken,
	}
	respondWithJSON(w, 200, respBody)
}

/*
	type RefreshToken struct {
		Token     string
		CreatedAt time.Time
		UpdatedAt time.Time
		UserID    uuid.UUID
		ExpiresAt time.Time
		RevokedAt sql.NullTime
	}
*/
func (cfg *apiConfig) handlerRevoke(w http.ResponseWriter, r *http.Request) {
	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		log.Printf("Error extracting token %v", err)
		respondWithError(w, http.StatusUnauthorized, "Incorrect Token", err)
		return
	}
	tokenRecord, err := cfg.db.GetRefreshTokenByToken(context.Background(), token)
	if err != nil {
		log.Printf("Error finding token in database %v", err)
		respondWithError(w, http.StatusUnauthorized, "Not authorised", err)
		return
	}

	cfg.db.RevokeToken(context.Background(), tokenRecord.Token)
	//check it worked
	/*
		updatedRecord, err := cfg.db.GetRefreshTokenByToken(context.Background(), tokenRecord.Token)
		if err != nil {
			log.Printf("Error finding token in database %v", err)

			return
		}
		log.Printf("updated record to check if it worked %+v\n", updatedRecord)
	*/
	w.WriteHeader(204)

}
func (cfg *apiConfig) handlerUsersUpdate(w http.ResponseWriter, r *http.Request) {
	// check the access token and extract the user id
	userJWT, err := cfg.AuthenticateUserByJWT(r.Header)
	if err != nil {
		log.Printf("Error extracting token %v\n", err)
		respondWithError(w, http.StatusUnauthorized, "Unauthorised", err)
		return
	}

	// create the params struct to update the database from the body of the request
	inputParams, err := userInput(w, r)
	if err != nil {
		log.Printf("Error validating user input parameters %v", err)
		respondWithError(w, http.StatusInternalServerError, "Something went wrong", err)
		return
	}
	hashedPassword, err := auth.HashPassword(inputParams.Password)
	if err != nil {
		log.Printf("Error hashing password %v", err)
		respondWithError(w, http.StatusInternalServerError, "Something went wrong", err)
		return
	}
	/*
			type UpdateUserParams struct {
			Email          string
			HashedPassword string
			ID             uuid.UUID
		}
	*/
	updateParams := database.UpdateUserParams{
		Email:          inputParams.Email,
		HashedPassword: hashedPassword,
		ID:             userJWT,
	}
	updatedUser, err := cfg.db.UpdateUser(context.Background(), updateParams)
	if err != nil {
		log.Printf("Error updating user record %v", err)
		respondWithError(w, http.StatusInternalServerError, "Something went wrong", err)
		return
	}
	respBody := userJsonStruct{
		ID:          updatedUser.ID,
		CreatedAt:   updatedUser.CreatedAt,
		UpdatedAt:   updatedUser.UpdatedAt,
		Email:       updatedUser.Email,
		IsChirpyRed: updatedUser.IsChirpyRed,
	}
	respondWithJSON(w, 200, respBody)

}

func (cfg *apiConfig) handlerPolkaWebhook(w http.ResponseWriter, r *http.Request) {
	/*
			request
			{
		  "event": "user.upgraded",
		  "data": {
		    "user_id": "3311741c-680c-4546-99f3-fc9efac2036c"
		  }
		}
	*/
	type InputPolkaWebhook struct {
		Event string `json:"event"`
		Data  struct {
			UserID string `json:"user_id"`
		} `json:"data"`
	}
	//fmt.Printf("request header in handler %v\n", r)
	fmt.Println()
	// check if the API key is correct
	apiKey, err := auth.GetAPIKey(r.Header)
	if err != nil {
		log.Printf("Error getting APIKey %v", err)
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	//fmt.Printf("incoming apikey and stored apikey __%s__, __%s__\n", apiKey, cfg.polkaKey)
	if apiKey != cfg.polkaKey {
		log.Printf("APIKey does not match %v", err)
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	decoder := json.NewDecoder(r.Body)
	params := InputPolkaWebhook{}
	err = decoder.Decode(&params)
	if err != nil {
		log.Printf("Error decoding request %v", err)
		respondWithError(w, http.StatusInternalServerError, "Something went wrong", err)
		return
	}
	if params.Event != "user.upgraded" {
		log.Printf("Incorrect Event received ")
		respondWithError(w, http.StatusNoContent, "", nil)
		return
	}

	userIDString := params.Data.UserID
	if userIDString == "" {
		log.Printf("No user ID found")
		respondWithError(w, http.StatusBadRequest, "No user id found", nil)
	}
	// convert incoming string ID to a uuid.UUID value
	userID, err := uuid.Parse(userIDString)
	if err != nil {
		log.Printf("Error converting user ID %v", err)
		respondWithError(w, http.StatusInternalServerError, "Something went wrong", nil)
	}
	_, err = cfg.db.SetRed(context.Background(), userID)
	if err != nil {
		log.Printf("Error updating user database to add is_chirpy_red to true %v", err)
		respondWithError(w, http.StatusNotFound, "User not found", nil)
	}
	// temporary check for user if needed otherwise use
	//_, err = cfg.db.SetRed(context.Background(), userID)
	//user, err := cfg.db.SetRed(context.Background(), userID)
	//fmt.Printf("user %+v", user)
	w.WriteHeader(http.StatusNoContent)
}
