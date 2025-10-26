package user

import (
	"app/app/model"
	"app/app/request"
	"app/app/response"
	"context"
	"database/sql"
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
	// 1) ตรวจว่ามีผู้ใช้จริง
	user := new(model.User)
	if err := s.db.NewSelect().
		Model(user).
		Where("id = ?", id.ID).
		Scan(ctx); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, false, errors.New("user not found")
		}
		return nil, false, err
	}

	// 2) เก็บคอลัมน์ที่จะอัปเดตแบบ dynamic
	updateCols := make([]string, 0)

	// 3) อัปเดตทีละฟิลด์ ถ้ามีค่ามา (trim space ก่อน)
	trim := func(s string) string { return strings.TrimSpace(s) }

	if v := trim(req.FirstName); v != "" {
		user.FirstName = v
		updateCols = append(updateCols, "first_name")
	}
	if v := trim(req.LastName); v != "" {
		user.LastName = v
		updateCols = append(updateCols, "last_name")
	}
	if v := trim(req.Email); v != "" {
		user.Email = v
		updateCols = append(updateCols, "email")
	}
	if v := trim(req.Phone); v != "" {
		user.Phone = v
		updateCols = append(updateCols, "phone")
	}
	if v := trim(req.Image_url); v != "" {
		user.Image_url = v
		updateCols = append(updateCols, "image_url")
	}

	// 4) position_name: หา position_id เฉพาะกรณีส่งมาและไม่ว่าง
	if v := trim(req.Position_Name); v != "" {
		var positionID string
		if err := s.db.NewSelect().
			Model((*model.Position)(nil)).
			Column("id").
			Where("name = ?", v).
			Scan(ctx, &positionID); err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return nil, true, errors.New("ไม่พบตำแหน่งงานที่ระบุ")
			}
			return nil, false, err
		}
		user.Position_ID = positionID
		updateCols = append(updateCols, "position_id")
	}

	// 5) password: hash เฉพาะเมื่อส่งมาและไม่ว่าง
	if v := trim(req.Password); v != "" {
		hashed, err := bcrypt.GenerateFromPassword([]byte(v), 10)
		if err != nil {
			return nil, false, err
		}
		user.Password = string(hashed)
		updateCols = append(updateCols, "password")
	}

	// 6) ถ้าไม่มีฟิลด์ไหนจะอัปเดตเลย ก็คืนค่ากลับ (ถือว่าสำเร็จแบบ no-op)
	if len(updateCols) == 0 {
		return user, false, nil
	}

	// อัปเดตเวลา
	user.SetUpdateNow()
	updateCols = append(updateCols, "updated_at")

	// 7) รันอัปเดตเฉพาะคอลัมน์ที่กำหนด
	_, err := s.db.NewUpdate().
		Model(user).
		Column(updateCols...).
		WherePK().
		Returning("*").
		Exec(ctx)
	if err != nil {
		// ตรวจ duplicate key (email/phone ซ้ำ)
		if strings.Contains(err.Error(), "duplicate key value") {
			return nil, true, errors.New("user already exists")
		}
		return nil, false, err
	}

	return user, false, nil
}

func (s *Service) List(ctx context.Context, req request.ListUser) ([]response.ListUser, int, error) {
	offset := (req.Page - 1) * req.Size
	m := []response.ListUser{}

	// Base query
	baseQuery := s.db.NewSelect().
		TableExpr("users as u").
		ColumnExpr("u.id").
		ColumnExpr("u.first_name").
		ColumnExpr("u.last_name").
		ColumnExpr("u.email").
		ColumnExpr("u.position_id").
		ColumnExpr("p.name AS position_name").
		ColumnExpr("r.name AS role_name").
		ColumnExpr("u.image_url").
		ColumnExpr("u.phone").
		ColumnExpr("u.created_at").
		ColumnExpr("u.updated_at").
		Join("LEFT JOIN positions AS p ON u.position_id::uuid = p.id::uuid").
		Join("LEFT JOIN user_roles AS ur ON u.id::uuid = ur.user_id::uuid").
		Join("LEFT JOIN roles AS r ON ur.role_id::uuid = r.id::uuid").
		Where("u.deleted_at IS NULL").
		OrderExpr("u.created_at ASC")

	// Filtering
	if req.Search != "" {
		search := "%" + strings.ToLower(req.Search) + "%"
		if req.SearchBy != "" {
			searchBy := strings.ToLower(req.SearchBy)
			baseQuery = baseQuery.Where(fmt.Sprintf("LOWER(u.%s) LIKE ?", searchBy), search)
		} else {
			baseQuery = baseQuery.Where("LOWER(u.first_name) LIKE ?", search)
		}
	}

	// Count query
	countQuery := s.db.NewSelect().
		TableExpr("users as u").
		ColumnExpr("COUNT(*)").
		Join("LEFT JOIN positions AS p ON u.position_id::uuid = p.id::uuid").
		Join("LEFT JOIN user_roles AS ur ON u.id::uuid = ur.user_id::uuid").
		Join("LEFT JOIN roles AS r ON ur.role_id::uuid = r.id::uuid").
		Where("u.deleted_at IS NULL")

	if req.Search != "" {
		search := "%" + strings.ToLower(req.Search) + "%"
		if req.SearchBy != "" {
			searchBy := strings.ToLower(req.SearchBy)
			countQuery = countQuery.Where(fmt.Sprintf("LOWER(u.%s) LIKE ?", searchBy), search)
		} else {
			countQuery = countQuery.Where("LOWER(u.first_name) LIKE ?", search)
		}
	}

	var count int
	err := countQuery.Scan(ctx, &count)
	if err != nil {
		return nil, 0, err
	}

	order := fmt.Sprintf("u.%s %s", req.SortBy, req.OrderBy)
	err = baseQuery.Order(order).Limit(req.Size).Offset(offset).Scan(ctx, &m)
	if err != nil {
		return nil, 0, err
	}

	return m, count, nil
}

func (s *Service) Get(ctx context.Context, id request.GetByIdUser) (*response.ListUser, error) {
	m := response.ListUser{}

	err := s.db.NewSelect().
		TableExpr("users as u").
		ColumnExpr("u.id").
		ColumnExpr("u.first_name").
		ColumnExpr("u.last_name").
		ColumnExpr("u.email").
		ColumnExpr("u.position_id").
		ColumnExpr("p.name AS position_name").
		ColumnExpr("r.name AS role_name").
		ColumnExpr("u.image_url").
		ColumnExpr("u.phone").
		ColumnExpr("u.created_at").
		ColumnExpr("u.updated_at").
		Join("LEFT JOIN positions AS p ON u.position_id::uuid = p.id::uuid").
		Join("LEFT JOIN user_roles AS ur ON u.id::uuid = ur.user_id::uuid").
		Join("LEFT JOIN roles AS r ON ur.role_id::uuid = r.id::uuid").
		Where("u.id = ?", id.ID).
		Where("u.deleted_at IS NULL").
		OrderExpr("u.created_at ASC").
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
