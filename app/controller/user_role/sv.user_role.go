package user_role

import (
	"context"
	"errors"
	"strings"

	"app/app/model"
	"app/app/request"
	"app/app/response"
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

func (s *Service) GetUserRolesByUserID(ctx context.Context, userID string) ([]response.UserRoleByUserID, error) {
	var result []response.UserRoleByUserID

	err := s.db.NewSelect().
		TableExpr("user_roles AS ur").
		ColumnExpr("ur.user_id AS user_id").
		ColumnExpr("ur.role_id AS role_id").
		ColumnExpr("u.first_name AS first_name").
		ColumnExpr("u.last_name AS last_name").
		ColumnExpr("r.name AS role_name").
		Join("JOIN users AS u ON u.id::uuid = ur.user_id::uuid"). 
		Join("JOIN roles AS r ON r.id::uuid = ur.role_id::uuid").
		Where("ur.user_id::uuid = ?", userID).
		Where("ur.deleted_at IS NULL").
		Where("u.deleted_at IS NULL").
		Where("r.deleted_at IS NULL").
		Scan(ctx, &result)

	if err != nil {
		return nil, err
	}

	return result, nil
}

