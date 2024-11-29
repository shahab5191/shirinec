package dto

import (
	"time"

	"github.com/google/uuid"
	"shirinec.com/internal/enums"
)

type ItemCreateRequest struct {
	Name       string `json:"name"`
	ImageID    *int   `json:"imageID"`
	CategoryID int    `json:"categoryID"`
}

type ItemJoinedResponse struct {
	ID              int                `json:"id"`
	UserID          uuid.UUID          `json:"userID"`
	Name            string             `json:"name"`
	ImageID         *int                `json:"imageID"`
	ImageURL        *string             `json:"imageURL"`
	ImageMetadata   *string             `json:"imageMetadata"`
	CategoryID      int                `json:"categoryID"`
	CategoryName    string             `json:"categoryName"`
	CategoryIconURL *string             `json:"categoryIconURL"`
	CategoryType    enums.CategoryType `json:"categoryType"`
	CreationDate    time.Time          `json:"creationDate"`
	UpdateDate      time.Time          `json:"updateDate"`
}

type ItemsListResponse struct {
	Pagination PaginationData `json:"pagination"`
	Items      *[]ItemJoinedResponse `json:"items"`
}
