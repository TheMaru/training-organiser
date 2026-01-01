package api

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/TheMaru/training-organiser/internal/auth"
	"github.com/TheMaru/training-organiser/internal/database"
	"github.com/google/uuid"
)

type UserResponse struct {
	ID           uuid.UUID `json:"id"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
	UserName     string    `json:"user_name"`
	Token        string    `json:"token"`
	RefreshToken string    `json:"refresh_token"`
}

func (cfg *ApiConfig) HandleRegisterUser(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		UserName string `json:"user_name"`
		Password string `json:"password"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		http.Error(w, "Error parsing JSON", http.StatusBadRequest)
		return
	}

	hashedPw, err := auth.HashPassword(params.Password)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error hashing password", err)
		return
	}

	dbUser, err := cfg.DB.CreateUser(r.Context(), database.CreateUserParams{
		UserName:     params.UserName,
		PasswordHash: hashedPw,
	})
	if err != nil {

		respondWithError(w, http.StatusInternalServerError, "Couldn't create user", err)
		return
	}

	user := UserResponse{
		ID:        dbUser.ID,
		CreatedAt: dbUser.CreatedAt.Time,
		UpdatedAt: dbUser.UpdatedAt.Time,
		UserName:  dbUser.UserName,
	}

	respondWithJSON(w, http.StatusCreated, user)
}
