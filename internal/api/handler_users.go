package api

import (
	"encoding/json"
	"net/http"

	"github.com/TheMaru/training-organiser/internal/auth"
	"github.com/TheMaru/training-organiser/internal/database"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

func dbUserToPublicUser(user database.User) UserPublicResponse {
	return UserPublicResponse{
		ID:           user.ID,
		CreatedAt:    user.CreatedAt.Time,
		UpdatedAt:    user.UpdatedAt.Time,
		UserName:     user.UserName,
		PlatformRole: user.PlatformRole,
	}
}

// @Summary Register a new user
// @Description Registers a new user and returns said user
// @Tags users
// @Accept json
// @Produce json
// @Param request body RegisterUserRequest true "Register data"
// @Success 201 {object} UserResponse
// @Failure 500 {object} ErrorResponse
// @Router /users [post]
func (cfg *ApiConfig) HandleRegisterUser(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	params := RegisterUserRequest{}
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

// @Summary Return all users
// @Description Returns all current registered users
// @Tags users
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param Authorization header string true "Bearer token"
// @Success 200 {array} UserPublicResponse
// @Failure 401 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /users [get]
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

// @Summary Return specific user
// @Description Returns a specific user
// @Tags users
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param Authorization header string true "Bearer token"
// @Param id path string true "User UUID"
// @Success 200 {object} UserPublicResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Router /users/{id} [get]
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

// @Summary Return logged in user
// @Description Returns the logged in user that sent the request
// @Tags users
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param Authorization header string true "Bearer token"
// @Success 200 {object} UserPublicResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Router /users/me [get]
func (cfg *ApiConfig) HandleGetMyProfile(w http.ResponseWriter, r *http.Request) {
	userID, err := auth.GetUserID(r.Context())
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Unable to extract user id from token", err)
		return
	}

	dbUser, err := cfg.DB.GetUserByID(r.Context(), userID)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "User not found", err)
		return
	}
	respondWithJSON(w, http.StatusOK, dbUserToPublicUser(dbUser))
}

// @Summary Updates user
// @Description updates the username of a user
// @Tags users
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param Authorization header string true "Bearer token"
// @Param id path string true "User UUID"
// @Param request body UpdateUserRequest true "Update user data"
// @Success 200 {object} UserPublicResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /users/{id} [put]
func (cfg *ApiConfig) HandleUpdateUser(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		UserName string `json:"user_name"`
	}
	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Couldn't decode parameters", err)
		return
	}

	idParam := chi.URLParam(r, "id")
	userID, err := uuid.Parse(idParam)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid user ID", err)
		return
	}

	requesterID, err := auth.GetUserID(r.Context())
	if err != nil || userID != requesterID {
		respondWithError(w, http.StatusUnauthorized, "Unauthorized to change other users", err)
		return
	}

	dbUser, err := cfg.DB.UpdateUser(r.Context(), database.UpdateUserParams{
		ID:       userID,
		UserName: params.UserName,
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Could not update user", err)
		return
	}

	respondWithJSON(w, http.StatusOK, dbUserToPublicUser(dbUser))
}

// @Summary Deletes a user
// @Description Deletes user
// @Tags users
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param Authorization header string true "Bearer token"
// @Param id path string true "User UUID"
// @Success 204 {string} string "No Content"
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /users/{id} [delete]
func (cfg *ApiConfig) HandleDeleteUser(w http.ResponseWriter, r *http.Request) {
	idParam := chi.URLParam(r, "id")
	userID, err := uuid.Parse(idParam)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid user ID", err)
		return
	}

	requesterID, err := auth.GetUserID(r.Context())
	if err != nil || userID != requesterID {
		respondWithError(w, http.StatusUnauthorized, "Unauthorized to delete other users", err)
		return
	}

	err = cfg.DB.DeleteUser(r.Context(), userID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Could not delete user", err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
