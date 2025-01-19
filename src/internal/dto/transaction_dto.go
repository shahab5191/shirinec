package dto

type TransferRequest struct {
	From   int     `json:"from" binding:"required,number"`
	Dest   int     `json:"dest" binding:"required,number"`
	Amount float64 `json:"amount" binding:"required,number"`
}
