package handler

import (
	"errors"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"shirinec.com/internal/dto"
	server_errors "shirinec.com/internal/errors"
	"shirinec.com/internal/services"
)

type AuthHandler struct{
    authService services.AuthService
}

func NewAuthHandler(authService services.AuthService) *AuthHandler{
    return &AuthHandler{authService: authService}
}

func (h *AuthHandler) Login(c *gin.Context) {
    var credentials dto.LoginRequest
    if err := c.ShouldBindJSON(&credentials); err != nil {
        log.Printf("ERROR: %s", err)
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid credentials format!"})
        return
    }

    token, err := h.authService.Login(credentials.Email, credentials.Password)
    var se *server_errors.SError
    if err != nil {
        if errors.As(err, &se){
            c.JSON(se.Code, gin.H{"error": se.Message})
            return
        }
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error!"})
        return
    }

    c.JSON(http.StatusOK, gin.H{"result": token})
    return
}
