package controller

import (
	"app/app/controller/buliding"
	"app/app/controller/buliding_room"
	"app/app/controller/product"
	"app/app/controller/room"
	"app/app/controller/user"
	"app/config"
)

type Controller struct {
	ProductCtl *product.Controller
	UserCtl *user.Controller
	RoomCtl *room.Controller 
	BuildingCtl *building.Controller
	Building_RoomCtl *building_room.Controller

	// Other controllers...
}

func New() *Controller {
	db := config.GetDB()
	return &Controller{

		ProductCtl: product.NewController(db),
		UserCtl: user.NewController(db),
		RoomCtl: room.NewController(db), 
		BuildingCtl: building.NewController(db),
		Building_RoomCtl: building_room.NewController(db),
		// Other controllers...
	}
}
