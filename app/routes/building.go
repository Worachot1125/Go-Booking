package routes

import (
	"app/app/controller"
	"app/app/middleware"

	"github.com/gin-gonic/gin"
)

func Building(router *gin.RouterGroup) {
	// Get the *bun.DB instance from config
	ctl := controller.New() // Pass the *bun.DB to the controller
	building := router.Group("")
	building.Use(middleware.AuthMiddleware())
	{
		building.POST("/create", ctl.BuildingCtl.Create)
		building.PATCH("/:id", ctl.BuildingCtl.Update)
		building.GET("/list", ctl.BuildingCtl.List)
		building.GET("/:id", ctl.BuildingCtl.Get)
		building.DELETE("/:id", ctl.BuildingCtl.Delete)
	}
}
