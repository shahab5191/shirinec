package handler

import (
	"github.com/gin-gonic/gin"
	"shirinec.com/internal/services"
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

func (h *transferHandler) Transfer(c *gin.Context){

}
