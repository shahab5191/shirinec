package middlewares

import (
	"errors"
	"strings"

	"github.com/gin-gonic/gin"
	server_errors "shirinec.com/internal/errors"
	"shirinec.com/internal/utils"
)

func AuthMiddleWare() gin.HandlerFunc {
    return func(ctx *gin.Context) {
        authHeader := ctx.GetHeader("Authorization")
        if authHeader == "" {
            ctx.JSON(server_errors.Unauthorized.Unwrap())
            ctx.Abort()
            return
        }

        tokenParts := strings.Split(authHeader, " ")
        if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
            ctx.JSON(server_errors.InvalidAuthorizationHeader .Unwrap())
            ctx.Abort()
            return
        }

        claims, err := utils.ParseAccessToken(tokenParts[1])
        if err != nil {
            var serverError *server_errors.SError
            if errors.As(err, &serverError) {
                ctx.JSON(serverError.Unwrap())
            }else{
                ctx.JSON(server_errors.InternalError.Unwrap())
            }
            ctx.Abort()
            return
        }

        id, ok := claims["id"].(string)
        if !ok {
            ctx.JSON(server_errors.InternalError.Unwrap())
            ctx.Abort()
            return
        }
        ctx.Set("user_id", id)

        ctx.Next()
    }
}
