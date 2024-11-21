package utils

import (
	"fmt"
	"log"
	"time"

	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"
	"shirinec.com/config"
	server_errors "shirinec.com/internal/errors"
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
    signedToken, err := token.SignedString([]byte(config.AppConfig.JWTSecret))
    if err != nil {
        return "", &server_errors.InternalError
    }
    return signedToken, nil
}

func ParseToken(refreshToken string) (jwt.MapClaims, error) {
    parsedToken, err := jwt.ParseWithClaims(refreshToken, jwt.MapClaims{}, func(token *jwt.Token) (interface{}, error) {
        return []byte(config.AppConfig.JWTSecret), nil
    })
    if err != nil {
        return nil, err
    }

    if claims, ok := parsedToken.Claims.(jwt.MapClaims); ok && parsedToken.Valid {
        return claims, nil
    }else{
        log.Printf("parsedToken: %+v\n", claims)
    }

    return nil, fmt.Errorf("Invalid token")
}
