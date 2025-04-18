package routes

import (
	"app/app/controller"

	"github.com/gin-gonic/gin"
)

func Permission(router *gin.RouterGroup) {
	// Get the *bun.DB instance from config
	ctl := controller.New() // Pass the *bun.DB to the controller
	permission := router.Group("")
	{
		permission.POST("/create",ctl.PermissionCtl.Create)
		permission.POST("/update/:id", ctl.PermissionCtl.Update)
		permission.DELETE("/delete/:id", ctl.PermissionCtl.Delete)
	}
}