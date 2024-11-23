package dto

import "github.com/google/uuid"

type LoginRequest struct {
    Email       string  `json:"email" binding:"required"`
    Password    string  `json:"password" binding:"required,min=8"`
}

type CreateUserRequest struct {
    Email       string  `json:"email" binding:"required"`
    Password    string  `json:"password" binding:"required,min=8"`
}

type LoginResponse struct {
    ID              uuid.UUID   `json:"id"`
    AccessToken     string      `json:"accessToken"`
    RefreshToken    string      `json:"refreshToken"`
}

type RefreshTokenRequest struct {
    RefreshToken    string  `json:"refreshToken"`
}

type ListIncomeCategoreisRequest struct {
    Limit   int
    Offset  int
}
