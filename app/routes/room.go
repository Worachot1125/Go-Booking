package routes

import (
	"app/app/controller"
	"app/app/middleware"

	"github.com/gin-gonic/gin"
)

func Room(router *gin.RouterGroup) {
	// Get the *bun.DB instance from config
	ctl := controller.New() // Pass the *bun.DB to the controller
	room := router.Group("")
	room.Use(middleware.AuthMiddleware())
	{
		room.POST("/create", ctl.RoomCtl.Create)
		room.PATCH("/:id", ctl.RoomCtl.Update)
		room.GET("/list", ctl.RoomCtl.List)
		room.GET("/:id", ctl.RoomCtl.Get)
		room.DELETE("/:id", ctl.RoomCtl.Delete)
	}
}
