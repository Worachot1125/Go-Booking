package routes

import (
	"app/app/controller"
	"app/app/middleware"

	"github.com/gin-gonic/gin"
)

func Role(router *gin.RouterGroup) {
	ctl := controller.New()
	md := middleware.AuthMiddleware()
	role := router.Group("")
	{
		role.POST("/create", md, ctl.RoleCtl.Create)
		role.GET("/list", md, ctl.RoleCtl.List)
		role.DELETE("/:id", md, ctl.RoleCtl.Delete)
	}
}
