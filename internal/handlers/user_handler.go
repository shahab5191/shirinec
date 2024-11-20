package handler

import (
	"shirinec.com/internal/repositories"
)

type UserHandler struct{
    userRepo repositories.UserRepository
}

func NewUserHandler(userRepo repositories.UserRepository) *UserHandler {
    return &UserHandler{userRepo: userRepo}
}
