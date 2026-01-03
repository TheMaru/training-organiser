package api

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/TheMaru/training-organiser/internal/auth"
	"github.com/TheMaru/training-organiser/internal/database"
	"github.com/go-chi/chi/v5"
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

type UserPublicResponse struct {
	ID           uuid.UUID `json:"id"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
	UserName     string    `json:"user_name"`
	PlatformRole string    `json:"platform_role"`
}

func dbUserToPublicUser(user database.User) UserPublicResponse {
	return UserPublicResponse{
		ID:           user.ID,
		CreatedAt:    user.CreatedAt.Time,
		UpdatedAt:    user.UpdatedAt.Time,
		UserName:     user.UserName,
		PlatformRole: user.PlatformRole,
	}
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
		respondWithError(w, http.StatusInternalServerError, "Error parsing JSON", err)
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

func (cfg *ApiConfig) HandleListUsers(w http.ResponseWriter, r *http.Request) {
	dbUsers, err := cfg.DB.GetUsers(r.Context())
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Could not get users", err)
		return
	}

	usersResponse := make([]UserPublicResponse, 0)

	for _, dbUser := range dbUsers {
		usersResponse = append(usersResponse, dbUserToPublicUser(dbUser))
	}

	respondWithJSON(w, http.StatusOK, usersResponse)
}

func (cfg *ApiConfig) HandleListUser(w http.ResponseWriter, r *http.Request) {
	idParam := chi.URLParam(r, "id")
	userID, err := uuid.Parse(idParam)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid user ID", err)
		return
	}

	dbUser, err := cfg.DB.GetUserByID(r.Context(), userID)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "User not found", err)
		return
	}

	respondWithJSON(w, http.StatusOK, dbUserToPublicUser(dbUser))
}
