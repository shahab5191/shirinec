package handler

import (
	"context"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"shirinec.com/src/internal/dto"
	"shirinec.com/src/internal/errors"
	"shirinec.com/src/internal/services"
	"shirinec.com/src/internal/utils"
)

type AuthHandler struct {
	authService services.AuthService
}

func NewAuthHandler(authService services.AuthService) *AuthHandler {
	return &AuthHandler{
		authService: authService,
	}
}

func (h *AuthHandler) SignUp(c *gin.Context) {
	var input dto.AuthSignupRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		if errList := server_errors.AsValidatorError(err); errList != nil {
			c.JSON(server_errors.ValidationErrorBuilder(errList).Unwrap())
			return
		}
		c.JSON(server_errors.InvalidInput.Unwrap())
		return
	}

	response, err := h.authService.CreateUser(context.Background(), &input, c.ClientIP())
	if err != nil {
		if errors.Is(err, &server_errors.UserAlreadyExistsError) {
			c.JSON(server_errors.UserAlreadyExistsError.Unwrap())
			return
		}
		utils.Logger.Errorf("Calling authService.CreateUser: %s", err.Error())
		c.JSON(server_errors.InternalError.Unwrap())
		return
	}

	c.JSON(http.StatusCreated, *response)
}

func (h *AuthHandler) Login(c *gin.Context) {
	var credentials dto.AuthLoginRequest
	if err := c.ShouldBindJSON(&credentials); err != nil {
		if errList := server_errors.AsValidatorError(err); errList != nil {
			c.JSON(server_errors.ValidationErrorBuilder(errList).Unwrap())
			return
		}
		utils.Logger.Errorf("binding input to dto.AuthLoginRequest: %s", err.Error())
		c.JSON(server_errors.CredentialError.Unwrap())
		return
	}

	loginResponse, err := h.authService.Login(credentials.Email, credentials.Password, c.ClientIP())
	if err != nil {
		c.JSON(err.(*server_errors.SError).Unwrap())
		return
	}

	c.JSON(http.StatusOK, loginResponse)
}

func (h *AuthHandler) RefreshToken(c *gin.Context) {
	var requestDTO dto.AuthRefreshTokenRequest
	if err := c.ShouldBindJSON(&requestDTO); err != nil {
		if errList := server_errors.AsValidatorError(err); errList != nil {
			c.JSON(server_errors.ValidationErrorBuilder(errList).Unwrap())
			return
		}
		utils.Logger.Errorf("binding request to dto: %s", err.Error())
		c.JSON(server_errors.InvalidInput.Unwrap())
		return
	}

	response, err := h.authService.Refresh(requestDTO.RefreshToken)
	if err != nil {
		c.JSON(err.(*server_errors.SError).Unwrap())
		return
	}
	c.JSON(http.StatusOK, response)
}
