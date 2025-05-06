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
			// ดึงจาก cookie
			cookieToken, err := ctx.Cookie("token")
			if err != nil || cookieToken == "" {
				ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Authorization token is required"})
				return
			}
			token = cookieToken
		}

		// ตรวจสอบ token
		claims, err := jwt.VerifyToken(token)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			return
		}

		mapClaims := claims 

		ctx.Set("claims", mapClaims)
		ctx.Next()
	}
}
