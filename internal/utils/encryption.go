package utils

import (
	"errors"
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
    return generateToken(id, email, expirationTime, []byte(config.AppConfig.JWTSecret))

}

func GenerateRefreshToken(id, email string) (string, error){
    expirationTime := time.Now().Add(config.AppConfig.RefreshTokenDuration)
    return generateToken(id, email, expirationTime, []byte(config.AppConfig.JWTRefreshSecret))
}

func generateToken(id, email string, exp time.Time, secret []byte) (string, error) {
    claims := jwt.MapClaims{
        "id": id,
        "email": email,
        "exp": exp.Unix(),
    }

    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    signedToken, err := token.SignedString(secret)
    if err != nil {
        return "", &server_errors.InternalError
    }
    return signedToken, nil
}

func ParseRefreshToken(refreshToken string) (jwt.MapClaims, error) {
    return parseToken(refreshToken, []byte(config.AppConfig.JWTRefreshSecret))
    
}

func ParseAccessToken(accessToken string) (jwt.MapClaims, error) {
    return parseToken(accessToken, []byte(config.AppConfig.JWTSecret))
}

func parseToken(token string, secret []byte) (jwt.MapClaims, error) {
    parsedToken, err := jwt.ParseWithClaims(token, jwt.MapClaims{}, func(token *jwt.Token) (interface{}, error) {
        return secret, nil
    })
    if err != nil {
        log.Printf("Error parsing refresh token: %+v\n", err)
        var validationErr *jwt.ValidationError
        if errors.As(err, &validationErr){
            if validationErr.Errors&jwt.ValidationErrorMalformed != 0 {
                return nil, &server_errors.TokenMalformed
            }else if validationErr.Errors&jwt.ValidationErrorExpired != 0 {
                return nil, &server_errors.TokenExpired
            }else if validationErr.Errors&jwt.ValidationErrorSignatureInvalid != 0{
                return nil, &server_errors.TokenSignatureInvalid
            }else{
                return nil, &server_errors.InternalError
            }
        }
        return nil, &server_errors.InternalError
    }

    if claims, ok := parsedToken.Claims.(jwt.MapClaims); ok && parsedToken.Valid {
        return claims, nil
    }else{
        log.Printf("parsedToken: %+v\n", claims)
    }

    return nil, &server_errors.InvalidToken
}
