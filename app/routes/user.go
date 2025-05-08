package routes

import (
	"app/app/controller"
	"app/app/middleware"

	"github.com/gin-gonic/gin"
)

func User(router *gin.RouterGroup) {
	ctl := controller.New()
	md := middleware.AuthMiddleware()
	user := router.Group("")
	{
		user.POST("/register", ctl.UserCtl.Create)
		user.PATCH("/:id", md, ctl.UserCtl.Update)
		user.GET("/list", md, ctl.UserCtl.List)
		user.GET("/:id", md, ctl.UserCtl.Get)
		user.DELETE("/:id", md, ctl.UserCtl.Delete)
	}
}
