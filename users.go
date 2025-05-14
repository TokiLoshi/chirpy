package main

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strings"
	"time"

	auth "github.com/TokiLoshi/chirpy/internal"
	"github.com/TokiLoshi/chirpy/internal/database"
	"github.com/google/uuid"
)

type User struct {
	ID uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Email string `json:"email"`
	Token string `json:"token"`
	RefreshToken string `json:"refresh_token"`
	IsChirpyRed bool `json:"is_chirpy_red"`
}

type ReturnTokens struct {
	ID uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Email string `json:"email"`
	Token string `json:"token"`
	RefreshToken string `json:"refresh_token"`
		IsChirpyRed bool `json:"is_chirpy_red"`
}

func (cfg *apiConfig) validateUser(w http.ResponseWriter, req *http.Request) {
	type parameters struct {
		Email string `json:"email"`
		Password string `json:"password"`
	}

	decoder := json.NewDecoder(req.Body)
	params := parameters{}

	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, 400, "error decoding json data")
		return
	}

	userEmail := params.Email
	userPassword := params.Password 

	hashedPassword, err := auth.HashPassword(userPassword)

	if err != nil {
		respondWithError(w, 500, "unable to hash password")
	}

	dbUser, err := cfg.dbQueries.CreateUser(req.Context(), database.CreateUserParams{
		Email: userEmail, 
		HashedPassword: hashedPassword,
	})

	if err != nil {
		respondWithError(w, 400, "couldn't create user: " + err.Error() )
		return
		} 
	
	user := User {
		ID: dbUser.ID,
		Email: dbUser.Email,
		CreatedAt: dbUser.CreatedAt,
		UpdatedAt: dbUser.UpdatedAt,
		IsChirpyRed: dbUser.IsChirpyRed.Bool,
	}

	respondWithJson(w, 201, user)

}



func (cfg *apiConfig) handleLogin(w http.ResponseWriter, req *http.Request) {
	type parameters struct {
		Email string `json:"email"`
		Password string `json:"password"`
	}
	decoder := json.NewDecoder(req.Body)
	params := parameters{}

	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, 400, "error decoding json data")
		return
	}

	userEmail := params.Email 
	userPassword := params.Password 

	userEntry, err := cfg.dbQueries.GetUserByEmail(req.Context(), userEmail)

	if err != nil {
		respondWithError(w, 401, "Unauthorized")
		return
	}

	error := auth.CheckPasswordHash(userEntry.HashedPassword, userPassword)

	if error != nil {
		respondWithError(w, 401, "Unauthorized")
		return
	}

	expirytDuration := time.Hour

	jwtToken, err := auth.MakeJWT(userEntry.ID, cfg.jwtSecret, expirytDuration)

	if err != nil {
		respondWithError(w, 500, "Unauthorized")
		return
	}

	refreshToken, err := auth.MakeRefreshToken() 
	if err != nil {
		respondWithError(w, 500, "Unauthorized")
		return
	}

	refreshEntry, err := cfg.dbQueries.CreateRefreshToken(req.Context(), database.CreateRefreshTokenParams{
		Token: refreshToken,
		UserID: userEntry.ID,
	})
	if err != nil {
		respondWithError(w, 500, "Could not save refreshToken")
		return
	}

	tokens :=  ReturnTokens {
		ID : userEntry.ID,
		CreatedAt: refreshEntry.CreatedAt,
		UpdatedAt: refreshEntry.UpdatedAt,
		Email: userEntry.Email,
		Token: jwtToken,
		RefreshToken: refreshEntry.Token,
		IsChirpyRed: userEntry.IsChirpyRed.Bool,
	}

	respondWithJson(w, 200, tokens)

}

type returnToken struct {
	token string;
}

func (cfg *apiConfig) hanldeRefresh(w http.ResponseWriter, req *http.Request) {
	// Look up in the database if the token exists 
	authHeader := req.Header.Get("Authorization")
	if authHeader == "" {
		respondWithError(w, 401, "Auth header missing")
		return
	}

	// Format for bearer 
	fields := strings.Split(authHeader, " ")
	if len(fields) != 2 || fields[0] != "Bearer" {
		respondWithError(w, 401, "Bearer missing from Authorization")
		return
	}

	token := fields[1]
	
	refreshToken, err := cfg.dbQueries.GetRefreshToken(req.Context(), token)
	if err != nil {
		respondWithError(w, 401, "Invalid RefreshToken")
		return
	}

	if refreshToken.RevokedAt.Valid {
		respondWithError(w, 401, "Token has been revoked")
		return 
	}

	if time.Now().After(refreshToken.ExpiresAt) {
		respondWithError(w, 401, "Token has expired")
		return
	}
	
	newJWT, err := auth.MakeJWT(refreshToken.UserID, cfg.jwtSecret, time.Hour)

	if err != nil {
		respondWithError(w, 500, "Error creating access token")
		return
	}

	response := struct {
		Token string `json:"token"`
	}{
		Token: newJWT,
	}

	respondWithJson(w, 200, response)
	
}

func (cfg *apiConfig) handleRevoke(w http.ResponseWriter, req *http.Request) {
	// require refresh token to be present 
	authHeader := req.Header.Get("Authorization")
	if authHeader == "" {
		respondWithError(w, 401, "Auth header missing")
		return
	}

	fields := strings.Split(authHeader, " ")
	if len(fields) != 2 || fields[0] != "Bearer" {
		respondWithError(w, 401, "Bearer missing from Authorization Header")
		return
	}

	token := fields[1]
	revokedAt := time.Now()
	_, err := cfg.dbQueries.RevokeToken(req.Context(), database.RevokeTokenParams{
		RevokedAt: sql.NullTime{Time: revokedAt, Valid: true},
		Token: token,
	})
	if err != nil {
		respondWithError(w, 500, "Error revoking token")
		return 
	}

	w.WriteHeader(204)

}

func (cfg *apiConfig) handleAuthentication(w http.ResponseWriter, req *http.Request) {
	type parameters struct {
		Email string `json:"email"`
		Password string `json:"password"`
	}

	// headers must have an access token 
	authHeader := req.Header.Get("Authorization")
	if authHeader == "" {
		respondWithError(w, 401, "Auth header is missing")
		return
	}

	// Format for Bearer 
	fields := strings.Split(authHeader, " ")
	if len(fields) != 2 || fields[0] != "Bearer" {
		respondWithError(w, 401, "Malformed Authorization")
		return
	}

	tokenString := fields[1]

	userId, err := auth.ValidateJWT(tokenString, cfg.jwtSecret)
	if err != nil {
		respondWithError(w, 401, "Invalid token")
		return
	}

	decoder := json.NewDecoder(req.Body)
	params := parameters{}

	err = decoder.Decode(&params)
	if err != nil {
		respondWithError(w, 500, "error decoding json data")
		return
	}

	newEmail := params.Email 
	newPassword := params.Password 

	// users should be able to update their own email and password
	// headers must have a new password and email 

	if newEmail == "" || newPassword == "" {
		respondWithError(w, 401, "Unauthorized")
		return
	}

	// hash the password 
	hashedPassword, err := auth.HashPassword(newPassword)
	if err != nil {
		respondWithError(w, 500, "unable to hash password")
		return
	}


	dbUser, err := cfg.dbQueries.UpdateUsers(req.Context(), database.UpdateUsersParams{
		Email: newEmail,
		HashedPassword: hashedPassword,
		UpdatedAt: time.Now(),
		ID: userId,
	})

	if err != nil {
		respondWithError(w, 500, "Error updating user password and email")
		return
	}

	updatedUser := User {
		ID: dbUser.ID, 
		Email: dbUser.Email,
		CreatedAt: dbUser.CreatedAt, 
		UpdatedAt: dbUser.UpdatedAt,
		IsChirpyRed: dbUser.IsChirpyRed.Bool,
	}

	respondWithJson(w, 200, updatedUser)

}

func (cfg *apiConfig) handleUpgrade(w http.ResponseWriter, req *http.Request) {
	type parameters struct {
		Event string `json:"event"`
		Data struct {
			UserId string `json:"user_id"`
		}
	}

	requestApiKey, err := auth.GetAPIkeys(req.Header)
	if err != nil {
		respondWithError(w, 401, "Unauthorized")
		return
	}

	secretApiKey := cfg.polkaKey
	if requestApiKey != secretApiKey {
		respondWithError(w, 401, "Unauthorized")
		return
	}


	// decode data 
	decoder := json.NewDecoder(req.Body)
	params := parameters{}
	err = decoder.Decode(&params)
	if err != nil {
		respondWithError(w, 400, "error decoding json data")
		return 
	}

	eventType := params.Event 
	userId := params.Data.UserId 

	// if the event is anything other than user.ugraded return 204 
	if eventType != "user.upgraded" {
		respondWithError(w, 204, "incorrect event type")
	}

	// parse userid 
	id, err := uuid.Parse(userId)
	if err != nil {
		respondWithError(w, 500, "error parsing userId")
	}
	
	// upgrade user in database 
	dbUser, err := cfg.dbQueries.UpgradeUsers(req.Context(), database.UpgradeUsersParams{
		ID: id, 
		UpdatedAt: time.Now(),
	})

	if err != nil {
		respondWithError(w, 404, "user could not be found")
	}

	upgradedUser := User {
		ID: dbUser.ID, 
		Email: dbUser.Email,
		CreatedAt: dbUser.CreatedAt,
		UpdatedAt: dbUser.UpdatedAt,
		IsChirpyRed: dbUser.IsChirpyRed.Bool,
	}

	respondWithJson(w, 204, upgradedUser)


	// if the user is upgraded respond with 204 and empty body
	// if user can't be found respond with 404 
	// update all the user resources to include is_chirpy_red field 
}