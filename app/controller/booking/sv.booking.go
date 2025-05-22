package booking

import (
	"app/app/enum"
	"app/app/model"
	"app/app/request"
	"app/app/response"
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/uptrace/bun"
)

type BookingService struct {
	db *bun.DB
}

func (s *Service) Create(ctx context.Context, req request.CreateBooking) (*model.Booking, bool, error) {

	startUnix := req.StartTime
	endUnix := req.EndTime

	exists, err := s.db.NewSelect().
		Model((*model.Booking)(nil)).
		Where("room_id = ?", req.RoomID).
		Where("start_time < ?", startUnix).
		Where("end_time > ?", endUnix).
		Where("deleted_at IS NULL").
		Exists(ctx)
	if err != nil {
		return nil, false, err
	}
	if exists {
		return nil, true, errors.New("booking time conflict")
	}

	// ถ้าไม่มี conflict, proceed
	m := model.Booking{
		UserID:      req.UserID,
		RoomID:      req.RoomID,
		Title:       req.Title,
		Description: req.Description,
		Phone:       req.Phone,
		StartTime:   startUnix,
		EndTime:     endUnix,
		ApprovedBy:  "",
		Status:      enum.BookingStatus("Pending"),
	}

	_, err = s.db.NewInsert().Model(&m).Exec(ctx)
	if err != nil {
		return nil, false, err
	}

	return &m, false, nil
}

func (s *Service) Update(ctx context.Context, req request.UpdateBooking, id request.GetByIdBooking) (*model.Booking, bool, error) {
	// ตรวจสอบว่า booking มีอยู่ไหม
	ex, err := s.db.NewSelect().Table("bookings").Where("id = ?", id.ID).Exists(ctx)
	if err != nil {
		return nil, false, err
	}
	if !ex {
		return nil, false, errors.New("booking not found")
	}

	// เตรียม model
	m := &model.Booking{ID: id.ID}
	q := s.db.NewUpdate().Model(m).WherePK().OmitZero().Returning("*")

	// อัปเดตฟิลด์ทั่วไป
	if req.RoomID != "" {
		m.RoomID = req.RoomID
		q.Set("room_id = ?room_id")
	}
	if req.Description != "" {
		m.Description = req.Description
		q.Set("description = ?description")
	}
	if req.Phone != "" {
		m.Phone = req.Phone
		q.Set("phone = ?phone")
	}
	if req.Title != "" {
		m.Title = req.Title
		q.Set("title = ?title")
	}
	if req.StartTime > 0 {
		m.StartTime = req.StartTime
		q.Set("start_time = ?start_time")
	}
	if req.EndTime > 0 {
		m.EndTime = req.EndTime
		q.Set("end_time = ?end_time")
	}

	if req.Status != "" {
		m.Status = enum.BookingStatus(req.Status)
		q.Set("status = ?status")

		if req.Status == string(enum.BookingApproved) || req.Status == string(enum.BookingCanceled) {
			userID, ok := ctx.Value("user_id").(string)
			if !ok || userID == "" {
				return nil, false, errors.New("unauthorized: user_id not found")
			}
			m.ApprovedBy = userID
			q.Set("approved_by = ?approved_by")
		}
	}

	// อัปเดตเวลาล่าสุด
	m.SetUpdateNow()
	q.Set("updated_at = ?updated_at")

	_, err = q.Exec(ctx)
	if err != nil {
		if strings.Contains(err.Error(), "duplicate key value") {
			return nil, true, errors.New("booking already exists")
		}
		return nil, false, err
	}

	return m, false, nil
}

func (s *Service) List(ctx context.Context, req request.ListBooking) ([]response.BookingResponse, int, error) {
	
	offset := (req.Page - 1) * req.Size
	m := []response.BookingResponse{}
	baseQuery := s.db.NewSelect().
		TableExpr("bookings as b").
		ColumnExpr("b.id as id").
		ColumnExpr("u.id as user_id").
		ColumnExpr("u.first_name as user_name").
		ColumnExpr("u.last_name as user_lastname").
		ColumnExpr("r.id as room_id").
		ColumnExpr("r.name as room_name").
		ColumnExpr("b.title as title").
		ColumnExpr("b.description as description").
		ColumnExpr("b.phone as phone").
		ColumnExpr("b.start_time as start_time").
		ColumnExpr("b.end_time as end_time").
		ColumnExpr("b.status as status").
		ColumnExpr("b.approved_by as approved_by").
		ColumnExpr("COALESCE(au.first_name || ' ' || au.last_name, '') as nameapproved_by").
		ColumnExpr("b.created_at as created_at").
		ColumnExpr("b.updated_at as updated_at").
		Join("JOIN users as u ON b.user_id::uuid = u.id").
		Join("JOIN rooms as r ON b.room_id::uuid = r.id").
		Join("LEFT JOIN users as au ON CAST(NULLIF(b.approved_by, '') AS uuid) = au.id").
		Where("b.deleted_at IS NULL").
		OrderExpr("b.created_at ASC")

	if req.Search != "" {
		search := "%" + strings.ToLower(req.Search) + "%"
		if req.SearchBy != "" {
			searchBy := strings.ToLower(req.SearchBy)
			baseQuery = baseQuery.Where(fmt.Sprintf("LOWER(b.%s) LIKE ?", searchBy), search)
		} else {
			baseQuery = baseQuery.Where("LOWER(b.title) LIKE ?", search)
		}
	}

	countQuery := s.db.NewSelect().
		TableExpr("bookings as b").
		ColumnExpr("COUNT(*)").
		Join("JOIN users as u ON b.user_id::uuid = u.id").
		Join("JOIN rooms as r ON b.room_id::uuid = r.id").
		Where("b.deleted_at IS NULL")

	if req.Search != "" {
		search := "%" + strings.ToLower(req.Search) + "%"
		if req.SearchBy != "" {
			searchBy := strings.ToLower(req.SearchBy)
			countQuery = countQuery.Where(fmt.Sprintf("LOWER(b.%s) LIKE ?", searchBy), search)
		} else {
			countQuery = countQuery.Where("LOWER(b.title) LIKE ?", search)
		}
	}

	var count int
	err := countQuery.Scan(ctx, &count)
	if err != nil {
		return nil, 0, err
	}

	order := fmt.Sprintf("b.%s %s", req.SortBy, req.OrderBy)
	err = baseQuery.Order(order).Limit(req.Size).Offset(offset).Scan(ctx, &m)
	if err != nil {
		return nil, 0, err
	}

	return m, count, nil
}

func (s *Service) ListHistory(ctx context.Context, req request.ListBooking) ([]response.BookingResponse, int, error) {
	offset := (req.Page - 1) * req.Size
	m := []response.BookingResponse{}

	baseQuery := s.db.NewSelect().
		TableExpr("bookings as b").
		ColumnExpr("b.id as id").
		ColumnExpr("u.id as user_id").
		ColumnExpr("u.first_name as user_name").
		ColumnExpr("u.last_name as user_lastname").
		ColumnExpr("r.id as room_id").
		ColumnExpr("r.name as room_name").
		ColumnExpr("b.title as title").
		ColumnExpr("b.description as description").
		ColumnExpr("b.phone as phone").
		ColumnExpr("b.start_time as start_time").
		ColumnExpr("b.end_time as end_time").
		ColumnExpr("b.status as status").
		ColumnExpr("b.approved_by as approved_by").
		ColumnExpr("COALESCE(au.first_name || ' ' || au.last_name, '') as nameapproved_by").
		ColumnExpr("b.created_at as created_at").
		ColumnExpr("b.updated_at as updated_at").
		Join("JOIN users as u ON b.user_id::uuid = u.id").
		Join("JOIN rooms as r ON b.room_id::uuid = r.id").
		Join("LEFT JOIN users as au ON CAST(NULLIF(b.approved_by, '') AS uuid) = au.id").
		Where("b.deleted_at IS NULL").
		OrderExpr("b.created_at ASC")

	if req.Search != "" {
		search := "%" + strings.ToLower(req.Search) + "%"
		if req.SearchBy != "" {
			searchBy := strings.ToLower(req.SearchBy)
			baseQuery = baseQuery.Where(fmt.Sprintf("LOWER(b.%s) LIKE ?", searchBy), search)
		} else {
			baseQuery = baseQuery.Where("LOWER(b.title) LIKE ?", search)
		}
	}

	countQuery := s.db.NewSelect().
		TableExpr("bookings as b").
		ColumnExpr("COUNT(*)").
		Join("JOIN users as u ON b.user_id::uuid = u.id").
		Join("JOIN rooms as r ON b.room_id::uuid = r.id").
		Where("b.deleted_at IS NULL")

	if req.Search != "" {
		search := "%" + strings.ToLower(req.Search) + "%"
		if req.SearchBy != "" {
			searchBy := strings.ToLower(req.SearchBy)
			countQuery = countQuery.Where(fmt.Sprintf("LOWER(b.%s) LIKE ?", searchBy), search)
		} else {
			countQuery = countQuery.Where("LOWER(b.title) LIKE ?", search)
		}
	}

	var count int
	err := countQuery.Scan(ctx, &count)
	if err != nil {
		return nil, 0, err
	}

	order := fmt.Sprintf("b.%s %s", req.SortBy, req.OrderBy)
	err = baseQuery.Order(order).Limit(req.Size).Offset(offset).Scan(ctx, &m)
	if err != nil {
		return nil, 0, err
	}

	return m, count, nil
}

func (s *Service) Get(ctx context.Context, id request.GetByIdBooking) (*response.BookingResponse, error) {
	m := response.BookingResponse{}

	err := s.db.NewSelect().
		TableExpr("bookings as b").
		ColumnExpr("b.id as id").
		ColumnExpr("u.id as user_id").
		ColumnExpr("u.first_name as user_name").
		ColumnExpr("u.last_name as user_lastname").
		ColumnExpr("r.id as room_id").
		ColumnExpr("r.name as room_name").
		ColumnExpr("b.title as title").
		ColumnExpr("b.description as description").
		ColumnExpr("b.phone as phone").
		ColumnExpr("b.start_time as start_time").
		ColumnExpr("b.end_time as end_time").
		ColumnExpr("b.status as status").
		ColumnExpr("b.approved_by as approved_by").
		ColumnExpr("COALESCE(au.first_name || ' ' || au.last_name, '') as nameapproved_by").
		ColumnExpr("b.updated_at as updated_at").
		Join("JOIN users as u ON b.user_id::uuid = u.id").
		Join("JOIN rooms as r ON b.room_id::uuid = r.id").
		Join("LEFT JOIN users as au ON CAST(NULLIF(b.approved_by, '') AS uuid) = au.id").
		Where("b.deleted_at IS NULL").
		OrderExpr("b.created_at ASC").
		Scan(ctx, &m)
	return &m, err
}

func (s *Service) GetByRoomId(ctx context.Context, req request.GetByRoomIdBooking) ([]response.BookingResponse, int, error) {
	// เพิ่มการรับค่า Page และ Size จาก req
	offset := (req.Page - 1) * req.Size
	m := []response.BookingResponse{}

	query := s.db.NewSelect().
		TableExpr("bookings as b").
		ColumnExpr("b.id as id").
		ColumnExpr("u.id as user_id").
		ColumnExpr("u.first_name as user_name").
		ColumnExpr("u.last_name as user_lastname").
		ColumnExpr("r.id as room_id").
		ColumnExpr("r.name as room_name").
		ColumnExpr("b.title as title").
		ColumnExpr("b.description as description").
		ColumnExpr("b.phone as phone").
		ColumnExpr("b.start_time as start_time").
		ColumnExpr("b.end_time as end_time").
		ColumnExpr("b.status as status").
		ColumnExpr("b.approved_by as approved_by").
		ColumnExpr("COALESCE(au.first_name || ' ' || au.last_name, '') as nameapproved_by").
		ColumnExpr("b.created_at as created_at").
		ColumnExpr("b.updated_at as updated_at").
		Join("JOIN users as u ON b.user_id::uuid = u.id").
		Join("JOIN rooms as r ON b.room_id::uuid = r.id").
		Join("LEFT JOIN users as au ON CAST(NULLIF(b.approved_by, '') AS uuid) = au.id").
		Where("b.deleted_at IS NULL").
		Where("r.id = ?", req.RoomID).
		OrderExpr("start_time ASC")

	count, err := query.Count(ctx)
	if err != nil {
		return nil, 0, err
	}

	err = query.OrderExpr("start_time ASC").Limit(req.Size).Offset(offset).Scan(ctx, &m)
	if err != nil {
		return nil, 0, err
	}

	return m, count, nil
}

func (s *Service) GetBookingByUserID(ctx context.Context, id request.GetByIdUser) ([]response.BookingbyUser, error) {
	var bookings []response.BookingbyUser

	err := s.db.NewSelect().
		TableExpr("bookings as b").
		ColumnExpr("b.id as id").
		ColumnExpr("u.first_name as user_name").
		ColumnExpr("u.last_name as user_lastname").
		ColumnExpr("r.id as room_id").
		ColumnExpr("r.name as room_name").
		ColumnExpr("b.title as title").
		ColumnExpr("b.description as description").
		ColumnExpr("b.phone as phone").
		ColumnExpr("b.start_time as start_time").
		ColumnExpr("b.end_time as end_time").
		ColumnExpr("b.status as status").
		ColumnExpr("b.approved_by as approved_by").
		ColumnExpr("COALESCE(au.first_name || ' ' || au.last_name, '') as nameapproved_by").
		ColumnExpr("b.created_at as created_at").
		ColumnExpr("b.updated_at as updated_at").
		ColumnExpr("b.deleted_at as deleted_at").
		Join("JOIN users as u ON b.user_id::uuid = u.id").
		Join("JOIN rooms as r ON b.room_id::uuid = r.id").
		Join("LEFT JOIN users as au ON CAST(NULLIF(b.approved_by, '') AS uuid) = au.id").
		Where("b.user_id = ?", id.ID).
		Where("b.deleted_at IS NULL").
		OrderExpr("b.created_at ASC").
		Scan(ctx, &bookings)
	return bookings, err
}

func (s *Service) GetBookingHistoryByUserID(ctx context.Context, id request.GetByIdUser) ([]response.BookingbyUser, error) {
	var bookings []response.BookingbyUser

	err := s.db.NewSelect().
		TableExpr("bookings as b").
		ColumnExpr("b.id as id").
		ColumnExpr("u.first_name as user_name").
		ColumnExpr("u.last_name as user_lastname").
		ColumnExpr("r.id as room_id").
		ColumnExpr("r.name as room_name").
		ColumnExpr("b.title as title").
		ColumnExpr("b.description as description").
		ColumnExpr("b.phone as phone").
		ColumnExpr("b.start_time as start_time").
		ColumnExpr("b.end_time as end_time").
		ColumnExpr("b.status as status").
		ColumnExpr("b.approved_by as approved_by").
		ColumnExpr("COALESCE(au.first_name || ' ' || au.last_name, '') as nameapproved_by").
		ColumnExpr("b.updated_at as updated_at").
		Join("JOIN users as u ON b.user_id::uuid = u.id").
		Join("JOIN rooms as r ON b.room_id::uuid = r.id").
		Join("LEFT JOIN users as au ON CAST(NULLIF(b.approved_by, '') AS uuid) = au.id").
		Where("b.user_id = ?", id.ID).
		Where("b.deleted_at IS NULL").
		OrderExpr("b.created_at ASC").
		Scan(ctx, &bookings)
	return bookings, err
}

func (s *Service) Delete(ctx context.Context, id request.GetByIdBooking) error {
	ex, err := s.db.NewSelect().Table("bookings").Where("id = ?", id.ID).Where("deleted_at IS NULL").Exists(ctx)
	if err != nil {
		return err
	}
	if !ex {
		return errors.New("booking not found")
	}

	_, err = s.db.NewDelete().Model((*model.Booking)(nil)).Where("id = ?", id.ID).Exec(ctx)
	return err
}

func NewBookingService(db *bun.DB) *BookingService {
	return &BookingService{db: db}
}

func (s *BookingService) AutoExpiredBookings() error {
	now := time.Now().Unix()

	_, err := s.db.NewUpdate().
		Model((*model.Booking)(nil)).
		Set("status = ?", enum.BookingFinished).
		Set("updated_at = ?", now).
		Where("end_time < ?", now).
		Where("status = ?", enum.BookingApproved).
		Where("deleted_at IS NULL").
		Exec(context.Background())

	return err
}
