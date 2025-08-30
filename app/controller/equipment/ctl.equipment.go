package equipment

import (
	"app/app/helper"
	"app/app/request"
	"app/app/response"
	"app/internal/logger"

	"github.com/gin-gonic/gin"
)

func (ctl *Controller) Create(ctx *gin.Context) {

	// รับไฟล์รูปภาพ
	file, err := ctx.FormFile("image_url")
	if err != nil {
		logger.Errf("No file uploaded: %v", err)
		response.BadRequest(ctx, "กรุณาเลือกไฟล์รูปภาพ")
		return
	}

	src, err := file.Open()
	if err != nil {
		logger.Errf("Cannot open uploaded file: %v", err)
		response.InternalError(ctx, "ไม่สามารถเปิดไฟล์ได้")
		return
	}
	defer src.Close()

	imageURL, err := helper.UploadImageToCloudinary(src)

	if err != nil {
		logger.Errf("Upload to Cloudinary failed: %v", err)
		response.InternalError(ctx, "ไม่สามารถอัปโหลดรูปภาพได้")
		return
	}

	// เตรียมสร้างข้อมูลห้อง
	req := request.CreateEquipment{
		Image_URL: imageURL,

	}

	// เรียก Service.Create
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
	id := request.ProductGetByID{}
	if err := ctx.BindUri(&id); err != nil {
		logger.Err(err.Error())
		response.BadRequest(ctx, err.Error())
		return
	}

	req := request.UpdateEquipment{}
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
	id := request.ProductGetByID{}
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
	id := request.ProductGetByID{}
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
	req := request.ProductListReuest{}
	if err := ctx.Bind(&req); err != nil {
		logger.Err(err.Error())
		response.BadRequest(ctx, err.Error())
		return
	}

	if req.Page == 0 {
		req.Page = 1
	}

	if req.Page == 0 {
		req.Page = 10
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

	response.SuccessWithPaginate(ctx, data, req.Size, req.Page, count)
}
