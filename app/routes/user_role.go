package routes

import (
	"app/app/controller"
	"app/app/middleware"

	"github.com/gin-gonic/gin"
)

func User_Role(router *gin.RouterGroup) {
	ctl := controller.New()
	md := middleware.AuthMiddleware()
	user_role := router.Group("")
	{
		user_role.GET("/:id", md, ctl.User_RoleCtl.GetByUserID)
	}
}
