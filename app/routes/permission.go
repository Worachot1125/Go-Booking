package routes

import (
	"app/app/controller"
	"app/app/middleware"

	"github.com/gin-gonic/gin"
)

func Permission(router *gin.RouterGroup) {
	// Get the *bun.DB instance from config
	ctl := controller.New() // Pass the *bun.DB to the controller
	md := middleware.AuthMiddleware()
	permission := router.Group("")
	{
		permission.POST("/create", md, ctl.PermissionCtl.Create)
		permission.PATCH("/:id", md, ctl.PermissionCtl.Update)
		permission.DELETE("/:id", md, ctl.PermissionCtl.Delete)
	}
}
