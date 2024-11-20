package utils

import (
	"time"

	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"
	"shirinec.com/config"
)


func HashPassword(password string) (string, error) {
    hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
    if err != nil {
        return "", err
    }
    return string(hashedPassword), nil
}

func GenerateAccessToken(id, email string) (string, error){
    expirationTime := time.Now().Add(config.AppConfig.AccessTokenDuration)
    return generateToken(id, email, expirationTime)

}

func GenerateRefreshToken(id, email string) (string, error){
    expirationTime := time.Now().Add(config.AppConfig.RefreshTokenDuration)
    return generateToken(id, email, expirationTime)
}

func generateToken(id, email string, exp time.Time) (string, error) {
    claims := jwt.MapClaims{
        "id": id,
        "email": email,
        "exp": exp.Unix(),
    }

    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    return token.SignedString([]byte(config.AppConfig.JWTSecret))
}
