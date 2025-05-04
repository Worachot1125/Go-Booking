package user_role

import (
	"app/app/request"
	"app/app/response"
	"app/internal/logger"
	"log"

	"github.com/gin-gonic/gin"
)

func (ctl *Controller) Create(ctx *gin.Context) {
	rp := request.User_RoleCreate{}
	if err := ctx.Bind(&rp); err != nil {
		response.BadRequest(ctx, err.Error())
		return
	}

	data, mserr, err := ctl.Service.Create(ctx, rp)
	if err != nil {
		ms := "Internal Server Error"
		if mserr {
			ms = err.Error()
		}
		logger.Err(err.Error())
		response.InternalError(ctx, ms)
		return
	}

	response.Success(ctx, data)
}

func (ctl *Controller) GetByUserID(ctx *gin.Context) {
	var req request.GetByIDUser
	if err := ctx.ShouldBindUri(&req); err != nil {
		log.Println("Error binding user_id:", err)
		response.BadRequest(ctx, "Invalid input")
		return
	}

	data, err := ctl.Service.GetUserRolesByUserID(ctx, req.ID)
	if err != nil {
		log.Println("Error fetching user role:", err)
		response.InternalError(ctx, "Failed to fetch user role")
		return
	}

	response.Success(ctx, data)
}
