package equipment

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"app/app/enum"
	"app/app/model"
	"app/app/request"
)

func (s *Service) Create(ctx context.Context, req request.CreateEquipment) (*model.Equipment, bool, error) {
	m := model.Equipment{
		Name:               req.Name,
		Image_URL:          req.Image_URL,
		Quantity:           req.Quantity,
		Status:             enum.EquipmentAvaliable,
	}
	_, err := s.db.NewInsert().Model(&m).Exec(ctx)
	if err != nil {
		if strings.Contains(err.Error(), "duplicate key value") {
			return nil, true, errors.New("equipment already exists")
		}
		return nil, false, err
	}
	return &m, false, nil
}

func (s *Service) Update(ctx context.Context, req request.UpdateEquipment, id string) (*model.Equipment, bool, error) {
	// ดึงข้อมูลเก่า
	var m model.Equipment
	ex, err := s.db.NewSelect().Model(&m).Where("id = ?", id).Exists(ctx)
	if err != nil {
		return nil, false, err
	}
	if !ex {
		return nil, true, errors.New("equipment not found")
	}

	if err := s.db.NewSelect().Model(&m).Where("id = ?", id).Scan(ctx); err != nil {
		return nil, false, err
	}

	// อัปเดตเฉพาะ field ที่ส่งมา
	if req.Name != nil {
		m.Name = *req.Name
	}
	if req.Image_URL != nil {
		m.Image_URL = *req.Image_URL
	}
	if req.Quantity != nil {
		m.Quantity = *req.Quantity
	}
	if req.Status != nil {
		m.Status = enum.EquipmentStatus(*req.Status)
	}

	m.SetUpdateNow()

	res, err := s.db.NewUpdate().Model(&m).Where("id = ?", id).Exec(ctx)
	if err != nil {
		if strings.Contains(err.Error(), "duplicate key value") {
			return nil, true, errors.New("equipment already exists")
		}
		return nil, false, err
	}
	affected, _ := res.RowsAffected()
	if affected == 0 {
		return nil, true, errors.New("equipment not updated")
	}

	return &m, false, nil
}

func (s *Service) Delete(ctx context.Context, id string) (*model.Equipment, bool, error) {
	var m model.Equipment
	ex, err := s.db.NewSelect().Model(&m).Where("id = ?", id).Exists(ctx)
	if err != nil {
		return nil, false, err
	}
	if !ex {
		return nil, true, errors.New("equipment not found")
	}

	_, err = s.db.NewDelete().Model(&m).Where("id = ?", id).Exec(ctx)
	return &m, false, err
}

func (s *Service) Get(ctx context.Context, id string) (*model.Equipment, error) {
	var m model.Equipment
	err := s.db.NewSelect().Model(&m).Where("id = ?", id).Scan(ctx)
	return &m, err
}

func (s *Service) List(ctx context.Context, req request.ListEquipment) ([]model.Equipment, int, error) {
	offset := (req.Page - 1) * req.Size
	m := []model.Equipment{}

	query := s.db.NewSelect().Model(&m)
	if req.Search != "" {
		query.Where("LOWER(name) LIKE ?", "%"+strings.ToLower(req.Search)+"%")
	}

	count, err := query.Count(ctx)
	if err != nil {
		return nil, 0, err
	}

	order := fmt.Sprintf("%s %s", req.SortBy, req.OrderBy)
	err = query.Offset(offset).Limit(req.Size).Order(order).Scan(ctx, &m)
	return m, count, err
}


func (s *Service) GetEquipmentAvailable(ctx context.Context, equipmentID string) (int, error) {
    var eq model.Equipment
    err := s.db.NewSelect().Model(&eq).Where("id = ?", equipmentID).Scan(ctx)
    if err != nil {
        return 0, err
    }

    var usedQty int
    err = s.db.NewSelect().
        Table("booking_equipments").
        Join("JOIN bookings as b ON booking_equipments.booking_id::uuid = b.id::uuid").
        Where("booking_equipments.equipment_id = ?", equipmentID).
        Where("b.status IN (?, ?)", "Pending", "Approved").
        ColumnExpr("COALESCE(SUM(booking_equipments.quantity), 0)").
        Scan(ctx, &usedQty)
    if err != nil {
        return 0, err
    }

    available := eq.Quantity - usedQty
    if available < 0 {
        available = 0
    }
    return available, nil
}