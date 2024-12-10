package handler

import (
	"github.com/gin-gonic/gin"
	"shirinec.com/internal/services"
)

type FinancialGroupHandler interface {
    Create(c *gin.Context)
}

type financialGroupHandler struct {
    financialGroupService services.FinancialGroupService
}

func NewFinancialGroupHandler(financialGroupService *services.FinancialGroupService) FinancialGroupHandler {
    return &financialGroupHandler {
        financialGroupService: *financialGroupService,
    }
}

func (h *financialGroupHandler) Create(c *gin.Context) {
}
