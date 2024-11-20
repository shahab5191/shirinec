package handler

import (
	"context"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"shirinec.com/internal/dto"
	"shirinec.com/internal/models"
	"shirinec.com/internal/repositories"
	"shirinec.com/internal/utils"
)

type UserHandler struct{
    userRepo repositories.UserRepository
}

func NewUserHandler(userRepo repositories.UserRepository) *UserHandler {
    return &UserHandler{userRepo: userRepo}
}

func (h *UserHandler) CreateUser(c *gin.Context) {
    var input dto.CreateUserDTO
    if err := c.ShouldBindJSON(&input); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Input!"})
        return
    }

    password, err := utils.HashPassword(input.Password)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not create user!"})
    }

    user := models.User {
        ID: uuid.New(),
        Email: input.Email,
        IP: c.ClientIP(),
        Password: password,
    }

    err = h.userRepo.Create(context.Background(), &user)
    if err != nil {
        log.Printf("ERROR: %s", err)
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not create user!"})
        return
    }

    c.JSON(http.StatusCreated, gin.H{"id": user.ID})
}
