package services

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"golang.org/x/crypto/bcrypt"
	"shirinec.com/src/internal/db"
	"shirinec.com/src/internal/dto"
	"shirinec.com/src/internal/errors"
	"shirinec.com/src/internal/repositories"
	"shirinec.com/src/internal/utils"
)

type UserService interface {
	NewPassword(ctx context.Context, input dto.UserUpdatePasswordRequest, userID uuid.UUID) error
	NewEmail(ctx context.Context, input dto.UserUpdateEmailRequest, userID uuid.UUID) (int, error)
	NewEmailVerification(ctx context.Context, verificationCode int, userID uuid.UUID) error
	SignupVerification(ctx context.Context, verificationCode int, userID uuid.UUID) error
}

type userService struct {
	userRepo repositories.UserRepository
}

func NewUserService(userRepo repositories.UserRepository) UserService {
	return &userService{userRepo: userRepo}
}

func (s *userService) NewPassword(ctx context.Context, input dto.UserUpdatePasswordRequest, userID uuid.UUID) error {
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return &server_errors.UserNotFound
		}
		utils.Logger.Errorf("userService.NewPassword - getting user by id:%s", err.Error())
		return &server_errors.InternalError
	}

	if err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.CurrentPassword)); err != nil {
		utils.Logger.Errorf("userService.NewPassword - Comparing password with saved hash:%s", err.Error())
		return &server_errors.CredentialError
	}

	hashedNewPassword, err := utils.HashPassword(input.NewPassword)
	if err != nil {
		utils.Logger.Errorf("userService.NewPassword - hashing newPassword: %s", err.Error())
		return &server_errors.InternalError
	}

	err = s.userRepo.UpdatePassword(context.Background(), hashedNewPassword, userID)
	if err != nil {
		utils.Logger.Errorf("userService.NewPassword - Updating password:%s", err.Error())
		return &server_errors.InternalError
	}

	return nil
}

func (s *userService) NewEmail(ctx context.Context, input dto.UserUpdateEmailRequest, userID uuid.UUID) (int, error) {
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, &server_errors.UserNotFound
		}
		utils.Logger.Errorf("userService.NewEmail - Getting user from repo%s", err.Error())
		return 0, &server_errors.InternalError
	}

	if err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.CurrentPassword)); err != nil {
		utils.Logger.Errorf("userService.NewEmail - Comparing password with saved hash:%s", err.Error())
		return 0, &server_errors.CredentialError
	}

	verificationCode := utils.GenerateVerificationCode()

	fields := map[string]interface{}{
		"userID":   userID.String(),
		"newEmail": input.NewEmail,
	}

	rKey := fmt.Sprintf("new_email:%d", verificationCode)
	_, err = db.Redis.HSet(ctx, rKey, fields).Result()
	if err != nil {
		utils.Logger.Errorf("userService.NewEmail - setting verification code to redis:%s", err.Error())
		return 0, &server_errors.InternalError
	}
	_, err = db.Redis.Expire(ctx, rKey, 5*time.Minute).Result()
	if err != nil {
		utils.Logger.Errorf("userService.NewEmail - setting Expire code to redis:%s", err.Error())
		return 0, &server_errors.InternalError
	}
	return verificationCode, nil
}

func (s *userService) NewEmailVerification(ctx context.Context, verificationCode int, userID uuid.UUID) error {
	rKey := fmt.Sprintf("new_email:%d", verificationCode)

	res, err := db.Redis.HGetAll(ctx, rKey).Result()
	if err != nil {
		if err == redis.Nil {
			return &server_errors.InvalidVerificationCode
		}
		utils.Logger.Errorf("userService.NewEMailVerification - getting email change inputs from redis:%s", err.Error())
		return &server_errors.InternalError
	}

	if userID.String() != res["userID"] {
		return &server_errors.InvalidVerificationCode
	}

	newEmail := res["newEmail"]
	if newEmail == "" {
		return &server_errors.InvalidVerificationCode
	}

	err = s.userRepo.UpdateEmail(ctx, newEmail, userID)
	if err != nil {
		utils.Logger.Errorf("userService.NewEmail - Updating email:%s", err.Error())
		return &server_errors.InternalError
	}
	return nil
}

func (s *userService) SignupVerification(ctx context.Context, verificationCode int, userID uuid.UUID) error {
	rKey := fmt.Sprintf("signup:%d", verificationCode)

	res, err := db.Redis.Get(ctx, rKey).Result()
	if err != nil {
		if err == redis.Nil {
			return &server_errors.InvalidVerificationCode
		}
		utils.Logger.Errorf("userService.SignupVerification - Getting email change inputs from redis:%s", err.Error())
		return &server_errors.InternalError
	}

	if userID.String() != res {
		return &server_errors.InvalidVerificationCode
	}

	db.Redis.Del(ctx, rKey)

	err = s.userRepo.VerifyUser(ctx, userID)
	return err
}
