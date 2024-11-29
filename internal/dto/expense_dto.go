package dto

import "shirinec.com/internal/models"

type ExpenseCategoriesListResponse struct {
	Pagination PaginationData    `json:"pagination"`
	Categories []models.Category `json:"categories"`
}
