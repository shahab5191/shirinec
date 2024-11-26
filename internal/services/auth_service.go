package services

import (
	"context"
	"database/sql"
	"errors"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"golang.org/x/crypto/bcrypt"
	"shirinec.com/internal/dto"
	"shirinec.com/internal/errors"
	"shirinec.com/internal/models"
	"shirinec.com/internal/repositories"
	"shirinec.com/internal/utils"
)

type AuthService interface {
	CreateUser(ctx context.Context, input *dto.CreateUserRequest, ip string) (dto.LoginResponse, error)
	Login(email, password, ip string) (dto.LoginResponse, error)
	Refresh(token string) (*dto.LoginResponse, error)
}

type authService struct {
	userRepo  repositories.UserRepository
	jwtSecret string
}

func NewAuthService(userRepo repositories.UserRepository, jwtSecret string) AuthService {
	return &authService{jwtSecret: jwtSecret, userRepo: userRepo}
}

func (s *authService) CreateUser(ctx context.Context, input *dto.CreateUserRequest, ip string) (dto.LoginResponse, error) {
	var response dto.LoginResponse
	password, err := utils.HashPassword(input.Password)
	if err != nil {
		log.Printf("Error hashing password: %+v\n", err)
		return response, &server_errors.InternalError
	}

	existingUser, err := s.userRepo.GetByEmail(ctx, input.Email)
	if err != nil {
		if !errors.Is(err, pgx.ErrNoRows) {
			log.Printf("Error getting user by email! %s\n", err)
			return response, &server_errors.InternalError
		}
	} else {
		return response, &server_errors.UserAlreadyExistsError
	}

	log.Printf("%+v\n", existingUser)

	user := models.User{
		ID:       uuid.New(),
		Email:    input.Email,
		IP:       ip,
		Password: password,
	}

	err = s.userRepo.Create(ctx, &user)
	if err != nil {
		log.Printf("Error creating user in datbase: %s", err)
		return response, &server_errors.InternalError
	}

	log.Printf("v+%\n", user)

	accessToken, err := utils.GenerateAccessToken(user.ID.String(), user.Email, user.LastPasswordChange)
	if err != nil {
		log.Printf("Error generating access token: %s", err)
		return response, &server_errors.InternalError
	}
	refreshToken, err := utils.GenerateRefreshToken(user.ID.String(), user.Email, user.LastPasswordChange)
	if err != nil {
		log.Printf("Error generating refresh token: %s", err)
		return response, &server_errors.InternalError
	}

	response = dto.LoginResponse{
		ID:           user.ID,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}

	return response, nil
}

func (s *authService) Login(email, password, ip string) (dto.LoginResponse, error) {
	var res dto.LoginResponse
	user, err := s.userRepo.GetByEmail(context.Background(), email)
	if err != nil {
		log.Printf("Error getting user by email from db: %s", err)
		return res, &server_errors.CredentialError
	}

	hashedPassword, err := utils.HashPassword(password)
	if err != nil {
		log.Printf("Error hashing password: %s", err)
		return res, &server_errors.InternalError
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		log.Printf("Password is not correct")
		log.Printf("Provided pass: %s\nDatabase pass: %s\n", hashedPassword, user.Password)
		return res, &server_errors.CredentialError
	}

	accessToken, err := utils.GenerateAccessToken(user.ID.String(), user.Email, user.LastPasswordChange)
	if err != nil {
		log.Printf("Error generating access token: %s", err)
		return res, &server_errors.InternalError
	}
	refreshToken, err := utils.GenerateRefreshToken(user.ID.String(), user.Email, user.LastPasswordChange)
	if err != nil {
		log.Printf("Error generating refresh token: %s", err)
		return res, &server_errors.InternalError
	}

    err = s.userRepo.Login(context.Background(), ip)

	response := dto.LoginResponse{
		ID:           user.ID,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}

	return response, nil
}

func (s *authService) Refresh(token string) (*dto.LoginResponse, error) {
	claims, err := utils.ParseRefreshToken(token)
	if err != nil {
		return nil, err
	}

	id, ok := claims["id"].(string)
	if !ok {
		log.Println("Error fetching id from token claims")
		return nil, &server_errors.InternalError
	}

	email, ok := claims["email"].(string)
	if !ok {
		log.Println("Error fetcing email from token claims")
		return nil, &server_errors.InternalError
	}

	uuID, err := uuid.Parse(id)
	if err != nil {
		log.Printf("[Error] - authService.Refresh - casting id to uuid %+v\n", uuID)
		return nil, &server_errors.InternalError
	}

	lastPasswordChangeUnixFloat, ok := claims["lastPasswordChange"].(float64)
	if !ok {
		log.Println("[Error] - authService.Refresh - could not get lastPasswordChange float64 from token claims")
		return nil, &server_errors.InternalError
	}
	lastPasswordChangeUnixInt := int64(lastPasswordChangeUnixFloat)
	lastPasswordChange := time.Unix(int64(lastPasswordChangeUnixInt), 0).UTC()
	if err != nil {
		log.Println("Error parsing lastPasswordChange to time format")
		return nil, &server_errors.InternalError
	}

	user, err := s.userRepo.GetByID(context.Background(), uuID)
	if err != nil {
		log.Printf("[Error] - authService.Refersh - getting user by id")
		if errors.Is(err, sql.ErrNoRows) {
			return nil, &server_errors.UserNotFound
		}
		return nil, &server_errors.InternalError
	}

	if user.ID != uuID || user.Email != email || user.LastPasswordChange != lastPasswordChange {
		log.Printf("user.LastpasswordChange: %v\nlastPasswordChange: %v\n", user.LastPasswordChange, lastPasswordChange)
		log.Printf("id: %b\nemail: %b\nlastPasswordChange: %b", user.ID != uuID, user.Email != email, user.LastPasswordChange != lastPasswordChange)
		return nil, &server_errors.CredentialError
	}

	accessToken, err := utils.GenerateAccessToken(id, email, lastPasswordChange)
	if err != nil {
		log.Printf("[Error] - Refresh - generating access token: %+v\n", err)
		return nil, err
	}

	uid, err := uuid.Parse(id)
	if err != nil {
		log.Printf("[Error] - Refresh - Error parsing id from string to uuid: %+v\n", err)
	}

	loginResponse := dto.LoginResponse{ID: uid, AccessToken: accessToken, RefreshToken: token}

	return &loginResponse, nil
}
