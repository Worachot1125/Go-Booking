package reportss

import (
	"app/app/model"
	"app/app/request"
	"app/app/response"
	"context"
	"errors"
	"fmt"
	"strings"
)

func (s *Service) Create(ctx context.Context, req request.CreateReport) (*response.ReportResponse, error) {
	m := &model.Report{
		UserID:      req.UserID,
		RoomID:      req.RoomID,
		Description: req.Description,
	}

	_, err := s.db.NewInsert().Model(m).Returning("*").Exec(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to create report: %w", err)
	}

	resp := &response.ReportResponse{
		ID:          m.ID,
		User_ID:     m.UserID,
		Room_ID:     m.RoomID,
		Description: m.Description,
		CreatedAt:   m.CreatedAt,
		UpdatedAt:   m.UpdatedAt,
	}

	return resp, nil
}

func (s *Service) Update(ctx context.Context, req request.UpdateReport, id request.GetByIDReport) (*response.ReportResponse, error) {
	m := &model.Report{
		ID:          id.ID,
		UserID:      req.UserID,
		RoomID:      req.RoomID,
		Description: req.Description,
	}
	m.SetUpdateNow()

	res, err := s.db.NewUpdate().
		Model(m).
		WherePK().
		OmitZero().
		Returning("*").
		Exec(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to update report: %w", err)
	}

	affected, _ := res.RowsAffected()
	if affected == 0 {
		return nil, errors.New("report not found")
	}

	resp := &response.ReportResponse{
		ID:          m.ID,
		User_ID:     m.UserID,
		Room_ID:     m.RoomID,
		Description: m.Description,
		CreatedAt:   m.CreatedAt,
		UpdatedAt:   m.UpdatedAt,
	}

	return resp, nil
}

func (s *Service) List(ctx context.Context, req request.ListReport) ([]response.ReportResponse, int, error) {
	offset := (req.Page - 1) * req.Size
	m := []response.ReportResponse{}

	query := s.db.NewSelect().
		TableExpr("reports AS rep").
		ColumnExpr("rep.id").
		ColumnExpr("rep.user_id").
		ColumnExpr("rep.room_id").
		ColumnExpr("rep.description").
		ColumnExpr("rep.created_at").
		ColumnExpr("rep.updated_at").
		Where("rep.deleted_at IS NULL")

	if req.Search != "" {
		search := "%" + strings.ToLower(req.Search) + "%"
		if req.SearchBy != "" {
			searchBy := strings.ToLower(req.SearchBy)
			query = query.Where(fmt.Sprintf("LOWER(rep.%s) LIKE ?", searchBy), search)
		} else {
			query = query.Where("LOWER(rep.description) LIKE ?", search)
		}
	}

	// นับจำนวนก่อน
	count, err := query.Count(ctx)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count reports: %w", err)
	}
	if count == 0 {
		return []response.ReportResponse{}, 0, nil
	}

	order := fmt.Sprintf("rep.%s %s", req.SortBy, req.OrderBy)

	err = query.Order(order).Limit(req.Size).Offset(offset).Scan(ctx, &m)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to list reports: %w", err)
	}

	return m, count, nil
}

func (s *Service) Get(ctx context.Context, id request.GetByIDReport) (*response.ReportResponse, error) {
	m := response.ReportResponse{}

	err := s.db.NewSelect().
		TableExpr("reports AS rep").
		ColumnExpr("rep.id").
		ColumnExpr("rep.user_id").
		ColumnExpr("rep.room_id").
		ColumnExpr("rep.description").
		ColumnExpr("rep.created_at").
		ColumnExpr("rep.updated_at").
		Where("rep.id = ?", id.ID).
		Where("rep.deleted_at IS NULL").
		Scan(ctx, &m)

	if err != nil {
		return nil, fmt.Errorf("failed to get report: %w", err)
	}

	// กัน panic: ถ้าไม่เจอ record
	if m.ID == "" {
		return nil, errors.New("report not found")
	}

	return &m, nil
}

func (s *Service) Delete(ctx context.Context, id request.GetByIDReport) error {
	ex, err := s.db.NewSelect().
		Table("reports").
		Where("id = ?", id.ID).
		Where("deleted_at IS NULL").
		Exists(ctx)
	if err != nil {
		return fmt.Errorf("failed to check report existence: %w", err)
	}
	if !ex {
		return errors.New("report not found")
	}

	_, err = s.db.NewDelete().
		Model((*model.Report)(nil)).
		Where("id = ?", id.ID).
		Exec(ctx)
	if err != nil {
		return fmt.Errorf("failed to delete report: %w", err)
	}
	return nil
}