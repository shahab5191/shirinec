package dto

import "shirinec.com/internal/models"

type ItemCreateRequest struct {
	Name       string `json:"name"`
	ImageID    *int   `json:"imageID"`
	CategoryID int    `json:"categoryID"`
}

type ItemsListResponse struct {
	Pagination PaginationData `json:"pagination"`
	Items      *[]models.Item  `json:"items"`
}
