package routes

import (
	"app/app/controller"
	"app/app/middleware"

	"github.com/gin-gonic/gin"
)

func Building_Room(router *gin.RouterGroup) {
	// Get the *bun.DB instance from config
	ctl := controller.New() // Pass the *bun.DB to the controller
	md := middleware.AuthMiddleware()
	building_room := router.Group("")

	{
		building_room.GET("/list", ctl.Building_RoomCtl.List)
		building_room.GET("/:id", ctl.Building_RoomCtl.Get)
		building_room.GET("/room/:id", ctl.Building_RoomCtl.GetByIDroom)
		building_room.GET("/buildRoom/:id", ctl.Building_RoomCtl.GetRoomsByBuildingID)
		building_room.POST("/create", md, ctl.Building_RoomCtl.Create)
		building_room.PATCH("/:id", md, ctl.Building_RoomCtl.Update)
		building_room.DELETE("/:id", md, ctl.Building_RoomCtl.Delete)
	}
}
