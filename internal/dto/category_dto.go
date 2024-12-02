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
	Name   string             `json:"name,required" validate:"required,alphanum"`
	Color  string             `json:"color,required" validate:"required,hexcolor"`
	IconID *int               `json:"iconID" validate:"omitempty,number"`
	Type   enums.CategoryType `json:"type,required" validate:"required,categoryCreateType"`
}

type CategoryUpdateRequest struct {
    Name   *string `json:"name" validate:"omitempty,alphanum"`
    Color  *string `json:"color" validate:"omitempty,hexcolor"`
    IconID *int    `json:"iconID" validate:"omitempty,number"`
}
