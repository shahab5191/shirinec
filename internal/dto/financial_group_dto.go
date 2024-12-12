package dto

type FinancialGroupCreateRequest struct {
    Name string `json:"name" binding:"required,alphaNumericSpace"`
    ImageID *int `json:"imageID" binding:"omitempty,number"`
}
