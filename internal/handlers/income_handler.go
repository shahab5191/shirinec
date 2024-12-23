package handler

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"shirinec.com/internal/dto"
	server_errors "shirinec.com/internal/errors"
	"shirinec.com/internal/services"
	"shirinec.com/internal/utils"
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
	var input dto.IncomeCreateRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		if errList := server_errors.AsValidatorError(err); errList != nil {
			c.JSON(server_errors.ValidationErrorBuilder(errList).Unwrap())
			return
		}

		utils.Logger.Errorf("incomeHandler.Create - Binding user input to dto.IncomeCreateRequest: %s", err.Error())
		c.JSON(server_errors.InternalError.Unwrap())
		return
	}

	userID, err := uuid.Parse(c.GetString("user_id"))
	if err != nil {
		utils.Logger.Errorf("incomeHandler.Create - Parsing uuid from user_id string %s", err.Error())
		c.JSON(server_errors.InternalError.Unwrap())
		return
	}

	response, err := h.incomeService.Create(context.Background(), &input, userID)
	if err != nil {
		c.JSON(err.(*server_errors.SError).Unwrap())
		return
	}

	c.JSON(
		http.StatusOK,
		dto.CreateResponse[dto.IncomeJoinedResponse]{
			Result: *response,
		},
	)
}
