package api

import (
	"encoding/json"
	"net/http"
	"os"
	"time"

	"github.com/TheMaru/training-organiser/internal/auth"
	"github.com/TheMaru/training-organiser/internal/database"
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
	token, err := auth.MakeJWT(dbUser.ID, secret, auth.JWT_DURATION)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't create session", err)
		return
	}

	refreshToken, err := auth.MakeRefreshToken()
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't create refresh token", err)
		return
	}

	_, err = cfg.DB.CreateRefreshToken(r.Context(), database.CreateRefreshTokenParams{
		Token:     refreshToken,
		UserID:    dbUser.ID,
		ExpiresAt: time.Now().Add(time.Hour * 24 * 60), // 60 days
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't save refresh token", err)
		return
	}

	user := UserResponse{
		ID:           dbUser.ID,
		CreatedAt:    dbUser.CreatedAt.Time,
		UpdatedAt:    dbUser.UpdatedAt.Time,
		UserName:     dbUser.UserName,
		Token:        token,
		RefreshToken: refreshToken,
	}
	respondWithJSON(w, http.StatusOK, user)
}
