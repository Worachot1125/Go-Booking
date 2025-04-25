package routes

import (
	"app/app/controller"
	"app/app/middleware"

	"github.com/gin-gonic/gin"
)

func Role(router *gin.RouterGroup) {
	ctl := controller.New() // Pass the *bun.DB to the controller
	role := router.Group("")
	role.Use(middleware.AuthMiddleware())
	{
		role.POST("/create", ctl.RoleCtl.Create)
		role.DELETE("/delete/:id", ctl.RoleCtl.Delete)
	}
}
