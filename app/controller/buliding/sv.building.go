package building

import (
	"app/app/model"
	"app/app/request"
	_"app/app/response"
	"context"
	"errors"
	_"fmt"
	"strings"
)

func (s *Service) Create(ctx context.Context, req request.CreateBuilding) (*model.Building, bool, error) {
	m := model.Building{
		Name:        req.Name,
	}
	_, err := s.db.NewInsert().Model(&m).Exec(ctx)
	if err != nil {
		if strings.Contains(err.Error(), "duplicate key value") {
			return nil, true, errors.New("building already exists")
		}
	}
	return &m, false, err
}

// func (s *Service) Update(ctx context.Context, req request.UpdateRoom, id request.GetByIdRoom) (*model.Room, bool, error) {
// 	ex, err := s.db.NewSelect().Table("rooms").Where("id = ?", id.ID).Exists(ctx)
// 	if err != nil {
// 		return nil, false, err
// 	}
// 	if !ex {
// 		return nil, false, err
// 	}

// 	m := &model.Room{
// 		ID:          id.ID,
// 		Name:        req.Name,
// 		Description: req.Description,
// 		Capacity:    req.Capacity,
// 		Image_url:   req.Image_url,
// 	}

// 	m.SetUpdateNow()

// 	_, err = s.db.NewUpdate().Model(m).
// 		Set("name = ?name").
// 		Set("description = ?description").
// 		Set("capacity = ?capacity").
// 		Set("image_url = ?image_url").
// 		Set("updated_at = ?updated_at").
// 		WherePK().
// 		OmitZero().
// 		Returning("*").
// 		Exec(ctx)
// 	if err != nil {
// 		if strings.Contains(err.Error(), "duplicate key value") {
// 			return nil, true, errors.New("room already exists")
// 		}
// 	}

// 	return m, false, err
// }

// func (s *Service) List(ctx context.Context, req request.ListRoom) ([]response.RooomResponse, int, error) {
// 	offset := (req.Page - 1) * req.Size
// 	m := []response.RooomResponse{}

// 	query := s.db.NewSelect().
// 		TableExpr("rooms as r").
// 		Column("r.id", "r.name", "r.description", "r.capacity", "r.updated_at").Where("deleted_at IS NULL")

// 	// Filtering
// 	if req.Search != "" {
// 		search := "%" + strings.ToLower(req.Search) + "%"
// 		if req.SearchBy != "" {
// 			searchBy := strings.ToLower(req.SearchBy)
// 			query = query.Where(fmt.Sprintf("LOWER(r.%s) LIKE ?", searchBy), search)
// 		} else {
// 			query = query.Where("LOWER(r.name) LIKE ?", search)
// 		}
// 	}

// 	// Count total before pagination
// 	count, err := query.Count(ctx)
// 	if err != nil {
// 		return nil, 0, err
// 	}

// 	// Order handling
// 	order := fmt.Sprintf("r.%s %s", req.SortBy, req.OrderBy)

// 	// Final query with order + pagination
// 	err = query.Order(order).Limit(req.Size).Offset(offset).Scan(ctx, &m)
// 	if err != nil {
// 		return nil, 0, err
// 	}

// 	return m, count, nil
// }

// func (s *Service) Get(ctx context.Context, id request.GetByIdRoom) (*response.RooomResponse, error) {
// 	m := response.RooomResponse{}

// 	err := s.db.NewSelect().
// 		TableExpr("rooms as r").
// 		Column("r.id", "r.name", "r.description", "r.capacity", "r.updated_at").Where("deleted_at IS NULL").
// 		Where("id = ?", id.ID).Where("deleted_at IS NULL").Scan(ctx, &m)
// 	return &m, err
// }

// func (s *Service) Delete(ctx context.Context, id request.GetByIdRoom) error {
// 	ex, err := s.db.NewSelect().Table("rooms").Where("id = ?", id.ID).Where("deleted_at IS NULL").Exists(ctx)
// 	if err != nil {
// 		return  err
// 	}
// 	if !ex {
// 		return errors.New("room not found")
// 	}

// 	_, err = s.db.NewDelete().Model((*model.Room)(nil)).Where("id = ?",id.ID).Exec(ctx)
// 	return err
// }
