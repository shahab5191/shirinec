package dto

import (
	"time"

	"github.com/google/uuid"
	"shirinec.com/internal/enums"
)

type AccountJoinedResponse struct {
	ID              int               `json:"id"`
	UserID          uuid.UUID         `json:"userID"`
	Name            string            `json:"name"`
	CategoryID      int               `json:"categoryID"`
	CategoryName    string            `json:"categoryName"`
	CategoryColor   string            `json:"categoryColor"`
	CategoryIconURL *string           `json:"categoryIconURL"`
	Balance         float64           `json:"balance"`
	CreationDate    time.Time         `json:"creationDate"`
	UpdateDate      time.Time         `json:"updateDate"`
	Type            enums.AccountType `json:"accountType"`
}

type AccountCreateRequest struct {
	Name       string            `json:"name" binding:"required,alphaNumericSpace"`
	CategoryID int               `json:"categoryID" binding:"required,number"`
	Balance    float64           `json:"balance" binding:"required,number"`
	Type       enums.AccountType `json:"accountType" binding:"omitempty,accountType"`
}

type AccountListResponse struct {
	Pagination PaginationData           `json:"pagination"`
	Accounts   *[]AccountJoinedResponse `json:"accounts"`
}

type AccountUpdateRequest struct {
	Name       *string  `json:"name" binding:"omitempty,alphaNumericSpace"`
	CategoryID *int     `json:"categoryID" binding:"omitempty,number"`
	Balance    *float64 `json:"balance" binding:"omitempty,number"`
}

type AccountTransferResultItem struct {
	ID      int               `json:"id"`
	Name    string            `json:"name"`
	Balance float64           `json:"balance"`
	Change  float64           `json:"change"`
	Type    enums.AccountType `json:"accountType"`
}

type AccountTransferResult struct {
	From AccountTransferResultItem `json:"from"`
	Dest AccountTransferResultItem `json:"dest"`
	Date time.Time                 `json:"date"`
}
