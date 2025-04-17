package routes

import (
	"app/app/controller"

	"github.com/gin-gonic/gin"
)

func Building(router *gin.RouterGroup) {
	// Get the *bun.DB instance from config
	ctl := controller.New() // Pass the *bun.DB to the controller
	building := router.Group("")
	{
		building.POST("/create", ctl.RoomCtl.Create)
	}
}