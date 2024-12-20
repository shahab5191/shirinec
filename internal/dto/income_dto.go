package dto

import (
	"time"

	"github.com/google/uuid"
)

type IncomeCreateRequest struct {
	AccountID   int     `json:"accountID" binding:"required,numeric"`
	Amount      float64 `json:"amount" binding:"required,numeric"`
	Description *string `json:"description" binding:"omitempty"`
}

type IncomeJoinedResponse struct {
	ID             int       `json:"id"`
	UserID         uuid.UUID `json:"userID"`
	AccountID      int       `json:"accountID"`
	AccountBalance float64   `json:"accountBalance"`
	Amount         float64   `json:"amount"`
	Description    string    `json:"description"`
	UpdateDate     time.Time `json:"updateDate"`
	CreationDate   time.Time `json:"creationDate"`
}
