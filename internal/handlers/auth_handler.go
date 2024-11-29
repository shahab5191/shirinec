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
	return &AuthHandler{authService: authService}
}

func (h *AuthHandler) SignUp(c *gin.Context) {
	var input dto.AuthSignupRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(server_errors.InvalidInput.Unwrap())
		return
	}

	response, err := h.authService.CreateUser(context.Background(), &input, c.ClientIP())
	if err != nil {
		if errors.Is(err, &server_errors.UserAlreadyExistsError) {
			c.JSON(server_errors.UserAlreadyExistsError.Code, server_errors.UserAlreadyExistsError.Message)
            return
		}
		log.Printf("Error Creating user: %+v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not create user!"})
		return
	}

	c.JSON(http.StatusCreated, *response)
}

func (h *AuthHandler) Login(c *gin.Context) {
	var credentials dto.AuthLoginRequest
	if err := c.ShouldBindJSON(&credentials); err != nil {
		log.Printf("ERROR: %s", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid credentials format!"})
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
	return
}

func (h *AuthHandler) RefreshToken(c *gin.Context) {
    var requestDTO dto.AuthRefreshTokenRequest
    if err := c.ShouldBindJSON(&requestDTO); err != nil {
        log.Printf("Error binding request to dto: %+v\n", err)
        c.JSON(server_errors.InvalidInput.Unwrap())
        return
    }

    response, err := h.authService.Refresh(requestDTO.RefreshToken)
    if err != nil {
        var serverError *server_errors.SError
        if errors.As(err, &serverError){
            c.JSON(serverError.Unwrap())
            return
        }
        c.JSON(server_errors.InternalError.Unwrap())
        return
    }
    c.JSON(http.StatusOK, response)
    return
}
