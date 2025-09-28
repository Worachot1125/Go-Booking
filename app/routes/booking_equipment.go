package routes

import (
	"app/app/controller"

	"github.com/gin-gonic/gin"
)

func Booking_Equipment(router *gin.RouterGroup) {
	// Get the *bun.DB instance from config
	ctl := controller.New() // Pass the *bun.DB to the controller
	//md := middleware.AuthMiddleware()
	booking_equipment := router.Group("")

	{
		booking_equipment.GET("/list", ctl.BookingEquipmentCtl.List)
		booking_equipment.GET("/:id", ctl.BookingEquipmentCtl.Get)
		booking_equipment.POST("/create", ctl.BookingEquipmentCtl.Create)
		booking_equipment.PATCH("/:id", ctl.BookingEquipmentCtl.Update)
		booking_equipment.DELETE("/:id", ctl.BookingEquipmentCtl.Delete)
	}
}
