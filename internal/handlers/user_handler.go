package handler

import (
	"context"
	"errors"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"shirinec.com/internal/dto"
	"shirinec.com/internal/errors"
	"shirinec.com/internal/services"
)

type UserHandler interface {
	NewPassword(c *gin.Context)
}

type userHandler struct {
	userService services.UserService
}

func NewUserHandler(userService services.UserService) UserHandler {
	return &userHandler{userService: userService}
}

func (h *userHandler) NewPassword(c *gin.Context) {
	var input dto.UpdatePasswordRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(server_errors.InvalidInput.Unwrap())
		return
	}

	userID, err := uuid.Parse(c.GetString("user_id"))
	if err != nil {
		c.JSON(server_errors.InternalError.Unwrap())
		return
	}

	err = h.userService.NewPassword(context.Background(), input, userID)
	if err != nil {
		log.Printf("[Error] - userHandler.NewPassword - calling userService.NewPassword: %+v\n", err)
		var sErr *server_errors.SError
		if errors.As(err, &sErr) {
			c.JSON(sErr.Unwrap())
		} else {
			c.JSON(server_errors.InternalError.Unwrap())
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"result": "Successfully updated the password"})
}
