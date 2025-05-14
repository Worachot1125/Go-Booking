package permission

import (
	"app/app/request"
	"app/app/response"
	"app/internal/logger"

	"github.com/gin-gonic/gin"
)

func (ctl *Controller) Create(ctx *gin.Context) {
	per := request.PermissionCreate{}

	if err := ctx.Bind(&per); err != nil {
		response.BadRequest(ctx,err.Error())
		return 
	}
	
	_, mserr, err := ctl.Service.Create(ctx,per)
	if err != nil{
		ms := "internal server error"
		if mserr{
			ms = err.Error()
		}
		logger.Errf(err.Error())
		response.InternalError(ctx, ms)
		return
	}

	response.Success(ctx, nil)
}

func (ctl *Controller) Update(ctx *gin.Context) {
	ID := request.PermissionGetByID{}
	if err := ctx.BindUri(&ID); err != nil {
		logger.Err(err.Error())
		response.BadRequest(ctx, err.Error())
		return
	}

	per := request.PermissionUpdate{}
	if err := ctx.Bind(&per); err != nil {
		logger.Err(err.Error())
		response.BadRequest(ctx, err.Error())
		return
	}
		
	_, mserr, err := ctl.Service.Update(ctx, per, ID)
	if err != nil{
		ms := "internal server error"
		if mserr{
			ms = err.Error()
		}
		logger.Errf(err.Error())
		response.InternalError(ctx, ms)
		return
	}
	response.Success(ctx, nil)
}

func (ctl *Controller) Delete(ctx *gin.Context) {
	id := request.PermissionGetByID{}
	if err := ctx.BindUri(&id); err != nil {
		logger.Err(err.Error())
		response.BadRequest(ctx, err.Error())
		return
	}

	data, mserr, err := ctl.Service.Delete(ctx, id.ID)
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