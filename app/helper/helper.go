package helper

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

// type user struct {
// 	ID int64 `json:"id"`
// }

func GetUserByToken(ctx *gin.Context) (int64, error) {
	// ตรวจสอบ claims ที่เก็บใน context
	claims, exist := ctx.Get("claims")
	if !exist {
		return 0, fmt.Errorf("no claims in context")
	}

	// แปลง claims เป็น map[string]interface{}
	mapClaims, ok := claims.(map[string]interface{})
	if !ok {
		return 0, fmt.Errorf("invalid token claims format")
	}

	// ดึง user_id จาก mapClaims และแปลงเป็น float64
	userIDFloat, ok := mapClaims["user_id"].(float64)
	if !ok {
		return 0, fmt.Errorf("user_id not found or invalid in token")
	}

	// คืนค่า user_id เป็น int64
	return int64(userIDFloat), nil
}
