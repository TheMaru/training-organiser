package api

import (
	"encoding/json"
	"net/http"
	"os"
	"time"

	"github.com/TheMaru/training-organiser/internal/auth"
)

func (cfg *ApiConfig) HandleLogin(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		UserName string `json:"user_name"`
		Password string `json:"password"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	if err := decoder.Decode(&params); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request", err)
		return
	}

	dbUser, err := cfg.DB.GetUserByUserName(r.Context(), params.UserName)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Invalid username or password", err)
		return
	}

	passwordMatch, err := auth.CheckPassword(params.Password, dbUser.PasswordHash)
	if err != nil || !passwordMatch {
		respondWithError(w, http.StatusInternalServerError, "Invalid username or password", err)
		return
	}

	secret := os.Getenv("JWT_SECRET")
	expirationDuration := time.Duration(1) * time.Hour
	token, err := auth.MakeJWT(dbUser.ID, secret, expirationDuration)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't create session", err)
		return
	}

	// TODO: Add refresh token at some point

	user := UserResponse{
		ID:        dbUser.ID,
		CreatedAt: dbUser.CreatedAt.Time,
		UpdatedAt: dbUser.UpdatedAt.Time,
		UserName:  dbUser.UserName,
		Token:     token,
	}
	respondWithJSON(w, http.StatusOK, user)
}
