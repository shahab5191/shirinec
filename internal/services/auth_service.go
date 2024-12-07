package services

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"golang.org/x/crypto/bcrypt"
	"shirinec.com/internal/db"
	"shirinec.com/internal/dto"
	"shirinec.com/internal/errors"
	"shirinec.com/internal/models"
	"shirinec.com/internal/repositories"
	"shirinec.com/internal/utils"
)

type AuthService interface {
	CreateUser(ctx context.Context, input *dto.AuthSignupRequest, ip string) (*dto.AuthLoginResponse, error)
	Login(email, password, ip string) (*dto.AuthLoginResponse, error)
	Refresh(token string) (*dto.AuthLoginResponse, error)
}

type authService struct {
	userRepo  repositories.UserRepository
	jwtSecret string
}

func NewAuthService(userRepo repositories.UserRepository, jwtSecret string) AuthService {
	return &authService{jwtSecret: jwtSecret, userRepo: userRepo}
}

func (s *authService) CreateUser(ctx context.Context, input *dto.AuthSignupRequest, ip string) (*dto.AuthLoginResponse, error) {
	password, err := utils.HashPassword(input.Password)
	if err != nil {
		utils.Logger.Errorf("authService.Create - Calling utils.HashingPassword: %s", err.Error())
		return nil, &server_errors.InternalError
	}

	_, err = s.userRepo.GetByEmail(ctx, input.Email)
	if err != nil {
		if !errors.Is(err, pgx.ErrNoRows) {
			utils.Logger.Errorf("ErrorauthService.Create - Calling userRepo.GetByEmail: %s", err.Error())
			return nil, &server_errors.InternalError
		}
	} else {
		return nil, &server_errors.UserAlreadyExistsError
	}

	user := models.User{
		ID:       uuid.New(),
		Email:    input.Email,
		IP:       ip,
		Password: password,
	}

	err = s.userRepo.Create(ctx, &user)
	if err != nil {
		utils.Logger.Errorf("authService.CreateUser - Calling userRepo.Create: %+v", err)
		return nil, &server_errors.InternalError
	}

	accessToken, err := utils.GenerateAccessToken(user.ID.String(), user.Email, user.LastPasswordChange)
	if err != nil {
		utils.Logger.Errorf("authService.CreateUser - Calling utils.GenerateAccessToken: %+v", err)
		return nil, &server_errors.InternalError
	}
	refreshToken, err := utils.GenerateRefreshToken(user.ID.String(), user.Email, user.LastPasswordChange)
	if err != nil {
		utils.Logger.Errorf("authService.CreateUser - Calling utils.GenerateRefreshToken: %+v", err)
		return nil, &server_errors.InternalError
	}

	verificationCode := utils.GenerateVerificationCode()

	rKey := fmt.Sprintf("signup:%d", verificationCode)
	_, err = db.Redis.SetEx(ctx, rKey, user.ID.String(), 5*time.Minute).Result()
	if err != nil {
		utils.Logger.Errorf("authService.CreateUser - Setting verification code to redis: %s", err.Error())
		return nil, &server_errors.InternalError
	}

	response := dto.AuthLoginResponse{
		ID:               user.ID,
		AccessToken:      accessToken,
		RefreshToken:     refreshToken,
		VerificationCode: verificationCode,
	}

	return &response, nil
}

func (s *authService) Login(email, password, ip string) (*dto.AuthLoginResponse, error) {
	var res dto.AuthLoginResponse
	user, err := s.userRepo.GetByEmail(context.Background(), email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			utils.Logger.Infof("authService.Login - Calling userRepo.GetByEmail - No rows error for email: %s", email)
			return &res, &server_errors.CredentialError
		}
		utils.Logger.Errorf("authService.Login - Calling userRepo.GetByEmail: %s", err.Error())
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		utils.Logger.Infof("authService.Login - Calling bcrypt.CompareHashAndPassword: %s", err.Error())
		return nil, &server_errors.CredentialError
	}

	accessToken, err := utils.GenerateAccessToken(user.ID.String(), user.Email, user.LastPasswordChange)
	if err != nil {
		utils.Logger.Errorf("authService.Login - Calling utils.GenerateAccessToken: %+v", err)
		return nil, &server_errors.InternalError
	}
	refreshToken, err := utils.GenerateRefreshToken(user.ID.String(), user.Email, user.LastPasswordChange)
	if err != nil {
		utils.Logger.Errorf("authService.Login - Calling utils.GenerateRefreshToken: %+v", err)
		return nil, &server_errors.InternalError
	}

	err = s.userRepo.Login(context.Background(), ip)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, &server_errors.UserNotFound
		}
		utils.Logger.Errorf("authService.Login - Calling userRepo.Login: %s", err.Error())
		return nil, &server_errors.InternalError
	}

	response := dto.AuthLoginResponse{
		ID:           user.ID,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}

	return &response, nil
}

func (s *authService) Refresh(token string) (*dto.AuthLoginResponse, error) {
	claims, err := utils.ParseRefreshToken(token)
	if err != nil {
		return nil, err
	}

	id, ok := claims["id"].(string)
	if !ok {
		utils.Logger.Error("authService.Refresh - Getting id from claims")
		return nil, &server_errors.InternalError
	}

	email, ok := claims["email"].(string)
	if !ok {
		utils.Logger.Error("authService.Refresh - Getting email from claims")
		return nil, &server_errors.InternalError
	}

	uuID, err := uuid.Parse(id)
	if err != nil {
		utils.Logger.Errorf("authService.Refresh - Parsing uuid from id string: %s", err.Error())
		return nil, &server_errors.InternalError
	}

	lastPasswordChangeUnixFloat, ok := claims["lastPasswordChange"].(float64)
	if !ok {
		utils.Logger.Error("authService.Refresh - Getting lastPasswordChange from calims")
		return nil, &server_errors.InternalError
	}
	lastPasswordChangeUnixInt := int64(lastPasswordChangeUnixFloat)
	lastPasswordChange := time.Unix(int64(lastPasswordChangeUnixInt), 0).UTC()

	user, err := s.userRepo.GetByID(context.Background(), uuID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, &server_errors.UserNotFound
		}
		utils.Logger.Errorf("authService.Refersh - Calling userRepo.GetById: %s", err.Error())
		return nil, &server_errors.InternalError
	}

	if user.ID != uuID || user.Email != email || user.LastPasswordChange != lastPasswordChange {
		return nil, &server_errors.CredentialError
	}

	accessToken, err := utils.GenerateAccessToken(id, email, lastPasswordChange)
	if err != nil {
		utils.Logger.Errorf("Refresh - Calling utils.GenerateAccessToken: %s", err.Error())
		return nil, err
	}

	uid, err := uuid.Parse(id)
	if err != nil {
		utils.Logger.Errorf("Refresh - Parsing uuid from id string: %s", err.Error())
		return nil, &server_errors.InternalError
	}

	loginResponse := dto.AuthLoginResponse{ID: uid, AccessToken: accessToken, RefreshToken: token}

	return &loginResponse, nil
}
