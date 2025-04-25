package routes

import (
	"app/app/controller"
	"github.com/gin-gonic/gin"
)

func Login(router *gin.RouterGroup) {
	ctl := controller.New()
	router.POST("/", ctl.LoginCtl.Login)
}