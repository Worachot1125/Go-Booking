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

		authHeader := ctx.GetHeader("Authorization")
		if authHeader != "" && strings.HasPrefix(authHeader, "Bearer ") {
			token = strings.TrimPrefix(authHeader, "Bearer ")
		} else {
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

		// ดึง user_id จาก claims และเซ็ตไว้ใน context
		mapClaims := claims
		userID, ok := mapClaims["user_id"].(string)
		if !ok || userID == "" {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid token: user_id not found"})
			return
		}

		ctx.Set("user_id", userID)
		ctx.Set("claims", mapClaims)
		ctx.Next()
	}
}
