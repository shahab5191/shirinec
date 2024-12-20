package handler

import (
	"github.com/gin-gonic/gin"
	"shirinec.com/internal/services"
)

type IncomeHandler interface {
	Create(ctx *gin.Context)
}

type incomeHandler struct {
	incomeService services.IncomeService
}

func NewIncomeHandler(incomeService services.IncomeService) IncomeHandler {
	return &incomeHandler{
		incomeService: incomeService,
	}
}

func (h *incomeHandler) Create(c *gin.Context) {

}
