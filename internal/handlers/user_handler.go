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
	var input dto.UserUpdatePasswordRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		if errList := server_errors.AsValidatorError(err); errList != nil {
			c.JSON(server_errors.ValidationErrorBuilder(errList).Unwrap())
			return
		}
        utils.Logger.Errorf("Binding input to dto.UserUpdatePasswordRequest: %s", err.Error())
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
		c.JSON(err.(*server_errors.SError).Unwrap())
		return
	}

	c.JSON(http.StatusOK, gin.H{"result": "Successfully updated the password"})
}

func (h *userHandler) NewEmail(c *gin.Context) {
	var input dto.UserUpdateEmailRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		if errList := server_errors.AsValidatorError(err); errList != nil {
			c.JSON(server_errors.ValidationErrorBuilder(errList).Unwrap())
			return
		}
		utils.Logger.Warnf("userHandler.NewEmail - Undefined error while binding input to dto.UserUpdateEmailRequest: %s", err.Error())
		c.JSON(server_errors.InvalidInput.Unwrap())
		return
	}

	userID, err := uuid.Parse(c.GetString("user_id"))
	if err != nil {
		c.JSON(server_errors.InternalError.Unwrap())
	}

	verificationCode, err := h.userService.NewEmail(context.Background(), input, userID)
	if err != nil {
		c.JSON(err.(*server_errors.SError).Unwrap())
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": verificationCode})
}

func (h *userHandler) NewEmailVerification(c *gin.Context) {
	var input dto.UserVerificationRequest
	err := c.ShouldBindJSON(&input)
	if err != nil {
		if errList := server_errors.AsValidatorError(err); errList != nil {
			c.JSON(server_errors.ValidationErrorBuilder(errList).Unwrap())
			return
		}
		utils.Logger.Errorf("userHandler.NewEmailVerification - Undefined error while binding input to dto.UserVerificationRequest: %s", err.Error())
		c.JSON(server_errors.InvalidInput.Unwrap())
		return
	}

	userID, err := uuid.Parse(c.GetString("user_id"))
	if err != nil {
        utils.Logger.Errorf("userHandler.NewEmailVerification - Parsing uuid from user_id string: %s", err.Error())
		c.JSON(server_errors.InternalError.Unwrap())
		return
	}

	err = h.userService.NewEmailVerification(context.Background(), input.VerificationCode, userID)
	if err != nil {
		c.JSON(err.(*server_errors.SError).Unwrap())
		return
	}

	c.JSON(http.StatusOK, gin.H{"result": "Email changed successfully"})
}

func (h *userHandler) SignupVerification(c *gin.Context) {
	var input dto.UserVerificationRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		if errList := server_errors.AsValidatorError(err); errList != nil {
			c.JSON(server_errors.ValidationErrorBuilder(errList).Unwrap())
			return
		}
		utils.Logger.Warnf("userHandler.SignupVerification - Undefined error while binding input to dto.UserVerificationRequest: %s", err.Error())
		c.JSON(server_errors.InvalidInput.Unwrap())
		return
	}

	userID, err := uuid.Parse(c.GetString("user_id"))
	if err != nil {
        utils.Logger.Errorf("[Error] - userHandler.SignupVerification - Parsing user_id to uuid: %s", err.Error())
		c.JSON(server_errors.InternalError.Unwrap())
		return
	}

	err = h.userService.SignupVerification(context.Background(), input.VerificationCode, userID)
	if err != nil {
		c.JSON(err.(*server_errors.SError).Unwrap())
		return
	}

	c.JSON(http.StatusOK, gin.H{"result": "User verified successfully"})
}
