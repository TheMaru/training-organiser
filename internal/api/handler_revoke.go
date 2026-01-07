package api

import (
	"net/http"

	"github.com/TheMaru/training-organiser/internal/auth"
)

// @Summary Invalidate refresh token
// @Description Revokes a refresh token to it can no longer be used
// @Tags auth
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer Refresh Token"
// @Success 204 {string} string "No Content"
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /revoke [post]
func (cfg *ApiConfig) HandleRevokeToken(w http.ResponseWriter, r *http.Request) {
	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid token format", err)
		return
	}

	err = cfg.DB.RevokeRefreshToken(r.Context(), token)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Could not revoke token", err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
