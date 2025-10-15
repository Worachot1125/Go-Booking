package routes

import (
	"app/app/controller"
	"app/app/middleware"

	"github.com/gin-gonic/gin"
)

func Line(router *gin.RouterGroup) {
    ctl := controller.New().LineCtl
    // ขอ pairing code (ต้อง Auth)
    router.POST("/pairing-code", middleware.AuthMiddleware(), ctl.IssuePairingCode)
    // webhook (ไม่ต้อง Auth)
    router.POST("/webhook/line", ctl.Webhook)
    router.GET("/pairing-code/:id", ctl.GetPairingCodeByUserID)
}