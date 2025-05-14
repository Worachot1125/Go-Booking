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
	ex, err := s.db.NewSelect().Table("rooms").Where("id = ?", id.ID).Exists(ctx)
	if err != nil {
		return nil, false, err
	}
	if !ex {
		return nil, false, err
	}

	m := &model.Room{
		ID:          id.ID,
		Name:        req.Name,
		Description: req.Description,
		Capacity:    req.Capacity,
		Image_url:   req.Image_url,
	}

	m.SetUpdateNow()

	_, err = s.db.NewUpdate().Model(m).
		Set("name = ?name").
		Set("description = ?description").
		Set("capacity = ?capacity").
		Set("image_url = ?image_url").
		Set("updated_at = ?updated_at").
		WherePK().
		OmitZero().
		Returning("*").
		Exec(ctx)
	if err != nil {
		if strings.Contains(err.Error(), "duplicate key value") {
			return nil, true, errors.New("room already exists")
		}
	}

	return m, false, err
}

func (s *Service) List(ctx context.Context, req request.ListRoom) ([]response.RooomResponse, int, error) {
	offset := (req.Page - 1) * req.Size
	m := []response.RooomResponse{}

	query := s.db.NewSelect().
		TableExpr("rooms as r").
		ColumnExpr("r.id AS id").
		ColumnExpr("r.name AS name").
		ColumnExpr("r.description AS description").
		ColumnExpr("r.capacity AS capacity").
		ColumnExpr("r.image_url AS image_url").
		ColumnExpr("r.created_at AS created_at").
		ColumnExpr("r.updated_at AS updated_at").
		Where("r.deleted_at IS NULL").
		OrderExpr("r.name ASC")

	if req.Search != "" {
		searchBy := strings.ToLower(req.SearchBy)
		search := req.Search

		switch searchBy {
		case "name", "description":
			query = query.Where(fmt.Sprintf("LOWER(r.%s) LIKE ?", searchBy), "%"+strings.ToLower(search)+"%")
		case "capacity":
			// แปลง string -> int
			if capValue, err := strconv.Atoi(search); err == nil {
				query = query.Where("r.capacity = ?", capValue)
			}
		default:
			// fallback ถ้าไม่กำหนดหรือไม่รองรับ
			query = query.Where("LOWER(r.name) LIKE ?", "%"+strings.ToLower(search)+"%")
		}
	}

	// Count total before pagination
	count, err := query.Count(ctx)
	if err != nil {
		return nil, 0, err
	}

	// Order handling
	order := fmt.Sprintf("r.%s %s", req.SortBy, req.OrderBy)

	// Final query with order + pagination
	err = query.Order(order).Limit(req.Size).Offset(offset).Scan(ctx, &m)
	if err != nil {
		return nil, 0, err
	}

	return m, count, nil
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
		ColumnExpr("r.created_at AS created_at").
		ColumnExpr("r.updated_at AS updated_at").
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
