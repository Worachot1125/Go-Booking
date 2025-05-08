package routes

import (
	"app/app/controller"
	"app/app/middleware"

	"github.com/gin-gonic/gin"
)

func Role_Permission(router *gin.RouterGroup) {
	ctl := controller.New()
	md := middleware.AuthMiddleware()
	role_permission := router.Group("")
	{
		role_permission.POST("/create", md, ctl.Role_Permissionctl.Create)
		role_permission.PATCH("/:id", md, ctl.Role_Permissionctl.Update)
	}
}
