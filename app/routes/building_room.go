package routes

// import (
// 	"app/app/controller"

// 	"github.com/gin-gonic/gin"
// )

// func Building_Room(router *gin.RouterGroup) {
// 	// Get the *bun.DB instance from config
// 	ctl := controller.New() // Pass the *bun.DB to the controller
// 	building_room := router.Group("")
// 	building_room.Use(middleware.AuthMiddleware())
// 	{
// 		building_room.POST("/create", ctl.Building_RoomCtl.Create)
// 	}
// }