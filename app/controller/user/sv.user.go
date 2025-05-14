package user

import (
	"app/app/model"
	"app/app/request"
	"app/app/response"
	"context"
	"errors"
	"fmt"
	"strings"

	"golang.org/x/crypto/bcrypt"
)

func (s *Service) Create(ctx context.Context, req request.CreateUser) (*model.User, bool, error) {
	position := &model.Position{}
	err := s.db.NewSelect().Model(position).Where("name = ?", req.Position_Name).Scan(ctx)
	if err != nil {
		return nil, true, fmt.Errorf("position '%s' not found", req.Position_Name)
	}

	bytes, err := bcrypt.GenerateFromPassword([]byte(req.Password), 10)
	if err != nil {
		return nil, false, err
	}

	user := &model.User{
		FirstName:   req.FirstName,
		LastName:    req.LastName,
		Email:       req.Email,
		Password:    string(bytes),
		Position_ID: position.ID,
		Image_url:   req.Image_url,
		Phone:       req.Phone,
	}
	user.SetCreatedNow()

	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, false, err
	}

	defer func() {
		if err != nil {
			_ = tx.Rollback()
		}
	}()

	_, err = tx.NewInsert().Model(user).Exec(ctx)
	if err != nil {
		if strings.Contains(err.Error(), "duplicate key value") {
			return nil, true, errors.New("user already exists")
		}
		return nil, false, err
	}

	role := &model.Role{}
	err = tx.NewSelect().Model(role).Where("name = ?", "Employee").Scan(ctx)
	if err != nil {
		return nil, false, errors.New("role 'Employee' not found")
	}

	userRole := &model.User_Role{
		User_ID: user.ID,
		Role_ID: role.ID,
	}
	_, err = tx.NewInsert().Model(userRole).Exec(ctx)
	if err != nil {
		return nil, false, err
	}

	err = tx.Commit()
	if err != nil {
		return nil, false, err
	}
	err = nil

	return user, false, nil
}

func (s *Service) Update(ctx context.Context, req request.UpdateUser, id request.GetByIdUser) (*model.User, bool, error) {
	// ตรวจสอบว่า user มีอยู่หรือไม่
	ex, err := s.db.NewSelect().Table("users").Where("id = ?", id.ID).Exists(ctx)
	if err != nil {
		return nil, false, err
	}
	if !ex {
		return nil, false, errors.New("user not found")
	}

	// หา position_id จาก position_name
	var positionID string
	err = s.db.NewSelect().
		Model((*model.Position)(nil)).
		Column("id").
		Where("name = ?", req.Position_Name).
		Scan(ctx, &positionID)
	if err != nil {
		return nil, false, errors.New("ไม่พบตำแหน่งงานที่ระบุ")
	}

	// เข้ารหัสรหัสผ่านใหม่
	bytes, err := bcrypt.GenerateFromPassword([]byte(req.Password), 10)
	if err != nil {
		return nil, false, err
	}

	// เตรียมข้อมูล user สำหรับอัปเดต
	m := &model.User{
		ID:          id.ID,
		FirstName:   req.FirstName,
		LastName:    req.LastName,
		Email:       req.Email,
		Password:    string(bytes),
		Position_ID: positionID,
		Image_url:   req.Image_url,
		Phone:       req.Phone,
	}
	m.SetUpdateNow()

	// อัปเดตข้อมูลในฐานข้อมูล
	_, err = s.db.NewUpdate().Model(m).
		Set("first_name = ?first_name").
		Set("last_name = ?last_name").
		Set("email = ?email").
		Set("password = ?password").
		Set("position_id = ?position_id").
		Set("image_url = ?image_url").
		Set("phone = ?phone").
		Set("updated_at = ?updated_at").
		WherePK().
		OmitZero().
		Returning("*").
		Exec(ctx)
	if err != nil {
		if strings.Contains(err.Error(), "duplicate key value") {
			return nil, true, errors.New("user already exists")
		}
		return nil, false, err
	}

	return m, false, nil
}

func (s *Service) List(ctx context.Context, req request.ListUser) ([]response.ListUser, int, error) {
	offset := (req.Page - 1) * req.Size
	m := []response.ListUser{}

	query := s.db.NewSelect().
		TableExpr("users as u").
		Column("u.id", "u.first_name", "u.last_name", "u.email", "u.position_id", "u.image_url", "u.phone", "u.created_at", "u.updated_at").Where("deleted_at IS NULL")

	// Filtering
	if req.Search != "" {
		search := "%" + strings.ToLower(req.Search) + "%"
		if req.SearchBy != "" {
			searchBy := strings.ToLower(req.SearchBy)
			query = query.Where(fmt.Sprintf("LOWER(u.%s) LIKE ?", searchBy), search)
		} else {
			query = query.Where("LOWER(u.first_name) LIKE ?", search)
		}
	}

	// Count total before pagination
	count, err := query.Count(ctx)
	if err != nil {
		return nil, 0, err
	}

	// Order handling
	order := fmt.Sprintf("u.%s %s", req.SortBy, req.OrderBy)

	// Final query with order + pagination
	err = query.Order(order).Limit(req.Size).Offset(offset).Scan(ctx, &m)
	if err != nil {
		return nil, 0, err
	}

	return m, count, nil
}

func (s *Service) Get(ctx context.Context, id request.GetByIdUser) (*response.ListUser, error) {
	m := response.ListUser{}

	err := s.db.NewSelect().
		TableExpr("users as u").
		Join("LEFT JOIN positions as p ON p.id = u.position_id::uuid AND p.deleted_at IS NULL").
		ColumnExpr("u.id").
		ColumnExpr("u.first_name").
		ColumnExpr("u.last_name").
		ColumnExpr("u.email").
		ColumnExpr("u.position_id").
		ColumnExpr("p.name AS position_name").
		ColumnExpr("u.image_url").
		ColumnExpr("u.phone").
		ColumnExpr("u.created_at").
		ColumnExpr("u.updated_at").
		Where("u.id = ?", id.ID).
		Where("u.deleted_at IS NULL").
		Scan(ctx, &m)

	return &m, err
}

func (s *Service) Delete(ctx context.Context, id request.GetByIdUser) error {
	ex, err := s.db.NewSelect().Table("users").Where("id = ?", id.ID).Where("deleted_at IS NULL").Exists(ctx)
	if err != nil {
		return err
	}
	if !ex {
		return errors.New("user not found")
	}

	_, err = s.db.NewDelete().Model((*model.User)(nil)).Where("id = ?", id.ID).Exec(ctx)
	return err
}
