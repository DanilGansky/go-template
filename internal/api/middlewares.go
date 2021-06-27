package api

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/littlefut/go-template/internal/auth"
	"github.com/littlefut/go-template/internal/hash"
	"github.com/littlefut/go-template/internal/user"
)

const TokenPrefix = "Bearer "

func AuthorizationMiddleware(userSvc user.Service, tokenSvc hash.TokenService) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		token := ctx.GetHeader("Authorization")
		if token == "" {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, auth.ErrTokenDoesNotValid)
			return
		}
		if strings.Contains(token, TokenPrefix) {
			token = token[len(TokenPrefix):]
		}

		id, isValid := tokenSvc.Validate(token)
		if !isValid {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, auth.ErrTokenDoesNotValid)
			return
		}

		_, err := userSvc.FindByID(ctx, id)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, user.ErrNotFound)
			return
		}

		ctx.Set("user_id", id)
		ctx.Next()
	}
}
