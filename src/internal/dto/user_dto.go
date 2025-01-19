package dto

import "github.com/google/uuid"

type UserUpdatePasswordRequest struct {
	NewPassword     string `json:"newPassword" binding:"required,min=8"`
	CurrentPassword string `json:"currentPassword" binding:"required,min=8"`
}

type UserUpdateEmailRequest struct {
	NewEmail        string `json:"newEmail" binding:"required,email"`
	CurrentPassword string `json:"currentPassword" binding:"required,min=8"`
}

type UserVerificationRequest struct {
	VerificationCode int `json:"VerificationCode" binding:"required,number,intLen=6"`
}

type UserGetResponse struct {
	ID                uuid.UUID `json:"id"`
	Name              *string    `json:"name"`
	FamilyName        *string    `json:"familyName"`
	ProfilePictureURL *string    `json:"profilePictureURL"`
}
