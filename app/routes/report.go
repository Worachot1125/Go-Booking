package routes

import (
	"app/app/controller"
	"app/app/middleware"

	"github.com/gin-gonic/gin"
)

func Report(router *gin.RouterGroup) {
	// Get the *bun.DB instance from config
	ctl := controller.New() // Pass the *bun.DB to the controller
	//md := middleware.AuthMiddleware()
	log := middleware.NewLogResponse()
	report := router.Group("", log)
	{
		report.POST("/create", ctl.ReportCtl.Create)
		report.GET("/list", ctl.ReportCtl.List)
		report.GET("/:id", ctl.ReportCtl.Get)
		report.PATCH("/:id", ctl.ReportCtl.Update)
		report.DELETE("/:id", ctl.ReportCtl.Delete)
	}
}
