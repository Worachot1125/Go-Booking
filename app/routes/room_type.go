package routes

import (
	"app/app/controller"

	"github.com/gin-gonic/gin"
)

func RoomType(router *gin.RouterGroup) {
	ctl := controller.New()
	//md := middleware.AuthMiddleware()
	room_type := router.Group("")
	{
		room_type.GET("/list", ctl.RoomTypeCtl.List)
		room_type.DELETE("/:id", ctl.RoomTypeCtl.Delete)
	}
}
