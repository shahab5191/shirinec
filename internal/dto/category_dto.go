package dto

import (
	"shirinec.com/internal/enums"
	"shirinec.com/internal/models"
)

type ListCategoriesResponse struct {
	Pagination PaginationData    `json:"pagination"`
	Categories []models.Category `json:"categories"`
}

type CreateCategoryRequest struct {
	Name   string             `json:"name,required"`
	Color  string             `json:"color,required"`
	IconID *int               `json:"iconID"`
	Type   enums.CategoryType `json:"type,required"`
}

type UpdateCategoryRequest struct {
	Name   *string `json:"name"`
	Color  *string `json:"color"`
	IconID *int    `json:"iconID"`
}
