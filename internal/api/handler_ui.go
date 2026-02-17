package api

import (
	"net/http"

	"github.com/TheMaru/training-organiser/internal/views"
)

func (cfg *ApiConfig) HandleHome(w http.ResponseWriter, r *http.Request) {
	component := views.Home()
	err := component.Render(r.Context(), w)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error rendering template", err)
	}
}
