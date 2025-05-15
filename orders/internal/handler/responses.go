package handler

import (
	"github.com/google/uuid"
)

type RefreshRequest struct {
	RefreshToken string `json:"refresh_token"`
}

type UserResponse struct {
	ID          uuid.UUID `json:"id"`
	PhoneNumber string    `json:"phone_number"`
	Name        string    `json:"name"`
	Mail        string    `json:"mail"`
	Birthday    string    `json:"birthday" example:"2006-01-02"`
	CreatedAt   string    `json:"created_at"`
}

type TokensResponse struct {
	Access  string `json:"access"`
	Refresh string `json:"refresh"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}

type AuthResponse struct {
	Tokens TokensResponse `json:"tokens"`
	User   UserResponse   `json:"user"`
}

type LogoutResponse struct {
	Message string `json:"message"`
}
