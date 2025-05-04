package routes

import (
	"app/app/controller"

	"github.com/gin-gonic/gin"
)

func User_Role(router *gin.RouterGroup) {
	// Get the *bun.DB instance from config
	ctl := controller.New() // Pass the *bun.DB to the controller
	user_role := router.Group("")
	{
		user_role.GET("/:id", ctl.User_RoleCtl.GetByUserID)
	}
}
