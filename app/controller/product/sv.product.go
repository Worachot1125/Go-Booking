package product

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"app/app/model"
	"app/app/request"
)

func (s *Service) Create(ctx context.Context, req request.ProductCeate) (*model.Product, bool, error) {
	m := model.Product{
		Name:        req.Name,
		Price:       req.Price,
		Description: req.Description,
	}
	_, err := s.db.NewInsert().Model(&m).Exec(ctx)
	if err != nil {
		if strings.Contains(err.Error(), "duplicate key value") {
			return nil, true, errors.New("product already exists")
		}
	}
	return &m, false, err
}

func (s *Service) Update(ctx context.Context, id int64, req request.ProductUpdate) (*model.Product, bool, error) {
	ex, err := s.db.NewSelect().Model((*model.Product)(nil)).Where("id = ?", id).Exists(ctx)
	if err != nil {
		return nil, false, err
	}

	if !ex {
		return nil, true, errors.New("product not found")
	}

	m := model.Product{
		ID:          id,
		Name:        req.Name,
		Price:       req.Price,
		Description: req.Description,
	}

	m.SetUpdateNow()

	_, err = s.db.NewUpdate().Model(&m).
		Set("name = ?name").
		Set("price = ?price").
		Set("description = ?description").
		Set("updated_at = ?updated_at").
		WherePK().
		OmitZero().
		Returning("*").
		Exec(ctx)

	return &m, false, err
}

func (s *Service) Delete(ctx context.Context, id int64) (*model.Product, bool, error) {
	ex, err := s.db.NewSelect().Model((*model.Product)(nil)).Where("id = ?", id).Exists(ctx)
	if err != nil {
		return nil, false, err
	}

	if !ex {
		return nil, true, errors.New("product not found")
	}

	_, err = s.db.NewDelete().Model((*model.Product)(nil)).Where("id = ?", id).Exec(ctx)

	return nil, false, err
}

func (s *Service) Get(ctx context.Context, id int64) (*model.Product, error) {
	m := model.Product{}

	err := s.db.NewSelect().Model(&m).Where("id = ?", id).Scan(ctx)
	return &m, err
}

func (s *Service) List(ctx context.Context, req request.ProductListReuest) ([]model.Product, int, error) {
	m := []model.Product{}

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
