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

func (s *Service) Create(ctx context.Context, req request.CreateRoom) (*response.RooomResponse, bool, error) {
	if req.Name == "" || req.Description == "" || req.Capacity == 0 || req.Image_url == "" {
		return nil, false, errors.New("ข้อมูลห้องไม่ครบถ้วน")
	}
	if req.RoomTypeID == "" {
		return nil, false, errors.New("room type is required")
	}

	// หลัง (ถูกจุดประสงค์: กันเฉพาะ "ห้องเดียวกัน")
	exists, err := s.db.NewSelect().
		Model((*model.Room)(nil)).
		Where("name = ?", req.Name).
		Where("start_room < ?", req.EndRoom).
		Where("end_room > ?", req.StartRoom).
		Where("deleted_at IS NULL").
		Exists(ctx)
		
	if err != nil {
		return nil, false, err
	}
	if exists {
		return nil, true, errors.New("room range conflict")
	}

	m := model.Room{
		RoomTypeID:       req.RoomTypeID,
		Name:             req.Name,
		Description:      req.Description,
		Capacity:         req.Capacity,
		Image_url:        req.Image_url,
		StartRoom:        req.StartRoom,
		EndRoom:          req.EndRoom,
		Maintenance_note: req.MaintenanceNote,
		Maintenance_eta:  req.MaintenanceETA,
	}

	_, err = s.db.NewInsert().Model(&m).Exec(ctx)
	if err != nil {
		if strings.Contains(err.Error(), "duplicate key value") {
			return nil, true, errors.New("room already exists")
		}
		return nil, false, err
	}

	// 🔑 map ไปเป็น response.RooomResponse
	res := response.RooomResponse{
		ID:              m.ID,
		RoomTypeID:      m.RoomTypeID,
		Name:            m.Name,
		Capacity:        m.Capacity,
		Description:     m.Description,
		ImageURL:        m.Image_url,
		StartRoom:       m.StartRoom,
		EndRoom:         m.EndRoom,
		Is_Available:    m.Is_Available,
		MaintenanceNote: m.Maintenance_note,
		MaintenanceETA:  m.Maintenance_eta,
		CreatedAt:       m.CreatedAt,
		UpdatedAt:       m.UpdatedAt,
	}

	return &res, false, nil
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

	m := &model.Room{ID: id.ID}
	q := s.db.NewUpdate().Model(m).WherePK().OmitZero().Returning("*")

	if req.RoomTypeID != "" {
		m.RoomTypeID = req.RoomTypeID
		q.Set("room_type_id = ?room_type_id")
	}
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
	if req.StartRoom > 0 {
		m.StartRoom = req.StartRoom
		q.Set("start_room = ?start_room")
	}
	if req.EndRoom > 0 {
		m.EndRoom = req.EndRoom
		q.Set("end_room = ?end_room")
	}
	if req.MaintenanceNote != "" {
		m.Maintenance_note = req.MaintenanceNote
		q.Set("maintenance_note = ?maintenance_note")
	}
	if req.MaintenanceETA != "" {
		m.Maintenance_eta = req.MaintenanceETA
		q.Set("maintenance_eta = ?maintenance_eta")
	}

	if req.Is_Available != nil { // ใช้ pointer เพื่อแยกกรณีไม่ได้ส่ง
		m.Is_Available = *req.Is_Available
		q.Set("is_available = ?is_available")
	}

	// update time
	m.SetUpdateNow()
	q.Set("updated_at = ?updated_at")

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
		ColumnExpr("r.room_type_id AS room_type_id").
		ColumnExpr("r.description AS description").
		ColumnExpr("r.capacity AS capacity").
		ColumnExpr("r.image_url AS image_url").
		ColumnExpr("r.start_room AS start_room").
		ColumnExpr("r.end_room AS end_room").
		ColumnExpr("r.maintenance_note AS maintenance_note").
		ColumnExpr("r.maintenance_eta AS maintenance_eta").
		ColumnExpr("r.is_available AS is_available").
		ColumnExpr("b.name AS building").
		ColumnExpr("r.created_at AS created_at").
		ColumnExpr("r.updated_at AS updated_at").
		Join("LEFT JOIN room_types as rt ON r.room_type_id::uuid = rt.id::uuid").
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
		ColumnExpr("r.room_type_id AS room_type_id").
		ColumnExpr("r.description AS description").
		ColumnExpr("r.capacity AS capacity").
		ColumnExpr("r.image_url AS image_url").
		ColumnExpr("r.start_room AS start_room").
		ColumnExpr("r.end_room AS end_room").
		ColumnExpr("r.maintenance_note AS maintenance_note").
		ColumnExpr("r.maintenance_eta AS maintenance_eta").
		ColumnExpr("r.is_available AS is_available").
		ColumnExpr("b.name AS building").
		ColumnExpr("r.created_at AS created_at").
		ColumnExpr("r.updated_at AS updated_at").
		Join("LEFT JOIN room_types as rt ON r.room_type_id::uuid = rt.id::uuid").
		Join("LEFT JOIN building_rooms as br ON r.id::uuid = br.room_id::uuid").
		Join("LEFT JOIN buildings as b ON br.building_id::uuid = b.id::uuid").
		Where("r.deleted_at IS NULL").
		Where("r.id = ?", id.ID).
		Scan(ctx, &m)
	return &m, err
}

func (s *Service) Delete(ctx context.Context, id request.GetByIdRoom) error {
	// เช็คว่ามี record นี้จริง (รวมที่ถูก soft delete แล้ว)
	exists, err := s.db.
		NewSelect().
		Model((*model.Room)(nil)).
		Where("id = ?", id.ID).
		WhereAllWithDeleted().
		Exists(ctx)
	if err != nil {
		return err
	}
	if !exists {
		return errors.New("room not found")
	}

	// ลบจริง (ข้าม soft delete)
	_, err = s.db.
		NewDelete().
		Model((*model.Room)(nil)).
		Where("id = ?", id.ID).
		ForceDelete().
		Exec(ctx)

	return err
}
