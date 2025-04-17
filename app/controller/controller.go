package controller

import (
	"app/app/controller/product"
	"app/app/controller/user"
	"app/config"
)

type Controller struct {
	ProductCtl *product.Controller
	UserCtl *user.Controller

	// Other controllers...
}

func New() *Controller {
	db := config.GetDB()
	return &Controller{

		ProductCtl: product.NewController(db),
		UserCtl: user.NewController(db),

		// Other controllers...
	}
}
