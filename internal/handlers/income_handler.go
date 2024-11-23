package handler

import "github.com/gin-gonic/gin"

type IncomeHandler struct {
    incomeService services.IncomeService
}

func NewIncomeHandler(incomeService services.IncomeService) *IncomeHandler {
    return &IncomeHandler{incomeService: incomeService}
}

func (h *IncomeHandler) Categories(c *gin.Context) {

}
