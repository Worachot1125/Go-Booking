package booking_equipment

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"app/app/enum"
	"app/app/model"
	"app/app/request"
)

func (s *Service) Create(ctx context.Context, bookingID string, equipments []request.EquipmentSelection) ([]model.BookingEquipment, error) {
    var result []model.BookingEquipment

    for _, eq := range equipments {
        // 1. ดึงข้อมูล equipment
        var equipment model.Equipment
        err := s.db.NewSelect().Model(&equipment).Where("id = ?", eq.EquipmentID).Scan(ctx)
        if err != nil {
            return nil, errors.New("equipment not found")
        }

        // 2. เช็คจำนวน
        if equipment.Available_Quantity < eq.Quantity || eq.Quantity <= 0 {
            return nil, errors.New("equipment " + equipment.Name + " is not available in requested quantity")
        }

        // 3. ลด available_​quantity
        equipment.Available_Quantity -= eq.Quantity
        _, err = s.db.NewUpdate().
            Model(&equipment).
            Set("available_quantity = ?", equipment.Available_Quantity).
            Where("id = ?", eq.EquipmentID).
            Exec(ctx)
        if err != nil {
            return nil, err
        }

        // 4. สร้าง booking_equipment
		be := model.BookingEquipment{
			BookingID:   &model.Booking{ID: bookingID},
			EquipmentID: &model.Equipment{ID: eq.EquipmentID},
			Quantity:    eq.Quantity,
		}
        _, err = s.db.NewInsert().Model(&be).Exec(ctx)
        if err != nil {
            return nil, err
        }
        result = append(result, be)
    }
    return result, nil
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
	if req.Available_Quantity != nil {
		m.Available_Quantity = *req.Available_Quantity
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
