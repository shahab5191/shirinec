package dto

import (
	"time"

	"github.com/google/uuid"
)

type AccountJoinedResponse struct {
	ID              int       `json:"id"`
	UserID          uuid.UUID `json:"userID"`
	Name            string    `json:"name"`
	CategoryID      int       `json:"categoryID"`
	CategoryName    string    `json:"categoryName"`
	CategoryColor   string    `json:"categoryColor"`
	CategoryIconURL *string   `json:"categoryIconURL"`
	Balance         float64   `json:"balance"`
	CreationDate    time.Time `json:"creationDate"`
	UpdateDate      time.Time `json:"updateDate"`
}

type AccountCreateRequest struct {
	Name       string  `json:"name"`
	CategoryID int     `json:"categoryID"`
	Balance    float64 `json:"balance"`
}

type AccountListResponse struct {
	Pagination PaginationData           `json:"pagination"`
	Accounts   *[]AccountJoinedResponse `json:"accounts"`
}

type AccountUpdateRequest struct {
	Name       *string  `json:"name"`
	CategoryID *int     `json:"categoryID"`
	Balance    *float64 `json:"balance"`
}
