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
	Name   string             `json:"name,required"`
	Color  string             `json:"color,required"`
	IconID *int               `json:"iconID"`
	Type   enums.CategoryType `json:"type,required"`
}

type CategoryUpdateRequest struct {
	Name   *string `json:"name"`
	Color  *string `json:"color"`
	IconID *int    `json:"iconID"`
}
