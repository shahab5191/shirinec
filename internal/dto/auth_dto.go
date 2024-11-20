package dto

import "github.com/google/uuid"

type LoginRequest struct {
    Email       string  `json:"email" binding:"required"`
    Password    string  `json:"password" binding:"required,min=8"`
}

type CreateUserDTO struct {
    Email       string  `json:"email" binding:"required, email"`
    Password    string  `json:"password" binding:"required,min=8"`
}

type LoginResponse struct {
    ID              uuid.UUID   `json:"id"`
    Token           string      `json:"token"`
    RefreshToken   string      `json:"refresh_token"`
}
