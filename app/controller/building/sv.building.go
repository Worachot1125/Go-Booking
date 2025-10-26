package building

import (
	"app/app/model"
	"app/app/request"
	"app/app/response"
	"context"
	"errors"
	"fmt"
	"strings"
)

func (s *Service) Create(ctx context.Context, req request.CreateBuilding) (*model.Building, bool, error) {
	m := model.Building{
		Name: req.Name,
		Image_url:   req.Image_url,
	}
	_, err := s.db.NewInsert().Model(&m).Exec(ctx)
	if err != nil {
		if strings.Contains(err.Error(), "duplicate key value") {
			return nil, true, errors.New("building already exists")
		}
	}
	return &m, false, err
}

func (s *Service) Update(ctx context.Context, req request.UpdateBuilding, id request.GetByIdBuilding) (*model.Building, bool, error) {
	ex, err := s.db.NewSelect().Table("buildings").Where("id = ?", id.ID).Exists(ctx)
	if err != nil {
		return nil, false, err
	}
	if !ex {
		return nil, false, err
	}

	m := &model.Building{
    ID:   id.ID,
    Name: req.Name,
    Image_url: req.Image_url, // เพิ่มบรรทัดนี้!
}

m.SetUpdateNow()

_, err = s.db.NewUpdate().Model(m).
    Set("name = ?name").
    Set("image_url = ?image_url"). // เพิ่มบรรทัดนี้!
    Set("updated_at = ?updated_at").
    WherePK().
    OmitZero().
    Returning("*").
    Exec(ctx)
	if err != nil {
		if strings.Contains(err.Error(), "duplicate key value") {
			return nil, true, errors.New("building already exists")
		}
	}

	return m, false, err
}

func (s *Service) List(ctx context.Context, req request.ListBuilding) ([]response.BuildingResponse, int, error) {
	offset := (req.Page - 1) * req.Size
	m := []response.BuildingResponse{}

	query := s.db.NewSelect().
		TableExpr("buildings AS b").
		ColumnExpr("b.id").
		ColumnExpr("b.name").
		ColumnExpr("b.image_url").
		ColumnExpr("b.created_at").
		ColumnExpr("b.updated_at").
		ColumnExpr(`COALESCE(
					json_agg(
						jsonb_build_object(
							'id', r.id,
							'name', r.name,
							'description', r.description,
							'capacity', r.capacity,
							'image_url', r.image_url,
							'created_at', r.created_at,
							'updated_at', r.updated_at
						) ORDER BY r.name
					) FILTER (WHERE r.id IS NOT NULL),
					'[]'
				) AS rooms`). // แก้ชื่อเป็น rooms
		Join("LEFT JOIN building_rooms AS br ON br.building_id::uuid = b.id::uuid").
		Join("LEFT JOIN rooms AS r ON r.id::uuid = br.room_id::uuid").
		Where("b.deleted_at IS NULL").
		Where("r.deleted_at IS NULL").
		GroupExpr("b.id, b.name, b.created_at, b.updated_at").
		OrderExpr("b.name ASC")

	// Filtering
	if req.Search != "" {
		search := "%" + strings.ToLower(req.Search) + "%"
		if req.SearchBy != "" {
			searchBy := strings.ToLower(req.SearchBy)
			query = query.Where(fmt.Sprintf("LOWER(b.%s) LIKE ?", searchBy), search)
		} else {
			query = query.Where("LOWER(b.name) LIKE ?", search)
		}
	}

	// Count total before pagination
	countQuery := s.db.NewSelect().
		TableExpr("buildings AS b").
		ColumnExpr("COUNT(*)").
		Where("b.deleted_at IS NULL")

	if req.Search != "" {
		search := "%" + strings.ToLower(req.Search) + "%"
		if req.SearchBy != "" {
			searchBy := strings.ToLower(req.SearchBy)
			countQuery = countQuery.Where(fmt.Sprintf("LOWER(b.%s) LIKE ?", searchBy), search)
		} else {
			countQuery = countQuery.Where("LOWER(b.name) LIKE ?", search)
		}
	}

	var count int
	if err := countQuery.Scan(ctx, &count); err != nil {
		return nil, 0, err
	}

	// Sorting and Pagination
	order := fmt.Sprintf("b.%s %s", req.SortBy, req.OrderBy)
	err := query.Order(order).Limit(req.Size).Offset(offset).Scan(ctx, &m)
	if err != nil {
		return nil, 0, err
	}

	return m, count, nil
}


func (s *Service) Get(ctx context.Context, id request.GetByIdBuilding) ([]response.BuildingResponse, error) {
	m := []response.BuildingResponse{}

	err := s.db.NewSelect().
		TableExpr("buildings AS b").
		ColumnExpr("b.id").
		ColumnExpr("b.name").
		ColumnExpr("b.created_at").
		ColumnExpr("b.updated_at").
		ColumnExpr(`COALESCE(
						json_agg(
							jsonb_build_object(
								'id', r.id,
								'name', r.name
							) ORDER BY r.name
						) FILTER (WHERE r.id IS NOT NULL),
						'[]'
					) AS rooms_name`).
		Join("LEFT JOIN building_rooms AS br ON br.building_id::uuid = b.id::uuid").
		Join("LEFT JOIN rooms AS r ON r.id::uuid = br.room_id::uuid").
		Where("b.deleted_at IS NULL").
		GroupExpr("b.id").
		Where("r.deleted_at IS NULL").
		OrderExpr("b.name ASC").
		Where("b.id = ?", id.ID).
		Scan(ctx, &m)
	return m, err
}

func (s *Service) Delete(ctx context.Context, id request.GetByIdBuilding) error {
	// เช็คว่ามี record นี้จริง (นับรวมที่เคย soft-delete ไปแล้วด้วย)
	exists, err := s.db.
		NewSelect().
		Model((*model.Building)(nil)).
		Where("id = ?", id.ID).
		WhereAllWithDeleted().
		Exists(ctx)
	if err != nil {
		return err
	}
	if !exists {
		return errors.New("building not found")
	}

	// ลบจริง (hard delete) — ข้าม soft delete
	_, err = s.db.
		NewDelete().
		Model((*model.Building)(nil)).
		Where("id = ?", id.ID).
		ForceDelete().
		Exec(ctx)
	return err
}