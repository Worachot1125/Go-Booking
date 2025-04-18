package routes

import (
	"app/app/controller"

	"github.com/gin-gonic/gin"
)

func Role(router *gin.RouterGroup) {
	ctl := controller.New() // Pass the *bun.DB to the controller
	role := router.Group("")
	{
		role.POST("/create", ctl.RoleCtl.Create)
		role.DELETE("/delete/:id", ctl.RoleCtl.Delete)
	}
}