package dto

import (
	"shirinec.com/internal/enums"
	"shirinec.com/internal/models"
)

type CategoriesListResponse struct {
	Pagination PaginationData    `json:"pagination"`
	Categories []models.Category `json:"categories"`
}

type CategoryCreateRequest struct {
	Name   string             `json:"name,required" binding:"required,alphaNumericSpace"`
	Color  string             `json:"color,required" binding:"required,hexcolor"`
	IconID *int               `json:"iconID" binding:"omitempty,number"`
	Type   enums.CategoryType `json:"type,required" binding:"required,categoryCreateType"`
}

type CategoryUpdateRequest struct {
	Name   *string `json:"name" binding:"omitempty,alphaNumericSpace"`
	Color  *string `json:"color" binding:"omitempty,hexcolor"`
	IconID *int    `json:"iconID" binding:"omitempty,number"`
}
