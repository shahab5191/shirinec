package dto

import "github.com/google/uuid"

type AuthLoginRequest struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required,min=8"`
}

type AuthSignupRequest struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required,min=8"`
}

type AuthLoginResponse struct {
	ID           uuid.UUID `json:"id"`
	AccessToken  string    `json:"accessToken"`
	RefreshToken string    `json:"refreshToken"`
    VerificationCode int `json:"verificationCode"`
}

type AuthRefreshTokenRequest struct {
	RefreshToken string `json:"refreshToken"`
}
