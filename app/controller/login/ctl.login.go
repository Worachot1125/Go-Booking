package login

import (
	"app/app/request"
	"app/app/response"
	"net/http"

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

	claims := jwt.MapClaims{
		"user_id": user.ID.String(),
		"email":   user.Email,
	}

	token, err := jwt.CreateToken(claims, viper.GetString("TOKEN_SECRET_USER"))
	if err != nil {
		response.InternalServerError(ctx, "Failed to generate token")
		return
	}

	// ✅ ส่ง token กลับ
	ctx.JSON(http.StatusOK, gin.H{
		"user":  user,
		"token": token,
	})
}
