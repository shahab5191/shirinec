package dto

type UserUpdatePasswordRequest struct {
	NewPassword     string `json:"newPassword"`
	CurrentPassword string `json:"currentPassword"`
}

type UserUpdateEmailRequest struct {
	NewEmail        string `json:"newEmail"`
	CurrentPassword string `json:"currentPassword"`
}

type UserVerificationRequest struct {
	VerificationCode int `json:"VerificationCode"`
}
