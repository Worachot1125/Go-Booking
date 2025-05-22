package user

import (
	"app/app/helper"
	"app/app/request"
	"app/app/response"
	"app/internal/logger"

	"github.com/gin-gonic/gin"
)

func (ctl *Controller) Create(ctx *gin.Context) {
	firstName := ctx.PostForm("first_name")
	lastName := ctx.PostForm("last_name")
	email := ctx.PostForm("email")
	password := ctx.PostForm("password")
	phone := ctx.PostForm("phone")
	positionName := ctx.PostForm("position_name")

	if firstName == "" || lastName == "" || email == "" || password == "" || phone == "" || positionName == "" {
		response.BadRequest(ctx, "ข้อมูลผู้ใช้ไม่ครบ")
		return
	}

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

	// เตรียม request
	req := request.CreateUser{
		FirstName:     firstName,
		LastName:      lastName,
		Email:         email,
		Password:      password,
		Phone:         phone,
		Position_Name: positionName,
		Image_url:     imageURL,
	}

	// เรียก Service
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
	// ดึง ID จาก path
	ID := request.GetByIdUser{}
	if err := ctx.BindUri(&ID); err != nil {
		logger.Err(err.Error())
		response.BadRequest(ctx, err.Error())
		return
	}

	// อ่านค่าทีละฟิลด์จาก multipart/form-data
	firstName := ctx.PostForm("first_name")
	lastName := ctx.PostForm("last_name")
	email := ctx.PostForm("email")
	password := ctx.PostForm("password")
	phone := ctx.PostForm("phone")
	positionName := ctx.PostForm("position_name")
	existingImageURL := ctx.PostForm("existing_image_url")

	if firstName == "" || lastName == "" || email == "" || phone == "" || positionName == "" {
		response.BadRequest(ctx, "ข้อมูลไม่ครบ")
		return
	}

	// เตรียม imageURL
	imageURL := existingImageURL
	file, err := ctx.FormFile("image_url")
	if err == nil {
		src, err := file.Open()
		if err != nil {
			logger.Errf("cannot open uploaded file: %v", err)
			response.InternalError(ctx, "ไม่สามารถเปิดไฟล์ได้")
			return
		}
		defer src.Close()

		imageURL, err = helper.UploadImageToCloudinary(src)
		if err != nil {
			logger.Errf("upload to cloudinary failed: %v", err)
			response.InternalError(ctx, "ไม่สามารถอัปโหลดรูปภาพได้")
			return
		}
	}

	// สร้าง request และเรียก service
	body := request.UpdateUser{
		CreateUser: request.CreateUser{
			FirstName:     firstName,
			LastName:      lastName,
			Email:         email,
			Password:      password,
			Phone:         phone,
			Position_Name: positionName,
			Image_url:     imageURL,
		},
	}

	_, mserr, err := ctl.Service.Update(ctx, body, ID)
	if err != nil {
		ms := "internal server error"
		if mserr {
			ms = err.Error()
		}
		logger.Errf(err.Error())
		response.InternalError(ctx, ms)
		return
	}

	response.Success(ctx, nil)
}

func (ctl *Controller) List(ctx *gin.Context) {
	req := request.ListUser{}
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
	ID := request.GetByIdUser{}
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

func (ctl *Controller) Delete(ctx *gin.Context) {
	ID := request.GetByIdUser{}
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
