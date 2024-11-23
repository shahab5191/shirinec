package handler

import (
	"github.com/gin-gonic/gin"
	"shirinec.com/internal/dto"
	server_errors "shirinec.com/internal/errors"
	"shirinec.com/internal/services"
)

type IncomeCategoryHandler interface {
    List(c *gin.Context)
}

type incomeCategoryHandler struct {
    incomeCategoryService services.IncomeCategoryService
}

func NewIncomeCategoryService(incomeCategoryService services.IncomeCategoryService) IncomeCategoryHandler {
    return &incomeCategoryHandler{incomeCategoryService: incomeCategoryService}
}

func (h *incomeCategoryHandler) List(c *gin.Context) {
    var input dto.ListIncomeCategoreisRequest
    if err := c.ShouldBindJSON(input); err != nil {
        c.JSON(server_errors.InvalidInput.Unwrap())
        return
    }
}
