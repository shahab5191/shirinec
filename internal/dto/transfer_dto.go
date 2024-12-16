package dto

type TransferRequest struct {
	From   int     `json:"from" binding:"required,number"`
	Dest   int     `json:"dest" binding:"required,number"`
	Type   string  `json:"type" binding:"required,accountType"`
	Amount float64 `json:"amount" binding:"required,number"`
}
