package routes

import (
	"app/app/controller"
	"app/app/middleware"

	"github.com/gin-gonic/gin"
)

func Reviews(router *gin.RouterGroup) {
	// Get the *bun.DB instance from config
	ctl := controller.New() // Pass the *bun.DB to the controller
	//md := middleware.AuthMiddleware()
	log := middleware.NewLogResponse()
	reviews := router.Group("", log)
	{
		reviews.POST("/create", ctl.ReviewsCtl.Create)
		reviews.GET("/list", ctl.ReviewsCtl.List)
		reviews.GET("/:id", ctl.ReviewsCtl.Get)
		reviews.GET("/bookings/:id", ctl.ReviewsCtl.GetByBookingID)
		reviews.PATCH("/:id", ctl.ReviewsCtl.Update)
		reviews.DELETE("/:id", ctl.ReviewsCtl.Delete)
	}
}
