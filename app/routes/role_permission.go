package routes

import (
	"app/app/controller"
	"app/app/middleware"

	"github.com/gin-gonic/gin"
)

func Role_Permission(router *gin.RouterGroup) {
	// Get the *bun.DB instance from config
	ctl := controller.New() // Pass the *bun.DB to the controller
	role_permission := router.Group("")
	role_permission.Use(middleware.AuthMiddleware())
	{
		role_permission.POST("/create", ctl.Role_Permissionctl.Create)
		role_permission.PATCH("/update/:id", ctl.Role_Permissionctl.Update)
	}
}
