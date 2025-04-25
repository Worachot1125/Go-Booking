package routes

import (
	"app/app/controller"
	"app/app/middleware"

	"github.com/gin-gonic/gin"
)

func Permission(router *gin.RouterGroup) {
	// Get the *bun.DB instance from config
	ctl := controller.New() // Pass the *bun.DB to the controller
	permission := router.Group("")
	permission.Use(middleware.AuthMiddleware())
	{
		permission.POST("/create", ctl.PermissionCtl.Create)
		permission.PATCH("/update/:id", ctl.PermissionCtl.Update)
		permission.DELETE("/delete/:id", ctl.PermissionCtl.Delete)
	}
}
