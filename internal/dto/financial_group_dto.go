package dto

import (
	"time"

	"github.com/google/uuid"
	"shirinec.com/internal/enums"
)

type FinancialGroupCreateRequest struct {
	Name    string `json:"name" binding:"required,alphaNumericSpace"`
	ImageID *int   `json:"imageID" binding:"omitempty,number"`
}

type FinancialGroupAddUser struct {
	UserID uuid.UUID `json:"userID" binding:"required,uuid4"`
}

type FinancialGroup struct {
	ID           int
	Name         string
	ImageURL     *string
	OwnerID      uuid.UUID
	Users        []*UserGetResponse
	UserRole     enums.FinancialGroupRole
	CreationDate time.Time
	UpdateDate   time.Time
}

type FinancialGroupListItem struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type FinancialGroupListResponse struct {
	Pagination      PaginationData           `json:"pagination"`
	FinancialGroups []FinancialGroupListItem `json:"financialGroups"`
}
