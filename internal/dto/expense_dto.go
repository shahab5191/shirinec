package dto

import "shirinec.com/internal/models"

type ListExpenseCategoriesResponse struct {
	Pagination PaginationData           `json:"pagination"`
	Categories []models.ExpenseCategory `json:"categories"`
}
