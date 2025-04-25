package role_permission

import (
	"context"
	"errors"
	"strings"

	"app/app/model"
	"app/app/request"
)

func (s *Service) Create(ctx context.Context, req request.Role_PermissionCreate) (*model.Role_Permission, bool, error) {
	rp := model.Role_Permission{
		Role_ID:       req.Role_ID,
		Permission_ID: req.Permission_ID,
	}
	rp.SetCreatedNow()
	_, err := s.db.NewInsert().Model(&rp).Exec(ctx)
	if err != nil {
		if strings.Contains(err.Error(), "duplicate key value") {
			return nil, true, errors.New("role_permission already exists")
		}
	}
	return &rp, false, err
}

func (s *Service) Update(ctx context.Context, id string, req request.Role_PermissionUpdate) (*model.Role_Permission, bool, error) {
	ex, err := s.db.NewSelect().Model((*model.Role_Permission)(nil)).Where("id = ?", id).Exists(ctx)
	if err != nil {
		return nil, false, err
	}

	if !ex {
		return nil, true, errors.New("role_permission not found")
	}

	m := model.Role_Permission{
		Role_ID:       req.Role_ID,
		Permission_ID: req.Permission_ID,
	}

	m.SetUpdateNow()

	_, err = s.db.NewUpdate().Model(&m).
		Set("role_id = ?role_id").
		Set("permission_id = ?permission_id").
		WherePK().
		OmitZero().
		Returning("*").
		Exec(ctx)

	return &m, false, err
}
