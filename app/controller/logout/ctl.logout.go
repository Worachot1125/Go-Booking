package logout

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (ctl *Controller) Logout(ctx *gin.Context) {
	// เรียกใช้งาน service logout
	err := ctl.Service.Logout(ctx)
	if err != nil {
		// ถ้ามีข้อผิดพลาด ให้ตอบกลับด้วย error
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to log out"})
		return
	}

	// ส่งผลลัพธ์การ logout
	ctx.JSON(http.StatusOK, gin.H{"message": "Logged out successfully"})
}
