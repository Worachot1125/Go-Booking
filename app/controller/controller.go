package controller

import (
	"app/app/controller/position"
	"app/app/controller/product"
	"app/app/controller/role"
	"app/app/controller/room"
	"app/app/controller/user"
	"app/config"
)

type Controller struct {
	ProductCtl  *product.Controller
	UserCtl     *user.Controller
	RoomCtl     *room.Controller
	PositionCtl *position.Controller
	RoleCtl     *role.Controller // Assuming RoomCtl is also a product controller for this example

	// Other controllers...
}

func New() *Controller {
	db := config.GetDB()
	return &Controller{

		ProductCtl:  product.NewController(db),
		UserCtl:     user.NewController(db),
		RoomCtl:     room.NewController(db),
		PositionCtl: position.NewController(db),
		RoleCtl:     role.NewController(db), // Assuming RoomCtl is also a product controller for this example

		// Other controllers...
	}
}
