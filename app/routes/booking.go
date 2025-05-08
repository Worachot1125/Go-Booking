package routes

import (
	"app/app/controller"
	"app/app/middleware"

	"github.com/gin-gonic/gin"
)

func Booking(router *gin.RouterGroup) {
	// Get the *bun.DB instance from config
	ctl := controller.New() // Pass the *bun.DB to the controller
	md := middleware.AuthMiddleware()
	booking := router.Group("")
	{
		booking.GET("/list", ctl.BookingCtl.List)
		booking.GET("/:id", ctl.BookingCtl.Get)
		booking.GET("/user/:id", md, ctl.BookingCtl.GetBookingByUserID)
		booking.GET("/history/list", ctl.BookingCtl.ListHistory)
		booking.GET("/history/:id", md, ctl.BookingCtl.GetBookingHistoryByUserID)
		booking.POST("/create", md, ctl.BookingCtl.Create)
		booking.PATCH("/:id", md, ctl.BookingCtl.Update)
		booking.GET("/room/:id", md, ctl.BookingCtl.GetByRoomId)
		booking.DELETE("/:id", md, ctl.BookingCtl.Delete)
	}
}
