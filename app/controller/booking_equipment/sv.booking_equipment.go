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
    var inserted []model.BookingEquipment

    for _, eq := range equipments {
        be := model.BookingEquipment{
            BookingID:   bookingID,
            EquipmentID: eq.EquipmentID,
            Quantity:    eq.Quantity,
        }
        _, err := s.db.NewInsert().Model(&be).Exec(ctx)
        if err != nil {
            return nil, err
        }
        inserted = append(inserted, be)
    }

    return inserted, nil
}

func (s *Service) Update(ctx context.Context, req request.UpdateEquipment, id string) (*model.Equipment, bool, error) {
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

func (s *Service) List(ctx context.Context, req request.ListEquipmentBooking) ([]model.BookingEquipment, int, error) {
    var data []model.BookingEquipment
    query := s.db.NewSelect().Model(&data)
    if req.SortBy != "" {
        query.Order(fmt.Sprintf("%s %s", req.SortBy, req.OrderBy))
    }
    if req.Size > 0 {
        query.Limit(req.Size).Offset((req.Page - 1) * req.Size)
    }
    count, err := query.ScanAndCount(ctx)
    if err != nil {
        return nil, 0, err
    }
    return data, count, nil
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
        Join("JOIN bookings as b ON booking_equipments.booking_id = b.id").
        Where("booking_equipments.equipment_id = ?", equipmentID).
        Where("b.status IN (?, ?)", "Pending", "Approved").
        ColumnExpr("COALESCE(SUM(booking_equipments.quantity), 0)").
        Scan(ctx, &usedQty)
    if err != nil {
        return 0, err
    }

    fmt.Println("EQ:", eq.Quantity, "USED:", usedQty)

    available := eq.Quantity - usedQty
    if available < 0 {
        available = 0
    }
    return available, nil
}