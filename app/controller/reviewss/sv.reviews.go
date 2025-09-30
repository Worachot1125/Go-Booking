package reviewss

import (
	"app/app/model"
	"app/app/request"
	"app/app/response"
	"context"
	"errors"
	"fmt"
	"strings"
)

// Create รีวิวใหม่
func (s *Service) Create(ctx context.Context, req request.CreateReviews) (*response.ReviewsResponse, bool, error) {
	// ตรวจสอบว่าผู้ใช้ได้รีวิวห้องนี้แล้วหรือไม่
	ex, err := s.db.NewSelect().
		Model((*model.Reviews)(nil)).
		Where("room_id = ?", req.RoomID).
		Where("user_id = ?", req.UserID).
		Exists(ctx)
	if err != nil {
		return nil, false, err
	}
	if ex {
		return nil, true, errors.New("คุณได้รีวิวห้องนี้แล้ว")
	}

	m := &model.Reviews{
		UserID:  req.UserID,
		RoomID:  req.RoomID,
		Rating:  req.Rating,
		Comment: req.Comment,
	}

	_, err = s.db.NewInsert().Model(m).Returning("*").Exec(ctx)
	if err != nil {
		return nil, false, err
	}

	// Map ไป struct response
	resp := &response.ReviewsResponse{
		ID:        m.ID,
		User_ID:   m.UserID,
		Room_ID:   m.RoomID,
		Rating:    int64(m.Rating),
		Comment:   m.Comment,
		CreatedAt: m.CreatedAt,
		UpdatedAt: m.UpdatedAt,
	}

	return resp, false, nil
}

// Update รีวิว
func (s *Service) Update(ctx context.Context, req request.UpdateReviews, id request.GetByIDReviews) (*response.ReviewsResponse, bool, error) {
	m := &model.Reviews{
		ID:      id.ID,
		Rating:  req.Rating,
		Comment: req.Comment,
	}
	m.SetUpdateNow()

	_, err := s.db.NewUpdate().
		Model(m).
		WherePK().
		OmitZero().
		Returning("*").
		Exec(ctx)
	if err != nil {
		return nil, false, err
	}

	resp := &response.ReviewsResponse{
		ID:        m.ID,
		User_ID:   m.UserID,
		Room_ID:   m.RoomID,
		Rating:    int64(m.Rating),
		Comment:   m.Comment,
		CreatedAt: m.CreatedAt,
		UpdatedAt: m.UpdatedAt,
	}

	return resp, false, nil
}

// List รีวิว
func (s *Service) List(ctx context.Context, req request.ListReviews) ([]response.ReviewsResponse, int, error) {
	offset := (req.Page - 1) * req.Size
	m := []response.ReviewsResponse{}

	query := s.db.NewSelect().
		TableExpr("reviews AS rv").
		ColumnExpr("rv.id").
		ColumnExpr("rv.user_id").
		ColumnExpr("rv.room_id").
		ColumnExpr("rv.rating").
		ColumnExpr("rv.comment").
		ColumnExpr("rv.created_at").
		ColumnExpr("rv.updated_at").
		Where("rv.deleted_at IS NULL")

	// Filtering
	if req.Search != "" {
		search := "%" + strings.ToLower(req.Search) + "%"
		if req.SearchBy != "" {
			searchBy := strings.ToLower(req.SearchBy)
			query = query.Where(fmt.Sprintf("LOWER(rv.%s) LIKE ?", searchBy), search)
		} else {
			query = query.Where("LOWER(rv.comment) LIKE ?", search)
		}
	}

	// Count
	count, err := query.Count(ctx)
	if err != nil {
		return nil, 0, err
	}

	// Order
	order := fmt.Sprintf("rv.%s %s", req.SortBy, req.OrderBy)

	// Scan เข้า response struct
	err = query.Order(order).Limit(req.Size).Offset(offset).Scan(ctx, &m)
	if err != nil {
		return nil, 0, err
	}

	return m, count, nil
}

// Get รีวิว
func (s *Service) Get(ctx context.Context, id request.GetByIDReviews) (*response.ReviewsResponse, error) {
	m := response.ReviewsResponse{}

	err := s.db.NewSelect().
		TableExpr("reviews AS rv").
		ColumnExpr("rv.id").
		ColumnExpr("rv.user_id").
		ColumnExpr("rv.room_id").
		ColumnExpr("rv.rating").
		ColumnExpr("rv.comment").
		ColumnExpr("rv.created_at").
		ColumnExpr("rv.updated_at").
		Where("rv.id = ?", id.ID).
		Where("rv.deleted_at IS NULL").
		Scan(ctx, &m)

	if err != nil {
		return nil, err
	}

	return &m, nil
}

// Delete รีวิว
func (s *Service) Delete(ctx context.Context, id request.GetByIDReviews) error {
	ex, err := s.db.NewSelect().
		Table("reviews").
		Where("id = ?", id.ID).
		Where("deleted_at IS NULL").
		Exists(ctx)
	if err != nil {
		return err
	}
	if !ex {
		return errors.New("review not found")
	}

	_, err = s.db.NewDelete().
		Model((*model.Reviews)(nil)).
		Where("id = ?", id.ID).
		Exec(ctx)
	return err
}
