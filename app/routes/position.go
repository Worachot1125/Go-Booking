package routes

import (
	"app/app/controller"
	
	"github.com/gin-gonic/gin"
)

func Position(router *gin.RouterGroup) {
	// Get the *bun.DB instance from config
	ctl := controller.New() // Pass the *bun.DB to the controller
	position := router.Group("")
	{
		position.POST("/create", ctl.PositionCtl.Create)
		position.PATCH("/:id", ctl.PositionCtl.Update)
		position.GET("/list", ctl.PositionCtl.List)
		position.GET("/:id", ctl.PositionCtl.Get)
		position.DELETE("/:id", ctl.PositionCtl.Delete)

	}
}
