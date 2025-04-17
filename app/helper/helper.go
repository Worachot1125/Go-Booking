package helper

import (
	"encoding/json"

	"github.com/gin-gonic/gin"
)

type user struct {
	ID int64 `json:"id"`
}

func GetUserByToken(ctx *gin.Context) (int64, error) {
	claims, exist := ctx.Get("claims")
	if !exist {
		return 0, nil
	}
	var user user
	err := json.Unmarshal(claims.([]byte), &user)
	if err != nil {
		return 0, err
	}

	return user.ID, nil
}
