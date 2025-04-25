package user_role

import (
	"context"
	"errors"
	"strings"

	"app/app/model"
	"app/app/request"
)

func (s *Service) Create(ctx context.Context, req request.User_RoleCreate) (*model.User_Role, bool, error) {
	rp := model.User_Role{
		Role_ID: req.Role_ID,
		User_ID: req.User_ID,
	}
	rp.SetCreatedNow()
	_, err := s.db.NewInsert().Model(&rp).Exec(ctx)
	if err != nil {
		if strings.Contains(err.Error(), "duplicate key value") {
			return nil, true, errors.New("user_role already exists")
		}
	}
	return &rp, false, err
}