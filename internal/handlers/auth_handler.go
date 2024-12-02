package handler

import (
	"context"
	"errors"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"shirinec.com/internal/dto"
	"shirinec.com/internal/errors"
	"shirinec.com/internal/services"
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
		log.Printf("[Error] - AuthHandler.SignUp - Calling authService.CreateUser: %+v", err)
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
		log.Printf("[Error] - AuthHandler.Login - binding input to dto.AuthLoginRequest: %s", err)
		c.JSON(server_errors.CredentialError.Unwrap())
		return
	}

	loginResponse, err := h.authService.Login(credentials.Email, credentials.Password, c.ClientIP())
	var se *server_errors.SError
	if err != nil {
		if errors.As(err, &se) {
			c.JSON(se.Unwrap())
			return
		}
		c.JSON(server_errors.InternalError.Unwrap())
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
		log.Printf("[Error] - AuthHandler.RefreshToken - binding request to dto: %+v\n", err)
		c.JSON(server_errors.InvalidInput.Unwrap())
		return
	}

	response, err := h.authService.Refresh(requestDTO.RefreshToken)
	if err != nil {
		var serverError *server_errors.SError
		if errors.As(err, &serverError) {
			c.JSON(serverError.Unwrap())
			return
		}
		c.JSON(server_errors.InternalError.Unwrap())
		return
	}
	c.JSON(http.StatusOK, response)
}
