package controller

import (
	"app/app/controller/product"
	"app/app/controller/user"
	"app/app/controller/room"
	"app/config"
)

type Controller struct {
	ProductCtl *product.Controller
	UserCtl *user.Controller
	RoomCtl *room.Controller // Assuming RoomCtl is also a product controller for this example

	// Other controllers...
}

func New() *Controller {
	db := config.GetDB()
	return &Controller{

		ProductCtl: product.NewController(db),
		UserCtl: user.NewController(db),
		RoomCtl: room.NewController(db), // Assuming RoomCtl is also a product controller for this example

		// Other controllers...
	}
}
