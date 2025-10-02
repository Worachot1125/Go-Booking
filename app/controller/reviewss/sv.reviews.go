package reviewss

import (
	"app/app/model"
	"app/app/request"
	"app/app/response"
	"context"
	"errors"
	"strings"
)

func (s *Service) Create(ctx context.Context, req request.CreateReviews) (*response.ReviewsResponse, bool, error) {
	ex, err := s.db.NewSelect().
		Model((*model.Reviews)(nil)).
		Where("room_id = ?", req.RoomID).
		Where("user_id = ?", req.UserID).
		Where("booking_id = ?", req.BookingID).
		Exists(ctx)
	if err != nil {
		return nil, false, err
	}
	if ex {
		return nil, true, errors.New("คุณได้รีวิวห้องนี้แล้ว")
	}

	m := &model.Reviews{
		UserID:    req.UserID,
		RoomID:    req.RoomID,
		BookingID: req.BookingID,
		Rating:    req.Rating,
		Comment:   req.Comment,
	}

	_, err = s.db.NewInsert().Model(m).Returning("*").Exec(ctx)
	if err != nil {
		return nil, false, err
	}

	resp := &response.ReviewsResponse{
		ID:        m.ID,
		User_ID:   m.UserID,
		Room_ID:   m.RoomID,
		Booking_ID: m.BookingID,
		Rating:    int64(m.Rating),
		Comment:   m.Comment,
		CreatedAt: m.CreatedAt,
		UpdatedAt: m.UpdatedAt,
	}

	return resp, false, nil
}

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
		Booking_ID: m.BookingID,
		Rating:    int64(m.Rating),
		Comment:   m.Comment,
		CreatedAt: m.CreatedAt,
		UpdatedAt: m.UpdatedAt,
	}

	return resp, false, nil
}

func (s *Service) List(ctx context.Context, req request.ListReviews) ([]response.ReviewsResponse, int, error) {
	offset := (req.Page - 1) * req.Size
	m := []response.ReviewsResponse{}

	query := s.db.NewSelect().
		TableExpr("reviews AS rv").
		ColumnExpr("rv.id").
		ColumnExpr("rv.user_id").
		ColumnExpr("rv.room_id").
		ColumnExpr("rv.booking_id").
		ColumnExpr("rv.rating").
		ColumnExpr("rv.comment").
		ColumnExpr("rv.created_at").
		ColumnExpr("rv.updated_at").
		Where("rv.deleted_at IS NULL")

	if req.Search != "" {
		search := "%" + strings.ToLower(req.Search) + "%"
		if req.SearchBy != "" {
			query = query.Where("LOWER(rv."+req.SearchBy+") LIKE ?", search)
		} else {
			query = query.Where("LOWER(rv.comment) LIKE ?", search)
		}
	}

	count, err := query.Count(ctx)
	if err != nil {
		return nil, 0, err
	}

	order := "rv." + req.SortBy + " " + req.OrderBy
	err = query.Order(order).Limit(req.Size).Offset(offset).Scan(ctx, &m)
	if err != nil {
		return nil, 0, err
	}

	return m, count, nil
}

func (s *Service) Get(ctx context.Context, id request.GetByIDReviews) (*response.ReviewsResponse, error) {
	m := response.ReviewsResponse{}
	err := s.db.NewSelect().
		TableExpr("reviews AS rv").
		ColumnExpr("rv.id").
		ColumnExpr("rv.user_id").
		ColumnExpr("rv.room_id").
		ColumnExpr("rv.booking_id").
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

func (s *Service) GetByBookingID(ctx context.Context, bookingID request.GetByBookingIDReviews) ([]response.ReviewsResponse, error) {
	var reviews []response.ReviewsResponse
	err := s.db.NewSelect().
		TableExpr("reviews AS rv").
		ColumnExpr("rv.id").
		ColumnExpr("rv.user_id").
		ColumnExpr("rv.room_id").
		ColumnExpr("rv.booking_id").
		ColumnExpr("rv.rating").
		ColumnExpr("rv.comment").
		ColumnExpr("rv.created_at").
		ColumnExpr("rv.updated_at").
		Where("rv.booking_id = ?", bookingID.BookingID).
		Where("rv.deleted_at IS NULL").
		Scan(ctx, &reviews)
	if err != nil {
		return nil, err
	}
	return reviews, nil
}

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
