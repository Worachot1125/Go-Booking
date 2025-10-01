package reportss

import (
	"app/app/request"
	"app/app/response"
	"app/internal/logger"

	"github.com/gin-gonic/gin"
)

func (ctl *Controller) Create(ctx *gin.Context) {
	req := request.CreateReport{}
	if err := ctx.ShouldBind(&req); err != nil {
		logger.Err(err.Error())
		response.BadRequest(ctx, err.Error())
		return
	}

	data, err := ctl.Service.Create(ctx, req)
	if err != nil {
		logger.Err(err.Error())
		response.InternalError(ctx, "Internal Server Error")
		return
	}

	response.Success(ctx, data)
}

func (ctl *Controller) Update(ctx *gin.Context) {
	ID := request.GetByIDReport{}
	if err := ctx.ShouldBindUri(&ID); err != nil {
		logger.Err(err.Error())
		response.BadRequest(ctx, err.Error())
		return
	}

	body := request.UpdateReport{}
	if err := ctx.ShouldBind(&body); err != nil {
		logger.Err(err.Error())
		response.BadRequest(ctx, err.Error())
		return
	}

	_, err := ctl.Service.Update(ctx, body, ID)
	if err != nil {
		logger.Err(err.Error())
		response.InternalError(ctx, "Internal Server Error")
		return
	}

	response.Success(ctx, nil)
}

func (ctl *Controller) List(ctx *gin.Context) {
	req := request.ListReport{}
	if err := ctx.ShouldBind(&req); err != nil {
		logger.Err(err.Error())
		response.BadRequest(ctx, err.Error())
		return
	}

	if req.Page == 0 {
		req.Page = 1
	}
	if req.Size == 0 {
		req.Size = 10
	}
	if req.OrderBy == "" {
		req.OrderBy = "asc"
	}
	if req.SortBy == "" {
		req.SortBy = "created_at"
	}

	data, total, err := ctl.Service.List(ctx, req)
	if err != nil {
		logger.Err(err.Error())
		response.InternalError(ctx, "Internal Server Error")
		return
	}
	response.SuccessWithPaginate(ctx, data, req.Size, req.Page, total)
}

func (ctl *Controller) Get(ctx *gin.Context) {
	ID := request.GetByIDReport{}
	if err := ctx.ShouldBindUri(&ID); err != nil {
		logger.Err(err.Error())
		response.BadRequest(ctx, err.Error())
		return
	}

	data, err := ctl.Service.Get(ctx, ID)
	if err != nil {
		logger.Err(err.Error())
		response.InternalError(ctx, "Internal Server Error")
		return
	}

	if data == nil {
		response.NotFound(ctx, "Report not found")
		return
	}

	response.Success(ctx, data)
}

func (ctl *Controller) Delete(ctx *gin.Context) {
	ID := request.GetByIDReport{}
	if err := ctx.ShouldBindUri(&ID); err != nil {
		logger.Err(err.Error())
		response.BadRequest(ctx, err.Error())
		return
	}

	if err := ctl.Service.Delete(ctx, ID); err != nil {
		logger.Err(err.Error())
		response.InternalError(ctx, "Internal Server Error")
		return
	}

	response.Success(ctx, nil)
}
