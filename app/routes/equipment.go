package routes

import (
	"app/app/controller"
	"app/app/middleware"

	"github.com/gin-gonic/gin"
)

func Equipment(router *gin.RouterGroup) {
	// Get the *bun.DB instance from config
	ctl := controller.New() // Pass the *bun.DB to the controller
	md := middleware.AuthMiddleware()
	equipment := router.Group("")

	{
		equipment.GET("/list", ctl.EquipmentCtl.List)
		equipment.GET("/:id", ctl.EquipmentCtl.Get)
		equipment.POST("/create", md, ctl.EquipmentCtl.Create)
		equipment.PATCH("/:id", md, ctl.EquipmentCtl.Update)
		equipment.DELETE("/:id", md, ctl.EquipmentCtl.Delete)
	}
}
