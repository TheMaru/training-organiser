package api

import (
	"time"

	"github.com/google/uuid"
)

type LoginRequest struct {
	UserName string `json:"user_name"`
	Password string `json:"password"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}

type TokenResponse struct {
	Token string `json:"token"`
}

type RegisterUserRequest struct {
	UserName string `json:"user_name"`
	Password string `json:"password"`
}

type UpdateUserRequest struct {
	UserName string `json:"user_name"`
}

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
