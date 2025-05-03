package middleware

import (
	"app/app/util/jwt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var token string

		// ดึงจาก Authorization header
		authHeader := ctx.GetHeader("Authorization")
		if authHeader != "" && strings.HasPrefix(authHeader, "Bearer ") {
			token = strings.TrimPrefix(authHeader, "Bearer ")
		} else {
			// ตรวจ token จาก cookie
			cookieToken, err := ctx.Cookie("token")
			if err != nil || cookieToken == "" {
				ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Authorization token is required"})
				return
			}
			token = cookieToken
		}

		claims, err := jwt.VerifyToken(token)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			return
		}

		ctx.Set("claims", claims)
		ctx.Next()
	}
}

