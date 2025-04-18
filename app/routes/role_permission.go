package routes

import (
	"app/app/controller"

	"github.com/gin-gonic/gin"
)

func Role_Permission(router *gin.RouterGroup) {
	// Get the *bun.DB instance from config
	ctl := controller.New() // Pass the *bun.DB to the controller
	role_permission := router.Group("")
	{
		role_permission.POST("/create",ctl.Role_Permissionctl.Create)
	}
}