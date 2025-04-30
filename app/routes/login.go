package routes

import (
	"app/app/controller"
	"net/http"

	"github.com/gin-gonic/gin"
)

func Login(router *gin.RouterGroup) {
	ctl := controller.New()

	// รองรับ preflight OPTIONS
	router.Handle("OPTIONS", "/", func(c *gin.Context) {
		c.Status(http.StatusNoContent)
	})
	router.Handle("OPTIONS", "", func(c *gin.Context) {
		c.Status(http.StatusNoContent)
	})
	router.POST("/", ctl.LoginCtl.Login)
	router.POST("", ctl.LoginCtl.Login)
}
