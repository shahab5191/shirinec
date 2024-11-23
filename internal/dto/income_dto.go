package dto

import (
	"shirinec.com/internal/models"
)

type ListIncomeCategoreisRequest struct {
	Page int `form:"page,default=0"`
	Size int `form:"size,default=10"`
}

type ListIncomeCategoriesResponse struct {
	Pagination PaginationData          `json:"pagination"`
	Categories []models.IncomeCategory `json:"categories"`
}

type CreateIncomeCategoryRequest struct {
	Name  string `json:"name"`
	Color string `json:"color"`
}
