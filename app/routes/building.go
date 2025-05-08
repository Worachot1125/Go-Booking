package routes

import (
	"app/app/controller"
	"app/app/middleware"

	"github.com/gin-gonic/gin"
)

func Building(router *gin.RouterGroup) {
	// Get the *bun.DB instance from config
	ctl := controller.New() // Pass the *bun.DB to the controller
	md := middleware.AuthMiddleware()
	building := router.Group("")
	{
		building.GET("/list", ctl.BuildingCtl.List)
		building.GET("/:id", ctl.BuildingCtl.Get)
		building.POST("/create", md, ctl.BuildingCtl.Create)
		building.PATCH("/:id", md, ctl.BuildingCtl.Update)
		building.DELETE("/:id", md, ctl.BuildingCtl.Delete)
	}
}
