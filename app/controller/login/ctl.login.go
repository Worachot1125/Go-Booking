package login

import (
	"app/app/request"
	"app/app/response"
	"net/http"
	"app/app/util/jwt"

	jwt5 "github.com/golang-jwt/jwt/v5"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

func (ctl *Controller) Login(ctx *gin.Context) {
	var req request.LoginRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.BadRequest(ctx, "Invalid input")
		return
	}

	user, err := ctl.Service.Login(ctx, req.Email, req.Password)
	if err != nil {
		response.Unauthorized(ctx, err.Error())
		return
	}

	claims := jwt5.MapClaims{
		"user_id": user.ID,
		"email":   user.Email,
	}

	token, err := jwt.CreateToken(claims, viper.GetString("TOKEN_SECRET_USER"))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"user":  user,
		"token": token,
	})
}
