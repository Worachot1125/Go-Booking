package equipment

import (
	"app/app/helper"
	"app/app/request"
	"app/app/response"
	"strconv"

	"github.com/gin-gonic/gin"
)

func (ctl *Controller) Create(ctx *gin.Context) {
	file, err := ctx.FormFile("image_url")
	if err != nil {
		response.BadRequest(ctx, "กรุณาเลือกไฟล์รูปภาพ")
		return
	}

	src, err := file.Open()
	if err != nil {
		response.InternalError(ctx, "ไม่สามารถเปิดไฟล์ได้")
		return
	}
	defer src.Close()

	imageURL, err := helper.UploadImageToCloudinary(src)
	if err != nil {
		response.InternalError(ctx, "ไม่สามารถอัปโหลดรูปภาพได้")
		return
	}

	name := ctx.PostForm("name")
	quantity, _ := strconv.Atoi(ctx.PostForm("quantity"))

	req := request.CreateEquipment{
		Name:               name,
		Image_URL:          imageURL,
		Quantity:           quantity,
	}

	data, _, err := ctl.Service.Create(ctx, req)
	if err != nil {
		response.InternalError(ctx, err.Error())
		return
	}

	res := response.EquipmentResponse{
		ID:                 data.ID,
		Name:               data.Name,
		Image_URL:          data.Image_URL,
		Quantity:           data.Quantity,
		CreatedAt:          data.CreatedAt,
		UpdatedAt:          data.UpdatedAt,
	}

	response.Success(ctx, res)
}

func (ctl *Controller) Update(ctx *gin.Context) {
    ID := request.GetByIdEquipment{}
    if err := ctx.BindUri(&ID); err != nil {
        response.BadRequest(ctx, "invalid id")
        return
    }

    var req request.UpdateEquipment

    // อ่านค่าจาก form-data
    if name := ctx.PostForm("name"); name != "" {
        req.Name = &name
    }
    if quantity := ctx.PostForm("quantity"); quantity != "" {
        if q, err := strconv.Atoi(quantity); err == nil {
            req.Quantity = &q
        }
    }
    if availableQuantity := ctx.PostForm("available_quantity"); availableQuantity != "" {
        if aq, err := strconv.Atoi(availableQuantity); err == nil {
            req.Available_Quantity = &aq
        }
    }
    if status := ctx.PostForm("status"); status != "" {
        req.Status = &status
    }

    // อัปโหลดไฟล์ใหม่ถ้ามี
    file, err := ctx.FormFile("image_url")
    if err == nil {
        src, err := file.Open()
        if err != nil {
            response.InternalError(ctx, "ไม่สามารถเปิดไฟล์ได้")
            return
        }
        defer src.Close()
        imageURL, err := helper.UploadImageToCloudinary(src)
        if err != nil {
            response.InternalError(ctx, "ไม่สามารถอัปโหลดรูปภาพได้")
            return
        }
        req.Image_URL = &imageURL
    }

    data, _, err := ctl.Service.Update(ctx, req, ID.ID)
    if err != nil {
        response.InternalError(ctx, err.Error())
        return
    }

    res := response.EquipmentResponse{
        ID:                 data.ID,
        Name:               data.Name,
        Image_URL:          data.Image_URL,
        Quantity:           data.Quantity,
        CreatedAt:          data.CreatedAt,
        UpdatedAt:          data.UpdatedAt,
    }

    response.Success(ctx, res)
}
func (ctl *Controller) Delete(ctx *gin.Context) {
	ID := request.GetByIdEquipment{}
	if err := ctx.BindUri(&ID); err != nil {
		response.BadRequest(ctx, err.Error())
		return
	}

	data, _, err := ctl.Service.Delete(ctx, ID.ID)
	if err != nil {
		response.InternalError(ctx, err.Error())
		return
	}

	response.Success(ctx, gin.H{"id": data.ID})
}

func (ctl *Controller) Get(ctx *gin.Context) {
	ID := request.GetByIdEquipment{}
	if err := ctx.BindUri(&ID); err != nil {
		response.BadRequest(ctx, err.Error())
		return
	}

	data, err := ctl.Service.Get(ctx, ID.ID)
	if err != nil {
		response.InternalError(ctx, err.Error())
		return
	}

	res := response.EquipmentResponse{
		ID:                 data.ID,
		Name:               data.Name,
		Image_URL:          data.Image_URL,
		Quantity:           data.Quantity,
		CreatedAt:          data.CreatedAt,
		UpdatedAt:          data.UpdatedAt,
	}

	response.Success(ctx, res)
}

func (ctl *Controller) List(ctx *gin.Context) {
	req := request.ListEquipment{}
	if err := ctx.Bind(&req); err != nil {
		response.BadRequest(ctx, err.Error())
		return
	}

	if req.Page == 0 {
		req.Page = 1
	}
	if req.Size == 0 {
		req.Size = 10
	}
	if req.SortBy == "" {
		req.SortBy = "created_at"
	}
	if req.OrderBy == "" {
		req.OrderBy = "asc"
	}

	data, count, err := ctl.Service.List(ctx, req)
	if err != nil {
		response.InternalError(ctx, err.Error())
		return
	}

	res := make([]response.EquipmentResponse, len(data))
	for i, d := range data {
		res[i] = response.EquipmentResponse{
			ID:                 d.ID,
			Name:               d.Name,
			Image_URL:          d.Image_URL,
			Quantity:           d.Quantity,
			CreatedAt:          d.CreatedAt,
			UpdatedAt:          d.UpdatedAt,
		}
	}

	response.SuccessWithPaginate(ctx, res, req.Size, req.Page, count)
}
