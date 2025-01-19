package dto

import (
	"time"

	"github.com/google/uuid"
	"shirinec.com/src/internal/enums"
)

type ItemCreateRequest struct {
	Name       string `json:"name" binding:"required,alphaNumericSpace"`
	ImageID    *int   `json:"imageID" binding:"omitempty,number"`
	CategoryID int    `json:"categoryID" binding:"number"`
}

type ItemJoinedResponse struct {
	ID              int                `json:"id"`
	UserID          uuid.UUID          `json:"userID"`
	Name            string             `json:"name"`
	ImageID         *int               `json:"imageID"`
	ImageURL        *string            `json:"imageURL"`
	ImageMetadata   *string            `json:"imageMetadata"`
	CategoryID      int                `json:"categoryID"`
	CategoryName    string             `json:"categoryName"`
	CategoryIconURL *string            `json:"categoryIconURL"`
	CategoryType    enums.CategoryType `json:"categoryType"`
	CreationDate    time.Time          `json:"creationDate"`
	UpdateDate      time.Time          `json:"updateDate"`
}

type ItemsListResponse struct {
	Pagination PaginationData        `json:"pagination"`
	Items      *[]ItemJoinedResponse `json:"items"`
}

type ItemUpdateRequest struct {
    Name       *string `json:"name" binding:"omitempty,alphaNumericSpace"`
    CategoryID *int    `json:"categoryID" bining:"omitEmpty,number"`
    ImageID    *int    `json:"imageID" binding:"omitempty,number"`
}
