package position

import (
	"app/app/request"
	"app/app/response"
	"app/internal/logger"

	"github.com/gin-gonic/gin"
)

func (ctl *Controller) Create(ctx *gin.Context) {
	req := request.PositionCreate{}
	if err := ctx.Bind(&req); err != nil {
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
	id := request.PositionGetByID{}
	if err := ctx.BindUri(&id); err != nil {
		logger.Err(err.Error())
		response.BadRequest(ctx, err.Error())
		return
	}

	req := request.PositionUpdate{}
	if err := ctx.Bind(&req); err != nil {
		logger.Err(err.Error())
		response.BadRequest(ctx, err.Error())
		return
	}

	data, mserr, err := ctl.Service.Update(ctx, id.ID, req)
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

func (ctl *Controller) Delete(ctx *gin.Context) {
	id := request.PositionGetByID{}
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

func (ctl *Controller) Get(ctx *gin.Context) {
	id := request.PositionGetByID{}
	if err := ctx.BindUri(&id); err != nil {
		logger.Err(err.Error())
		response.BadRequest(ctx, err.Error())
		return
	}

	data, err := ctl.Service.Get(ctx, id.ID)
	if err != nil {
		logger.Err(err.Error())
		response.InternalError(ctx, err.Error())
		return
	}

	response.Success(ctx, data)
}

func (ctl *Controller) List(ctx *gin.Context) {
	req := request.PositionListRequest{}
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

	data, count, err := ctl.Service.List(ctx, req)
	if err != nil {
		logger.Err(err.Error())
		response.InternalError(ctx, err.Error())
		return
	}

	// 🔽 เพิ่มส่วนนี้เพื่อแปลงข้อมูลก่อนส่งออก
	positions := make([]response.ListPosition, 0, len(data))
	for _, p := range data {
		positions = append(positions, response.ListPosition{
			ID:        p.ID,
			Name:      p.Name,
			CreatedAt: p.CreatedAt,
			UpdatedAt: p.UpdatedAt,
		})
	}

	response.SuccessWithPaginate(ctx, positions, req.Size, req.Page, count)
}
