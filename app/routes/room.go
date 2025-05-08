package routes

import (
	"app/app/controller"
	"app/app/middleware"

	"github.com/gin-gonic/gin"
)

func Room(router *gin.RouterGroup) {
	ctl := controller.New()
	md := middleware.AuthMiddleware()
	room := router.Group("")
	{
		room.GET("/list", ctl.RoomCtl.List)
		room.GET("/:id", ctl.RoomCtl.Get)
		room.POST("/create", md, ctl.RoomCtl.Create)
		room.PATCH("/:id", md, ctl.RoomCtl.Update)
		room.DELETE("/:id", md, ctl.RoomCtl.Delete)
		room.POST("/upload", md, ctl.RoomCtl.UploadImage)
	}
}
