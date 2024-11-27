package dto

type UpdatePasswordRequest struct {
	NewPassword     string `json:"newPassword"`
	CurrentPassword string `json:"currentPassword"`
}

type UpdateEmailRequest struct {
	NewEmail        string `json:"newEmail"`
	CurrentPassword string `json:"currentPassword"`
}

type VerificationRequest struct {
	VerificationCode int `json:"VerificationCode"`
}
