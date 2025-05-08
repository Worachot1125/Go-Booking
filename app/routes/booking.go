package routes

import (
	"app/app/controller"

	"github.com/gin-gonic/gin"
)

func Booking(router *gin.RouterGroup) {
	// Get the *bun.DB instance from config
	ctl := controller.New() // Pass the *bun.DB to the controller
	booking := router.Group("")
	{
		booking.GET("/list", ctl.BookingCtl.List)
		booking.GET("/:id", ctl.BookingCtl.Get)
		booking.GET("/user/:id", ctl.BookingCtl.GetBookingByUserID)
		booking.GET("/history/list", ctl.BookingCtl.ListHistory)
		booking.GET("/history/:id", ctl.BookingCtl.GetBookingHistoryByUserID)
	}
	{
		booking.POST("/create", ctl.BookingCtl.Create)
		booking.PATCH("/:id", ctl.BookingCtl.Update)
		booking.GET("/room/:id", ctl.BookingCtl.GetByRoomId)
		booking.DELETE("/:id", ctl.BookingCtl.Delete)
	}
}
