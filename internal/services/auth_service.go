package services

import (
	"context"
	"log"
	"time"

	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"
	"shirinec.com/internal/dto"
	"shirinec.com/internal/errors"
	"shirinec.com/internal/repositories"
	"shirinec.com/internal/utils"
)

type AuthService struct{
    userRepo     repositories.UserRepository
    jwtSecret   string
}

func NewAuthService(userRepo repositories.UserRepository, jwtSecret string) *AuthService{
    return &AuthService{jwtSecret: jwtSecret, userRepo: userRepo}
}


func (s *AuthService) Login(email, password string) (dto.LoginResponse, error) {
    var res dto.LoginResponse
    user, err := s.userRepo.GetByEmail(context.Background(), email)
    if err != nil {
        log.Printf("Error getting user by email from db: %s", err)
        return res, &server_errors.CredentialError
    }

    hashedPassword, err := utils.HashPassword(password)
    if err != nil {
        log.Printf("ERROR: %s", err)
        return res, &server_errors.InternalError
    }

    if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
        log.Printf("Password is not correct")
        log.Printf("Provided pass: %s\nDatabase pass: %s\n", hashedPassword, user.Password)
        return res, &server_errors.CredentialError
    }

    token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
        "id": user.ID,
        "email": user.Email,
        "exp": time.Now().Add(time.Hour * 24).Unix(),
    })

    tokenString, err := token.SignedString([]byte(s.jwtSecret))
    if err != nil {
        log.Printf("ERROR: %s", err)
        return res, &server_errors.InternalError
    }

    return tokenString, nil
}
