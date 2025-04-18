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
	_, err := s.db.NewInsert().Model(&rp).Exec(ctx)
	if err != nil {
		if strings.Contains(err.Error(), "duplicate key value") {
			return nil, true, errors.New("role_permission already exists")
		}
	}
	return &rp, false, err
}
