package reviewss

import (
	"app/app/request"
	"app/app/response"
	"app/internal/logger"

	"github.com/gin-gonic/gin"
)

func (ctl *Controller) Create(ctx *gin.Context) {
	req := request.CreateReviews{}
	if err := ctx.ShouldBind(&req); err != nil {
		response.BadRequest(ctx, err.Error())
		return
	}

	data, mserr, err := ctl.Service.Create(ctx, req)
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

func (ctl *Controller) Update(ctx *gin.Context) {
	ID := request.GetByIDReviews{}
	if err := ctx.BindUri(&ID); err != nil {
		logger.Err(err.Error())
		response.BadRequest(ctx, err.Error())
		return
	}
	body := request.UpdateReviews{}
	if err := ctx.Bind(&body); err != nil {
		logger.Err(err.Error())
		response.BadRequest(ctx, err.Error())
		return
	}

	_, _, err := ctl.Service.Update(ctx, body, ID)
	if err != nil {
		logger.Err(err.Error())
		response.InternalError(ctx, err.Error())
		return
	}

	response.Success(ctx, nil)
}

func (ctl *Controller) List(ctx *gin.Context) {
	req := request.ListReviews{}
	if err := ctx.Bind(&req); err != nil {
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
		logger.Errf(err.Error())
		response.InternalError(ctx, err.Error())
		return
	}
	response.SuccessWithPaginate(ctx, data, req.Size, req.Page, total)
}

func (ctl *Controller) Get(ctx *gin.Context) {
	ID := request.GetByIDReviews{}
	if err := ctx.BindUri(&ID); err != nil {
		logger.Err(err.Error())
		response.BadRequest(ctx, err.Error())
		return
	}

	data, err := ctl.Service.Get(ctx, ID)
	if err != nil {
		logger.Errf(err.Error())
		response.InternalError(ctx, err.Error())
		return
	}
	response.Success(ctx, data)
}

func (ctl *Controller) GetByBookingID(ctx *gin.Context) {
	req := request.GetByBookingIDReviews{}
	if err := ctx.BindUri(&req); err != nil {
		response.BadRequest(ctx, err.Error())
		return
	}

	data, err := ctl.Service.GetByBookingID(ctx, req)
	if err != nil {
		response.InternalError(ctx, err.Error())
		return
	}

	response.Success(ctx, data)
}

func (ctl *Controller) Delete(ctx *gin.Context) {
	ID := request.GetByIDReviews{}
	if err := ctx.BindUri(&ID); err != nil {
		logger.Err(err.Error())
		response.BadRequest(ctx, err.Error())
		return
	}

	err := ctl.Service.Delete(ctx, ID)
	if err != nil {
		logger.Errf(err.Error())
		response.InternalError(ctx, err.Error())
		return
	}
	response.Success(ctx, nil)
}
