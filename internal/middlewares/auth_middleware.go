package middlewares

import (
	"context"
	"database/sql"
	"errors"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"shirinec.com/internal/enums"
	server_errors "shirinec.com/internal/errors"
	"shirinec.com/internal/repositories"
	"shirinec.com/internal/utils"
)

type AuthMiddleWareFlags struct {
    ShouldBeActive bool
}

func AuthMiddleWare(flags AuthMiddleWareFlags, db *pgxpool.Pool) gin.HandlerFunc {
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

        if flags.ShouldBeActive {
            userRepo := repositories.NewUserRepository(db)
            uid, err := uuid.Parse(id)
            if err != nil {
                ctx.JSON(server_errors.InternalError.Unwrap())
                ctx.Abort()
                return
            }
            user, err := userRepo.GetByID(context.Background(), uid)
            if err != nil {
                if errors.Is(err, sql.ErrNoRows){
                    ctx.JSON(server_errors.UserNotFound.Unwrap())
                    ctx.Abort()
                    return
                }
                ctx.JSON(server_errors.InternalError.Unwrap())
                ctx.Abort()
                return
            }

            if user.Status != enums.StatusVerified{
                ctx.JSON(server_errors.AccountIsNotActive.Unwrap())
                ctx.Abort()
                return
            }
        }

        ctx.Set("user_id", id)

        ctx.Next()
    }
}
