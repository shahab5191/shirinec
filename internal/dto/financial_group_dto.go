package dto

import (
	"time"

	"github.com/google/uuid"
)

type FinancialGroupCreateRequest struct {
    Name string `json:"name" binding:"required,alphaNumericSpace"`
    ImageID *int `json:"imageID" binding:"omitempty,number"`
}

type FinancialGroupAddUser struct {
    UserID uuid.UUID `json:"userID" binding:"required,uuid4"`
}

type FinancialGroup struct {
    ID int
    Name string
    ImageURL *string
    UserID uuid.UUID
    Users []*UserGetResponse
    CreationDate time.Time
    UpdateDate time.Time
}
