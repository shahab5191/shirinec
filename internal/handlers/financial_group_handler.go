package handler

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"shirinec.com/internal/dto"
	"shirinec.com/internal/errors"
	"shirinec.com/internal/services"
	"shirinec.com/internal/utils"
)

type FinancialGroupHandler interface {
	Create(c *gin.Context)
}

type financialGroupHandler struct {
	financialGroupService services.FinancialGroupService
}

func NewFinancialGroupHandler(financialGroupService *services.FinancialGroupService) FinancialGroupHandler {
	return &financialGroupHandler{
		financialGroupService: *financialGroupService,
	}
}

func (h *financialGroupHandler) Create(c *gin.Context) {
	var input dto.FinancialGroupCreateRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		if errList := server_errors.AsValidatorError(err); errList != nil {
			utils.Logger.Infof("Error is ValidatorError: %s", err.Error())
			c.JSON(server_errors.ValidationErrorBuilder(errList).Unwrap())
			return
		}

		utils.Logger.Errorf("financialGroupHandler.Create - Binding request body to dto.FinancialGroupCreateRequest: %s", err.Error())
		c.JSON(server_errors.InvalidAuthorizationHeader.Unwrap())
		return
	}

	uid, err := uuid.Parse(c.GetString("user_id"))
	if err != nil {
		utils.Logger.Errorf("Parsing uuid from user_id string: %s", err.Error())
		c.JSON(server_errors.InternalError.Unwrap())
		return
	}

	item, err := h.financialGroupService.Create(context.Background(), &input, uid)
	if err != nil {
		c.JSON(err.(*server_errors.SError).Unwrap())
		return
	}

	c.JSON(http.StatusOK, item)
}
