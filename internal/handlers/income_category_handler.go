package handler

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"shirinec.com/internal/dto"
	"shirinec.com/internal/errors"
	"shirinec.com/internal/services"
)

type IncomeCategoryHandler interface {
    List(c *gin.Context)
}

type incomeCategoryHandler struct {
    incomeCategoryService services.IncomeCategoryService
}

func NewIncomeCategoryHandler(incomeCategoryService services.IncomeCategoryService) IncomeCategoryHandler {
    return &incomeCategoryHandler{incomeCategoryService: incomeCategoryService}
}

func (h *incomeCategoryHandler) List(c *gin.Context) {
    var input dto.ListIncomeCategoreisRequest
    if err := c.ShouldBindQuery(&input); err != nil {
        log.Printf("[Error] - incomeCategoryHandler.List - Bind query %+v\n", err)
        c.JSON(server_errors.InvalidInput.Unwrap())
        return
    }

    userID, err := uuid.Parse(c.GetString("user_id"))
    if err != nil {
        c.JSON(server_errors.InternalError.Unwrap())
        return
    }

    categories, err := h.incomeCategoryService.ListCategories(userID, input.Limit, input.Offset)
    if err != nil {
        c.JSON(server_errors.InternalError.Unwrap())
    }

    c.JSON(http.StatusOK, categories)
}
