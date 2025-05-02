package controller

import (
	"app/app/controller/booking"
	"app/app/controller/building"
	"app/app/controller/building_room"
	"app/app/controller/login"
	"app/app/controller/permission"
	"app/app/controller/position"
	"app/app/controller/product"
	"app/app/controller/role"
	"app/app/controller/role_permission"
	"app/app/controller/room"
	"app/app/controller/user"
	"app/config"
)

type Controller struct {
	ProductCtl         *product.Controller
	UserCtl            *user.Controller
	RoomCtl            *room.Controller
	PositionCtl        *position.Controller
	RoleCtl            *role.Controller
	PermissionCtl      *permission.Controller
	Role_Permissionctl *role_permission.Controller
	BuildingCtl        *building.Controller
	Building_RoomCtl   *building_room.Controller
	BookingCtl         *booking.Controller
	LoginCtl           *login.Controller 

}

func New() *Controller {
	db := config.GetDB()
	return &Controller{

		ProductCtl:         product.NewController(db),
		UserCtl:            user.NewController(db),
		RoomCtl:            room.NewController(db),
		PositionCtl:        position.NewController(db),
		RoleCtl:            role.NewController(db),
		PermissionCtl:      permission.NewController(db),
		Role_Permissionctl: role_permission.NewController(db),
		BuildingCtl:        building.NewController(db),
		Building_RoomCtl:   building_room.NewController(db),
		LoginCtl:           login.NewController(db),
		BookingCtl:         booking.NewController(db), 

	}
}
