package api

import (
	"net/http"
	"os"
	"time"

	"github.com/TheMaru/training-organiser/internal/auth"
)

func (cfg *ApiConfig) HandleRefreshToken(w http.ResponseWriter, r *http.Request) {
	refreshToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Couldn't find token", err)
		return
	}

	refreshTokenDB, err := cfg.DB.GetRefreshToken(r.Context(), refreshToken)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Invalid refresh token", err)
		return
	}

	if time.Now().After(refreshTokenDB.ExpiresAt) {
		respondWithError(w, http.StatusUnauthorized, "Token expired", nil)
		return
	}

	if refreshTokenDB.RevokedAt.Valid {
		respondWithError(w, http.StatusUnauthorized, "Token revoked", nil)
		return
	}

	newAccessToken, err := auth.MakeJWT(refreshTokenDB.UserID, os.Getenv("JWT_SECRET"), auth.JWT_DURATION)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't create new access token", err)
		return
	}

	respondWithJSON(w, http.StatusOK, struct {
		Token string `json:"token"`
	}{
		Token: newAccessToken,
	})
}
