package dto

import "github.com/google/uuid"

type FinancialGroupCreateRequest struct {
    Name string `json:"name" binding:"required,alphaNumericSpace"`
    ImageID *int `json:"imageID" binding:"omitempty,number"`
}

type FinancialGroupAddUser struct {
    UserID uuid.UUID `json:"userID" binding:"required,uuid4"`
}