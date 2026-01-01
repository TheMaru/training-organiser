package api

import (
	"net/http"

	"github.com/TheMaru/training-organiser/internal/auth"
)

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
