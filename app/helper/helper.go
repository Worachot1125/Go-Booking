package helper

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
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

// คืน Cloudinary instance พร้อมใช้งาน
func NewCloudinary() (*cloudinary.Cloudinary, error) {
	cld, err := cloudinary.NewFromParams(
		os.Getenv("CLOUDINARY_CLOUD_NAME"),
		os.Getenv("CLOUDINARY_API_KEY"),
		os.Getenv("CLOUDINARY_API_SECRET"),
	)
	if err != nil {
		return nil, fmt.Errorf("cloudinary config error: %v", err)
	}
	return cld, nil
}

// อัปโหลดรูปภาพ พร้อมคืน URL
func UploadImageToCloudinary(file interface{}) (string, error) {
	cld, err := NewCloudinary()
	if err != nil {
		return "", err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	uploadResult, err := cld.Upload.Upload(ctx, file, uploader.UploadParams{
		Folder:   "user",
		PublicID: fmt.Sprintf("user_%d", time.Now().UnixNano()),
	})
	if err != nil {
		return "", fmt.Errorf("upload failed: %v", err)
	}

	return uploadResult.SecureURL, nil
}
