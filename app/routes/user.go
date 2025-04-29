package routes

import (
	"app/app/controller"
	"app/app/middleware"

	"github.com/gin-gonic/gin"
)

func User(router *gin.RouterGroup) {
	// Get the *bun.DB instance from config
	ctl := controller.New() // Pass the *bun.DB to the controller
	user := router.Group("")
	{
		user.POST("/register", ctl.UserCtl.Create)
	}
		user.Use(middleware.AuthMiddleware())
	{
		user.PATCH("/:id", ctl.UserCtl.Update)
		user.GET("/list", ctl.UserCtl.List)
		user.GET("/:id", ctl.UserCtl.Get)
		user.DELETE("/:id", ctl.UserCtl.Delete)

	}
}
