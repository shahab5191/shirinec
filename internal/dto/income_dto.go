package dto

type ListIncomeCategoreisRequest struct {
    Limit   int     `form:"limit"`
    Offset  int     `form:"offset"`
}
