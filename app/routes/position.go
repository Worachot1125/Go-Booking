package routes

import (
	"app/app/controller"
	"app/app/middleware"

	"github.com/gin-gonic/gin"
)

func Position(router *gin.RouterGroup) {
	// Get the *bun.DB instance from config
	ctl := controller.New() // Pass the *bun.DB to the controller
	md := middleware.AuthMiddleware()
	position := router.Group("")
	{
		position.GET("/list", ctl.PositionCtl.List)
		position.GET("/:id", ctl.PositionCtl.Get)
		position.POST("/create", md, ctl.PositionCtl.Create)
		position.PATCH("/:id", md, ctl.PositionCtl.Update)
		position.DELETE("/:id", md, ctl.PositionCtl.Delete)

	}
}
