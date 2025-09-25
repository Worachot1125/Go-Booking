package room

import (
	"app/app/model"
	"app/app/request"
	"app/app/response"
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"
)

func (s *Service) Create(ctx context.Context, req request.CreateRoom) (*model.Room, bool, error) {
	// ตรวจสอบข้อมูลห้องก่อนที่จะบันทึก
	if req.Name == "" || req.Description == "" || req.Capacity == 0 || req.Image_url == "" {
		return nil, false, errors.New("ข้อมูลห้องไม่ครบถ้วน")
	}

	// สร้างโมเดลห้องใหม่
	m := model.Room{
		Name:        req.Name,
		Description: req.Description,
		Capacity:    req.Capacity,
		Image_url:   req.Image_url,
	}

	// ทำการบันทึกข้อมูลห้องลงในฐานข้อมูล
	_, err := s.db.NewInsert().Model(&m).Exec(ctx)
	if err != nil {
		if strings.Contains(err.Error(), "duplicate key value") {
			return nil, true, errors.New("room already exists")
		}
		return nil, false, err
	}

	return &m, false, nil
}

func (s *Service) Update(ctx context.Context, req request.UpdateRoom, id request.GetByIdRoom) (*model.Room, bool, error) {
	// ตรวจสอบว่ามีห้องนี้ไหม
	ex, err := s.db.NewSelect().Table("rooms").Where("id = ?", id.ID).Exists(ctx)
	if err != nil {
		return nil, false, err
	}
	if !ex {
		return nil, false, errors.New("room not found")
	}

	// เตรียม model สำหรับอัปเดต
	m := &model.Room{ID: id.ID}
	q := s.db.NewUpdate().Model(m).WherePK().OmitZero().Returning("*")

	// อัปเดตเฉพาะ field ที่ส่งมา
	if req.Name != "" {
		m.Name = req.Name
		q.Set("name = ?name")
	}
	if req.Description != "" {
		m.Description = req.Description
		q.Set("description = ?description")
	}
	if req.Capacity != 0 {
		m.Capacity = req.Capacity
		q.Set("capacity = ?capacity")
	}
	if req.Image_url != "" {
		m.Image_url = req.Image_url
		q.Set("image_url = ?image_url")
	}

	// อัปเดตเวลา
	m.SetUpdateNow()
	q.Set("updated_at = ?updated_at")

	// รัน query
	_, err = q.Exec(ctx)
	if err != nil {
		if strings.Contains(err.Error(), "duplicate key value") {
			return nil, true, errors.New("room already exists")
		}
		return nil, false, err
	}

	return m, false, nil
}

func (s *Service) List(ctx context.Context, req request.ListRoom) ([]response.RooomResponse, map[string]int, error) {
	// ตรวจสอบ Page / Size
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.Size <= 0 {
		req.Size = 10
	}
	offset := (req.Page - 1) * req.Size

	rooms := []response.RooomResponse{}

	// Base query
	query := s.db.NewSelect().
		TableExpr("rooms as r").
		ColumnExpr("r.id AS id").
		ColumnExpr("r.name AS name").
		ColumnExpr("r.description AS description").
		ColumnExpr("r.capacity AS capacity").
		ColumnExpr("r.image_url AS image_url").
		ColumnExpr("b.name AS building").
		ColumnExpr("r.created_at AS created_at").
		ColumnExpr("r.updated_at AS updated_at").
		Join("LEFT JOIN building_rooms as br ON r.id::uuid = br.room_id::uuid").
		Join("LEFT JOIN buildings as b ON br.building_id::uuid = b.id::uuid").
		Where("r.deleted_at IS NULL")

	// Search
	if req.Search != "" {
		searchBy := strings.ToLower(req.SearchBy)
		search := strings.ToLower(req.Search)

		switch searchBy {
		case "name", "description":
			query = query.Where(fmt.Sprintf("LOWER(r.%s) LIKE ?", searchBy), "%"+search+"%")
		case "capacity":
			if capVal, err := strconv.Atoi(search); err == nil {
				query = query.Where("r.capacity = ?", capVal)
			}
		default:
			query = query.Where("LOWER(r.name) LIKE ?", "%"+search+"%")
		}
	}

	// Count total
	count, err := query.Clone().Count(ctx)
	if err != nil {
		return nil, nil, err
	}

	// Sort
	sortBy := "name"
	if req.SortBy != "" {
		sortBy = req.SortBy
	}
	orderBy := "ASC"
	if strings.ToUpper(req.OrderBy) == "DESC" {
		orderBy = "DESC"
	}
	query = query.Order(fmt.Sprintf("r.%s %s", sortBy, orderBy))

	// Pagination
	err = query.Limit(req.Size).Offset(offset).Scan(ctx, &rooms)
	if err != nil {
		return nil, nil, err
	}

	// Pagination info
	pagination := map[string]int{
		"page":       req.Page,
		"size":       req.Size,
		"total":      count,
		"total_page": (count + req.Size - 1) / req.Size,
	}

	return rooms, pagination, nil
}


func (s *Service) Get(ctx context.Context, id request.GetByIdRoom) (*response.RooomResponse, error) {
	m := response.RooomResponse{}

	err := s.db.NewSelect().
		TableExpr("rooms as r").
		ColumnExpr("r.id AS id").
		ColumnExpr("r.name AS name").
		ColumnExpr("r.description AS description").
		ColumnExpr("r.capacity AS capacity").
		ColumnExpr("r.image_url AS image_url").
		ColumnExpr("b.name AS building").
		ColumnExpr("r.created_at AS created_at").
		ColumnExpr("r.updated_at AS updated_at").
		Join("JOIN building_rooms as br ON r.id::uuid = br.room_id::uuid ").
		Join("JOIN buildings as b ON br.building_id::uuid = b.id::uuid ").
		Where("r.deleted_at IS NULL").
		Where("r.id = ?", id.ID).
		Scan(ctx, &m)
	return &m, err
}

func (s *Service) Delete(ctx context.Context, id request.GetByIdRoom) error {
	ex, err := s.db.NewSelect().Table("rooms").Where("id = ?", id.ID).Where("deleted_at IS NULL").Exists(ctx)
	if err != nil {
		return err
	}
	if !ex {
		return errors.New("room not found")
	}

	_, err = s.db.NewDelete().Model((*model.Room)(nil)).Where("id = ?", id.ID).Exec(ctx)
	return err
}
