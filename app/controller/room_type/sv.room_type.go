package room_type

import (
	"app/app/model"
	"app/app/request"
	"app/app/response"
	"context"
	"errors"
	"fmt"
	"strings"
)

func (s *Service) List(ctx context.Context, req request.ListRoomType) ([]response.RoomTypeResponse, int, error) {
	offset := (req.Page - 1) * req.Size
	m := []response.RoomTypeResponse{}

	query := s.db.NewSelect().
		TableExpr("room_types as rt").
		Column("rt.id", "rt.name").
		Column("rt.created_at").
		Column("rt.updated_at").
		Where("deleted_at IS NULL")

	// Filtering
	if req.Search != "" {
		search := "%" + strings.ToLower(req.Search) + "%"
		if req.SearchBy != "" {
			searchBy := strings.ToLower(req.SearchBy)
			query = query.Where(fmt.Sprintf("LOWER(rt.%s) LIKE ?", searchBy), search)
		} else {
			query = query.Where("LOWER(rt.name) LIKE ?", search)
		}
	}

	// Count total before pagination
	count, err := query.Count(ctx)
	if err != nil {
		return nil, 0, err
	}

	// Order handling
	order := fmt.Sprintf("rt.%s %s", req.SortBy, req.OrderBy)

	// Final query with order + pagination
	err = query.Order(order).Limit(req.Size).Offset(offset).Scan(ctx, &m)
	if err != nil {
		return nil, 0, err
	}

	return m, count, nil
}

func (s *Service) Delete(ctx context.Context, id string) (*model.RoomType, bool, error) {
	ex, err := s.db.NewSelect().Model((*model.RoomType)(nil)).Where("id = ?", id).Exists(ctx)
	if err != nil {
		return nil, false, err
	}

	if !ex {
		return nil, true, errors.New("role not found")
	}

	_, err = s.db.NewDelete().Model((*model.RoomType)(nil)).Where("id = ?", id).Exec(ctx)

	return nil, false, err
}
