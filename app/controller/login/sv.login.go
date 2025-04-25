package login

import (
	"app/app/model"
	"context"
	"errors"
)

func (s *Service) Login(ctx context.Context, email, password string) (*model.User, error) {
	user := new(model.User) // สร้าง user instance ก่อน

	err := s.db.NewSelect().Model(user).
		Where("email = ?", email).
		Limit(1).
		Scan(ctx)
	if err != nil {
		return nil, errors.New("user not found")
	}

	if !user.CheckPassword(password) {
		return nil, errors.New("invalid password")
	}

	return user, nil
}

