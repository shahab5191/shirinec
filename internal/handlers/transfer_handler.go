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

type TransferHandler interface {
	Transfer(c *gin.Context)
}

type transferHandler struct {
	transferService services.TransferService
}

func NewTransferHandler(transferService services.TransferService) TransferHandler {
	return &transferHandler{
		transferService: transferService,
	}
}

func (h *transferHandler) Transfer(c *gin.Context) {
	var input dto.TransferRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		if errList := server_errors.AsValidatorError(err); errList != nil {
			c.JSON(server_errors.ValidationErrorBuilder(errList).Unwrap())
			return
		}

		utils.Logger.Errorf("transferHandler.Transfer - Binding user input to dto.TransferRequest: %s", err.Error())
		c.JSON(server_errors.InternalError.Unwrap())
		return
	}

	userID, err := uuid.Parse(c.GetString("user_id"))
	if err != nil {
		utils.Logger.Errorf("transferHandler.Transfer - Parsing uuid from user_id string: %s", err.Error())
		c.JSON(server_errors.InvalidAuthorizationHeader.Unwrap())
		return
	}

	transferResult, err := h.transferService.Transfer(context.Background(), &input, userID)
	if err != nil {
		c.JSON(err.(*server_errors.SError).Unwrap())
		return
	}

	c.JSON(
		http.StatusOK,
		dto.CreateResponse[dto.AccountTransferResult]{
			Result: *transferResult,
		},
	)
}
