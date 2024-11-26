package services

import (
	"context"
	"database/sql"
	"errors"
	"log"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"shirinec.com/internal/dto"
	"shirinec.com/internal/errors"
	"shirinec.com/internal/repositories"
	"shirinec.com/internal/utils"
)

type UserService interface {
	NewPassword(ctx context.Context, input dto.UpdatePasswordRequest, userID uuid.UUID) error
	NewEmail(ctx context.Context, input dto.UpdateEmailRequest, userID uuid.UUID) error
}

type userService struct {
	userRepo repositories.UserRepository
}

func NewUserService(userRepo repositories.UserRepository) UserService {
	return &userService{userRepo: userRepo}
}

func (s *userService) NewPassword(ctx context.Context, input dto.UpdatePasswordRequest, userID uuid.UUID) error {
	user, err := s.userRepo.GetByID(context.Background(), userID)
	if err != nil {
		log.Printf("[Error] - userService.NewPassword - getting user by id: %+v\n", err)
		if errors.Is(err, sql.ErrNoRows) {
			return &server_errors.UserNotFound
		}
		return &server_errors.InternalError
	}

	if err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.CurrentPassword)); err != nil {
		log.Printf("[Error] - userService.NewPassword - Comparing password with saved hash: %+v\n", err)
		return &server_errors.CredentialError
	}

	hashedNewPassword, err := utils.HashPassword(input.NewPassword)
	if err != nil {
		log.Printf("[Error] - userService.NewPassword - hashing newPassword:  %+v\n", err)
		return &server_errors.InternalError
	}

	err = s.userRepo.UpdatePassword(context.Background(), hashedNewPassword, userID)
	if err != nil {
        log.Printf("[Error] - userService.NewPassword - Updating password: %+v\n", err)
		return &server_errors.InternalError
	}

	return nil
}

func (s *userService) NewEmail(ctx context.Context, input dto.UpdateEmailRequest, userID uuid.UUID) error {
	return nil
}
