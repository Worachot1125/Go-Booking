package login

import (
	"app/app/request"
	"app/app/response"
	"app/app/util/jwt"
	"net/http"

	"github.com/gin-gonic/gin"
	jwt5 "github.com/golang-jwt/jwt/v5"
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

	// สร้าง token
	claims := jwt5.MapClaims{
		"user_id": user.ID,
		"email":   user.Email,
	}

	token, err := jwt.CreateToken(claims, viper.GetString("TOKEN_SECRET_USER"))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	// ตั้งค่า cookie สำหรับ token
	ctx.SetCookie(
		"token", // ชื่อของ cookie
		token,   // ค่า (value) ของ cookie
		3600,    // อายุของ cookie (ในที่นี้คือ 1 ชั่วโมง)
		"/",     // path ที่ cookie จะใช้ได้
		"",      // domain (ถ้าต้องการสามารถกำหนดได้)
		false,   // secure (ตั้งเป็น true ถ้าใช้ https)
		true,    // httpOnly (ตั้งเป็น true เพื่อป้องกันการเข้าถึง cookie จาก JavaScript)
	)

	// ส่งกลับ user และ token ใน response (หากจำเป็นต้องส่ง)
	ctx.JSON(http.StatusOK, gin.H{
		"user":  user,
		"token": token,
	})
}
