package login

import (
	"app/app/model"
	"context"

	"github.com/uptrace/bun"
)

type Repository struct {
	db *bun.DB
}

func NewRepository(db *bun.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) FindByEmail(ctx context.Context, email string) (*model.User, error) {
	user := new(model.User)
	err := r.db.NewSelect().Model(user).
		Where("email = ?", email).
		Limit(1).
		Scan(ctx)
	if err != nil {
		return nil, err
	}
	return user, nil
}
