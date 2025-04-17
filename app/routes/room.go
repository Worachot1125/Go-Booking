package routes

import (
	"app/app/controller"

	"github.com/gin-gonic/gin"
)

func Room(router *gin.RouterGroup) {
	// Get the *bun.DB instance from config
	ctl := controller.New() // Pass the *bun.DB to the controller
	room := router.Group("")
	{
		room.POST("/create", ctl.RoomCtl.Create)
	}
}
