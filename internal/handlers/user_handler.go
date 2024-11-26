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
	NewEmail(c *gin.Context)
    NewEmailVerification(c *gin.Context)
    SignupVerification(c *gin.Context)
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

func (h *userHandler) NewEmail(c *gin.Context) {
	var input dto.UpdateEmailRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(server_errors.InvalidInput.Unwrap())
		return
	}

	userID, err := uuid.Parse(c.GetString("user_id"))
	if err != nil {
		c.JSON(server_errors.InternalError.Unwrap())
	}

	verificationCode, err := h.userService.NewEmail(context.Background(), input, userID)
	if err != nil {
		log.Printf("[Error] - userHandler.NewEmail - Calling userService.NewEmail: %+v\n", err)
		var sErr *server_errors.SError
		if errors.As(err, &sErr) {
			c.JSON(sErr.Unwrap())
		} else {
			c.JSON(server_errors.InternalError.Unwrap())
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": verificationCode})
}

func (h *userHandler) NewEmailVerification(c *gin.Context) {
	var input dto.VerificationRequest
	err := c.ShouldBindJSON(&input)
	if err != nil {
        log.Printf("Verification error: %+v\n", err)
		c.JSON(server_errors.InvalidInput.Unwrap())
		return
	}

	userID, err := uuid.Parse(c.GetString("user_id"))
	if err != nil {
		c.JSON(server_errors.InternalError.Unwrap())
		return
	}

	err = h.userService.NewEmailVerification(context.Background(), input.VerificationCode, userID)
	if err != nil {
		var sErr *server_errors.SError
		if errors.As(err, &sErr) {
			c.JSON(sErr.Unwrap())
			return
		}
		c.JSON(server_errors.InternalError.Unwrap())
		return
	}

	c.JSON(http.StatusOK, gin.H{"result": "Email changed successfully"})
	return
}


func (h *userHandler) SignupVerification(c *gin.Context) {
    var input dto.VerificationRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(server_errors.InvalidInput.Unwrap())
		return
	}

	userID, err := uuid.Parse(c.GetString("user_id"))
	if err != nil {
		c.JSON(server_errors.InternalError.Unwrap())
		return
	}

    err = h.userService.SignupVerification(context.Background(), input.VerificationCode, userID)
    if err != nil {
        var sErr *server_errors.SError
        if errors.As(err, &sErr){
            c.JSON(sErr.Unwrap())
            return
        }
        c.JSON(server_errors.InternalError.Unwrap())
        return
    }

    c.JSON(http.StatusOK, gin.H{"result": "User verified successfully"})
    return
}
