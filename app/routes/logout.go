package routes

import (
	"app/app/controller"

	"github.com/gin-gonic/gin"
)

func Logout(router *gin.RouterGroup) {
	ctl := controller.New()

	router.POST("/", ctl.LogoutCtl.Logout)
	router.POST("", ctl.LogoutCtl.Logout)
}
