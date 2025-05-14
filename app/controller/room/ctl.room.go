package room

import (
	"app/app/helper"
	"app/app/request"
	"app/app/response"
	"app/internal/logger"
	"context"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"

	"github.com/gin-gonic/gin"
)

func (ctl *Controller) Create(ctx *gin.Context) {
	// อ่านฟิลด์ทีละตัว ไม่ใช้ ShouldBind
	name := ctx.PostForm("name")
	description := ctx.PostForm("description")
	capacityStr := ctx.PostForm("capacity")

	if name == "" || description == "" || capacityStr == "" {
		response.BadRequest(ctx, "ข้อมูลห้องไม่ครบ")
		return
	}

	// แปลง capacity จาก string เป็น int
	capacity, err := strconv.Atoi(capacityStr)
	if err != nil {
		response.BadRequest(ctx, "จำนวนคนต้องเป็นตัวเลข")
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

	// เตรียมสร้างข้อมูลห้อง
	req := request.CreateRoom{
		Name:        name,
		Description: description,
		Capacity:    int64(capacity),
		Image_url:   imageURL,
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
	// ดึง ID จาก path
	ID := request.GetByIdRoom{}
	if err := ctx.BindUri(&ID); err != nil {
		logger.Err(err.Error())
		response.BadRequest(ctx, err.Error())
		return
	}

	// อ่านค่าจาก multipart/form-data
	name := ctx.PostForm("name")
	description := ctx.PostForm("description")
	capacityStr := ctx.PostForm("capacity")
	existingImageURL := ctx.PostForm("existing_image_url")

	if name == "" || description == "" || capacityStr == "" {
		response.BadRequest(ctx, "ข้อมูลห้องไม่ครบ")
		return
	}

	// แปลง string → int
	capacity, err := strconv.Atoi(capacityStr)
	if err != nil {
		response.BadRequest(ctx, "จำนวนคนต้องเป็นตัวเลข")
		return
	}

	// เตรียมค่า image_url (ใช้ค่าที่มีอยู่)
	imageURL := existingImageURL

	// ถ้ามีไฟล์รูปใหม่ → อัปโหลด Cloudinary
	file, err := ctx.FormFile("image_url")
	if err == nil {
		src, err := file.Open()
		if err != nil {
			logger.Errf("Cannot open uploaded file: %v", err)
			response.InternalError(ctx, "ไม่สามารถเปิดไฟล์ได้")
			return
		}
		defer src.Close()

		imageURL, err = helper.UploadImageToCloudinary(src)

		if err != nil {
			logger.Errf("Upload to Cloudinary failed: %v", err)
			response.InternalError(ctx, "ไม่สามารถอัปโหลดรูปภาพได้")
			return
		}

	}

	// ส่งไป service
	req := request.UpdateRoom{
		CreateRoom: request.CreateRoom{
			Name:        name,
			Description: description,
			Capacity:    int64(capacity),
			Image_url:   imageURL,
		},
	}

	_, mserr, err := ctl.Service.Update(ctx, req, ID)
	if err != nil {
		ms := "Internal Server Error"
		if mserr {
			ms = err.Error()
		}
		logger.Err(err.Error())
		response.InternalError(ctx, ms)
		return
	}

	response.Success(ctx, nil)
}

func (ctl *Controller) List(ctx *gin.Context) {
	req := request.ListRoom{}
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

	data, total, err := ctl.Service.List(ctx, req)
	if err != nil {
		logger.Errf(err.Error())
		response.InternalError(ctx, err.Error())
		return
	}
	response.SuccessWithPaginate(ctx, data, req.Size, req.Page, total)

}

func (ctl *Controller) Get(ctx *gin.Context) {
	ID := request.GetByIdRoom{}
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
	ID := request.GetByIdRoom{}
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

func (ctl *Controller) UploadImage(ctx *gin.Context) {
	file, err := ctx.FormFile("image_url")
	if err != nil {
		logger.Errf("No file uploaded: %v", err)
		response.BadRequest(ctx, "กรุณาเลือกไฟล์รูปภาพ")
		return
	}

	// เปิดไฟล์
	src, err := file.Open()
	if err != nil {
		logger.Errf("Cannot open uploaded file: %v", err)
		response.InternalError(ctx, "ไม่สามารถเปิดไฟล์ได้")
		return
	}
	defer src.Close()

	// สร้าง Cloudinary client
	cld, err := cloudinary.NewFromParams(
		os.Getenv("CLOUDINARY_CLOUD_NAME"),
		os.Getenv("CLOUDINARY_API_KEY"),
		os.Getenv("CLOUDINARY_API_SECRET"),
	)
	if err != nil {
		logger.Errf("Cloudinary config error: %v", err)
		response.InternalError(ctx, "การตั้งค่า Cloudinary ไม่ถูกต้อง")
		return
	}

	// กำหนด timeout สำหรับ upload
	uploadCtx, cancel := context.WithTimeout(ctx, 15*time.Second)
	defer cancel()

	// อัปโหลดไปยัง Cloudinary
	uploadResult, err := cld.Upload.Upload(uploadCtx, src, uploader.UploadParams{
		Folder:   "room",
		PublicID: fmt.Sprintf("room_%d", time.Now().UnixNano()), // ตั้งชื่อให้ unique
	})
	if err != nil {
		logger.Errf("Upload to Cloudinary failed: %v", err)
		response.InternalError(ctx, "ไม่สามารถอัปโหลดรูปภาพได้")
		return
	}

	// ส่งกลับ URL
	response.Success(ctx, gin.H{
		"url": uploadResult.SecureURL,
	})
}
