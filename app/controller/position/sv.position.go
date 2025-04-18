package position

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"app/app/model"
	"app/app/request"
)

func (s *Service) Create(ctx context.Context, req request.PositionCreate) (*model.Position, bool, error) {
	m := model.Position{
		Name: req.Name,
	}
	_, err := s.db.NewInsert().Model(&m).Exec(ctx)
	if err != nil {
		if strings.Contains(err.Error(), "duplicate key value") {
			return nil, true, errors.New("position already exists")
		}
	}
	return &m, false, err
}

func (s *Service) Update(ctx context.Context, id string, req request.PositionUpdate) (*model.Position, bool, error) {
	ex, err := s.db.NewSelect().Model((*model.Position)(nil)).Where("id = ?", id).Exists(ctx)
	if err != nil {
		return nil, false, err
	}

	if !ex {
		return nil, true, errors.New("product not found")
	}

	m := model.Position{
		ID:   id,
		Name: req.Name,
	}

	m.SetUpdateNow()

	_, err = s.db.NewUpdate().Model(&m).
		Set("name = ?name").
		Set("updated_at = ?updated_at").
		WherePK().
		OmitZero().
		Returning("*").
		Exec(ctx)

	return &m, false, err
}

func (s *Service) Delete(ctx context.Context, id string) (*model.Position, bool, error) {
	ex, err := s.db.NewSelect().Model((*model.Position)(nil)).Where("id = ?", id).Exists(ctx)
	if err != nil {
		return nil, false, err
	}

	if !ex {
		return nil, true, errors.New("position not found")
	}

	_, err = s.db.NewDelete().Model((*model.Position)(nil)).Where("id = ?", id).Exec(ctx)

	return nil, false, err
}

func (s *Service) Get(ctx context.Context, id string) (*model.Position, error) {
	m := model.Position{}

	err := s.db.NewSelect().Model(&m).Where("id = ?", id).Scan(ctx)
	return &m, err
}

func (s *Service) List(ctx context.Context, req request.PositionListRequest) ([]model.Position, int, error) {
	m := []model.Position{}

	var (
		offset = (req.Page - 1) * req.Size
		limit  = req.Size
	)

	query := s.db.NewSelect().Model(&m)

	if req.Search != "" {
		search := fmt.Sprint("%" + strings.ToLower(req.Search) + "%")
		query.Where("LOWER(name) Like ?", search)
	}

	count, err := query.Count(ctx)
	if count == 0 {
		return m, 0, err
	}

	order := fmt.Sprintf("%s %s", req.SortBy, req.OrderBy)
	err = query.Offset(offset).Limit(limit).Order(order).Scan(ctx, &m)

	return m, count, err
}
