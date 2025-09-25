package building

import (
	"app/app/helper"
	"app/app/request"
	"app/app/response"
	"app/internal/logger"
	"context"
	"fmt"
	"os"
	"time"

	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
	"github.com/gin-gonic/gin"
)

func (ctl *Controller) Create(ctx *gin.Context) {
	// อ่านฟิลด์ทีละตัว ไม่ใช้ ShouldBind
	name := ctx.PostForm("name")
	existingImageURL := ctx.PostForm("image_url")

	if name == "" {
		response.BadRequest(ctx, "ชื่ออาคารไม่ได้ระบุ")
		return
	}

	// รับไฟล์รูปภาพ - ใช้ field name เดียวกับ Room
	imageURL := existingImageURL // ใช้ URL เดิมเป็นค่าเริ่มต้น
	file, err := ctx.FormFile("image_url")
	if err == nil {
		// มีไฟล์ใหม่ส่งมา
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
	} else {
		// แสดง log เพื่อ debug
		logger.Errf("No file uploaded or file field not found: %v", err)
	}

	// เตรียมสร้างข้อมูลอาคาร
	req := request.CreateBuilding{
		Name:      name,
		Image_url: imageURL,
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
	ID := request.GetByIdBuilding{}
	if err := ctx.BindUri(&ID); err != nil {
		logger.Err(err.Error())
		response.BadRequest(ctx, err.Error())
		return
	}

	name := ctx.PostForm("name")
	existingImageURL := ctx.PostForm("existing_image_url")

	// ใช้รูปเดิมเป็นค่าเริ่มต้น
	imageURL := existingImageURL
	file, err := ctx.FormFile("image_url")
	if err == nil {
		// มีไฟล์ใหม่ส่งมา
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

	// เตรียม request
	body := request.UpdateBuilding{
    CreateBuilding: request.CreateBuilding{
        Name:      name,
        Image_url: imageURL,
    },
}

	_, mserr, err := ctl.Service.Update(ctx, body, ID)
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
	req := request.ListBuilding{}
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
	ID := request.GetByIdBuilding{}
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
	ID := request.GetByIdBuilding{}
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
		Folder:   "building",
		PublicID: fmt.Sprintf("building_%d", time.Now().UnixNano()), // ตั้งชื่อให้ unique
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