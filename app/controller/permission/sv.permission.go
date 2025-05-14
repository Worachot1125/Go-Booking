package permission

import (
	"app/app/model"
	"app/app/request"
	"context"
	"errors"
	"strings"
)

func (s *Service) Create(ctx context.Context, req request.PermissionCreate) (*model.Permission, bool, error) {
	rol := model.Permission{
		Name:        req.Name,
		Description: req.Description,
	}
	rol.SetCreatedNow()
	_, err := s.db.NewInsert().Model(&rol).Exec(ctx)
	if err != nil {
		if strings.Contains(err.Error(), "duplicate key value") {
			return nil, true, errors.New("permission already exists")
		}
	}
	return &rol, false, err
}

func (s *Service) Update(ctx context.Context, req request.PermissionUpdate, id request.PermissionGetByID) (*model.Permission, bool, error) {
	ex, err := s.db.NewSelect().Table("permissions").Where("id = ?", id.ID).Exists(ctx)
	if err != nil {
		return nil, false, err
	}
	if !ex {
		return nil, false, err
	}

	m := &model.Permission{
		ID:          id.ID,
		Name:        req.Name,
		Description: req.Description,
	}
	m.SetUpdateNow()
	_, err = s.db.NewUpdate().Model(m).
		Set("name = ?name").
		Set("description = ?description").
		WherePK().
		OmitZero().
		Returning("*").
		Exec(ctx)
	if err != nil {
		if strings.Contains(err.Error(), "duplicate key value") {
			return nil, true, errors.New("permission already exists")
		}
	}
	return m, false, err
}

func (s *Service) Delete(ctx context.Context, id string) (*model.Permission, bool, error) {
	ex, err := s.db.NewSelect().Model((*model.Permission)(nil)).Where("id = ?", id).Exists(ctx)
	if err != nil {
		return nil, false, err
	}

	if !ex {
		return nil, true, errors.New("permission not found")
	}

	_, err = s.db.NewDelete().Model((*model.Permission)(nil)).Where("id = ?", id).Exec(ctx)

	return nil, false, err
}
