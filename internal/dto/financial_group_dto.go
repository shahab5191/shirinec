package dto

type FinancialGroupCreateRequest struct {
    Name string `json:"name" binding:"required,alphaNumSpace"`
    ImageID int `json:"imageID" binding:"omitempty,number"`
}
