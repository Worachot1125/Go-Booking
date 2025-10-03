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
    name := ctx.PostForm("name")
    description := ctx.PostForm("description")
    capacityStr := ctx.PostForm("capacity")
    roomType := ctx.PostForm("room_type_id")
    startRoomStr := ctx.PostForm("start_room")
    endRoomStr := ctx.PostForm("end_room")

    if name == "" || description == "" || capacityStr == "" || roomType == "" || startRoomStr == "" || endRoomStr == "" {
        response.BadRequest(ctx, "ข้อมูลห้องไม่ครบ")
        return
    }

    capacity, err := strconv.Atoi(capacityStr)
    if err != nil {
        response.BadRequest(ctx, "จำนวนคนต้องเป็นตัวเลข")
        return
    }

    startRoom, err := strconv.ParseInt(startRoomStr, 10, 64)
    if err != nil {
        response.BadRequest(ctx, "ค่า start_room ต้องเป็น Unix timestamp")
        return
    }

    endRoom, err := strconv.ParseInt(endRoomStr, 10, 64)
    if err != nil {
        response.BadRequest(ctx, "ค่า end_room ต้องเป็น Unix timestamp")
        return
    }

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

    req := request.CreateRoom{
        Name:        name,
        Description: description,
        Capacity:    int64(capacity),
        RoomTypeID:  roomType,
        Image_url:   imageURL,
        StartRoom:   startRoom,
        EndRoom:     endRoom,
		Is_Available: true, // ห้องใหม่จะพร้อมใช้งานเสมอ
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
	ID := request.GetByIdRoom{}
	if err := ctx.BindUri(&ID); err != nil {
		logger.Err(err.Error())
		response.BadRequest(ctx, err.Error())
		return
	}

	name := ctx.PostForm("name")
	description := ctx.PostForm("description")
	capacityStr := ctx.PostForm("capacity")
	existingImageURL := ctx.PostForm("existing_image_url")
	isAvailableStr := ctx.PostForm("is_available")
	maintenanceNote := ctx.PostForm("maintenance_note")
	maintenanceETA := ctx.PostForm("maintenance_eta")

	var capacity int64
	if capacityStr != "" {
		c, err := strconv.Atoi(capacityStr)
		if err != nil {
			response.BadRequest(ctx, "จำนวนคนต้องเป็นตัวเลข")
			return
		}
		capacity = int64(c)
	}

	isAvailable := true // default
    if isAvailableStr != "" {
        isAvailable = isAvailableStr == "true" || isAvailableStr == "1"
    }
	isAvailablePtr := &isAvailable // สร้าง pointer

	imageURL := existingImageURL
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

	// เตรียม request แบบยืดหยุ่น
	req := request.UpdateRoom{
		CreateRoom: request.CreateRoom{
			Name:        name,
			Description: description,
			Capacity:    capacity, // จะเป็น 0 ถ้าไม่ได้กรอก
			Image_url:   imageURL,
			MaintenanceNote: maintenanceNote,
			MaintenanceETA:  maintenanceETA,
		},
		Is_Available: isAvailablePtr,
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
	if err := ctx.BindQuery(&req); err != nil {
		logger.Err(err.Error())
		response.BadRequest(ctx, err.Error())
		return
	}

	data, pagination, err := ctl.Service.List(ctx, req)
	if err != nil {
		logger.Err(err.Error())
		response.InternalError(ctx, "Failed to get room list")
		return
	}

	response.Success(ctx, gin.H{
		"data":       data,
		"pagination": pagination,
	})
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
