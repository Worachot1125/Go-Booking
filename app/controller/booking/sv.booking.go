package booking

import (
	linectrl "app/app/controller/line"
	"app/app/enum"
	"app/app/helper"
	"app/app/model"
	"app/app/request"
	"app/app/response"
	"app/internal/logger"
	"context"
	"errors"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/line/line-bot-sdk-go/v7/linebot"
	"github.com/uptrace/bun"
)

type BookingService struct {
	db *bun.DB
}

func (s *Service) Create(ctx context.Context, req request.CreateBooking) (*model.Booking, bool, error) {
	startUnix := req.StartTime
	endUnix := req.EndTime

	// ✅ ตรวจชนกันแบบครอบคลุมทุกกรณีของการทับช่วงเวลา
	exists, err := s.db.NewSelect().
		Model((*model.Booking)(nil)).
		Where("room_id = ?", req.RoomID).
		Where("deleted_at IS NULL").
		Where("(start_time < ?) AND (end_time > ?)", endUnix, startUnix).
		Exists(ctx)
	if err != nil {
		return nil, false, err
	}
	if exists {
		return nil, true, errors.New("booking time conflict")
	}

	// สร้าง booking
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

	var info struct {
		LineUserID string
		RoomName   string
	}
	_ = s.db.NewSelect().
		TableExpr("users AS u").
		Join("JOIN rooms AS r ON r.id = ?", req.RoomID).
		ColumnExpr("u.line_user_id, r.name AS room_name").
		Where("u.id = ?", req.UserID).
		Where("u.line_opt_in = TRUE").
		Scan(ctx, &info)

	// ส่งแจ้งเตือนถ้ามี LINE
	if info.LineUserID != "" {
		loc, _ := time.LoadLocation("Asia/Bangkok")
		msg := fmt.Sprintf(
			"📅 สร้างการจองสำเร็จ\nเรื่อง: %s\nห้อง: %s\nเวลาเริ่ม: %s\nเวลาสิ้นสุด: %s",
			m.Title,
			info.RoomName,
			time.Unix(m.StartTime, 0).In(loc).Format("02/01 15:04"),
			time.Unix(m.EndTime, 0).In(loc).Format("02/01 15:04"),
		)
		_ = linectrl.PushTextToLineUser(info.LineUserID, msg)
	}

	return &m, false, nil
}

func (s *Service) Update(ctx context.Context, req request.UpdateBooking, id request.GetByIdBooking) (*model.Booking, bool, error) {
	// 1) มี booking ไหม
	ex, err := s.db.NewSelect().Table("bookings").Where("id = ?", id.ID).Exists(ctx)
	if err != nil {
		return nil, false, err
	}
	if !ex {
		return nil, false, errors.New("booking not found")
	}

	// 2) prev
	prev := new(model.Booking)
	if err := s.db.NewSelect().Model(prev).Where("id = ?", id.ID).Scan(ctx); err != nil {
		return nil, false, err
	}

	// 3) update fields
	m := &model.Booking{ID: id.ID}
	q := s.db.NewUpdate().Model(m).WherePK().OmitZero().Returning("*")

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

	// 3.1) สถานะ
	var statusProvided bool
	var newStatus enum.BookingStatus
	rawStatus := strings.TrimSpace(req.Status)
	if rawStatus != "" {
		statusProvided = true
		switch strings.ToLower(rawStatus) {
		case "approved":
			newStatus = enum.BookingApproved
		case "canceled", "cancelled":
			newStatus = enum.BookingCanceled
		case "pending":
			newStatus = enum.BookingPending
		case "finished":
			newStatus = enum.BookingFinished
		default:
			newStatus = enum.BookingStatus(req.Status)
		}
		m.Status = newStatus
		q.Set("status = ?status")

		if newStatus == enum.BookingApproved || newStatus == enum.BookingCanceled {
			userID, ok := ctx.Value("user_id").(string)
			if !ok || strings.TrimSpace(userID) == "" {
				logger.Errf("[booking:%s] missing ctx user_id for status=%s", id.ID, newStatus)
				return nil, false, errors.New("unauthorized: user_id not found")
			}
			m.ApprovedBy = userID
			q.Set("approved_by = ?approved_by")
		}
	}

	// 3.2) updated_at
	m.SetUpdateNow()
	q.Set("updated_at = ?updated_at")

	// 4) exec
	if _, err = q.Exec(ctx); err != nil {
		if strings.Contains(err.Error(), "duplicate key value") {
			return nil, true, errors.New("booking already exists")
		}
		return nil, false, err
	}

	logger.Infof("[booking:%s] update ok | prev=%s -> new=%s | statusProvided=%v",
		id.ID, prev.Status, m.Status, statusProvided)

	// 5) push เฉพาะเมื่อส่งสถานะมา
	if statusProvided {
		// ดึง line_user_id
		var lineUserID string
		_ = s.db.NewSelect().
			Table("users").
			Column("line_user_id").
			Where("id = ?", prev.UserID).
			Where("line_opt_in = TRUE").
			Scan(ctx, &lineUserID)

		if strings.TrimSpace(lineUserID) == "" {
			logger.Infof("[booking:%s] skip push: empty line_user_id or not opted in", id.ID)
			return m, false, nil
		}

		switch m.Status {
		case enum.BookingApproved:
			if prev.Status != enum.BookingApproved {
				msg := fmt.Sprintf(
					"✅ อนุมัติการจองแล้ว\nเรื่อง: %s\nเวลา: %s–%s",
					helper.FirstNonEmpty(m.Title, prev.Title),
					helper.FormatTS(helper.FirstNonZero(m.StartTime, prev.StartTime)),
					helper.FormatTS(helper.FirstNonZero(m.EndTime, prev.EndTime)),
				)
				if err := linectrl.PushTextToLineUser(lineUserID, msg); err != nil {
					logger.Errf("[booking:%s] push approved error: %v", id.ID, err)
				} else {
					logger.Infof("[booking:%s] pushed approved to %s", id.ID, lineUserID)
				}
			} else {
				logger.Infof("[booking:%s] approved transition skipped (already approved)", id.ID)
			}

		case enum.BookingCanceled:
			msg := fmt.Sprintf(
				"❌ การจองถูกยกเลิก\nเรื่อง: %s\nเวลา: %s–%s",
				helper.FirstNonEmpty(m.Title, prev.Title),
				helper.FormatTS(helper.FirstNonZero(m.StartTime, prev.StartTime)),
				helper.FormatTS(helper.FirstNonZero(m.EndTime, prev.EndTime)),
			)
			if err := linectrl.PushTextToLineUser(lineUserID, msg); err != nil {
				logger.Errf("[booking:%s] push canceled error: %v", id.ID, err)
			} else {
				logger.Infof("[booking:%s] pushed canceled to %s", id.ID, lineUserID)
			}

		default:
			logger.Infof("[booking:%s] status=%s no push rule", id.ID, m.Status)
		}
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
		OrderExpr("b.created_at DESC")

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
		OrderExpr("b.created_at DESC").
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

func (s *BookingService) WarnExpiringBookings() error {
	now := time.Now().Unix()
	from := now + 15*60 // เป้าหมาย: เหลือ 15 นาที
	to := from + 90     // เผื่อหน้าต่าง ~90 วิ กันพลาดรอบ tick

	// ดึงรายการที่ยังไม่ได้เตือน
	type row struct {
		ID      string
		UserID  string
		Title   string
		EndTime int64
	}
	var rows []row

	err := s.db.NewSelect().
		TableExpr("bookings AS b").
		ColumnExpr("b.id, b.user_id, b.title, b.end_time").
		Where("b.status = ?", enum.BookingApproved).
		Where("b.deleted_at IS NULL").
		Where("b.expire_warn_sent_at IS NULL").
		Where("b.end_time BETWEEN ? AND ?", from, to).
		Scan(context.Background(), &rows)
	if err != nil {
		return err
	}
	if len(rows) == 0 {
		return nil
	}

	// เตรียม LINE bot client (reuse ตลอดรอบ)
	secret := os.Getenv("LINE_CHANNEL_SECRET")
	token := os.Getenv("LINE_CHANNEL_TOKEN")
	bot, _ := linebot.New(secret, token)

	for _, bk := range rows {
		// หา line_user_id ของเจ้าของ (เฉพาะคนที่ opt-in)
		var lineUserID string
		_ = s.db.NewSelect().
			Table("users").
			Column("line_user_id").
			Where("id = ?", bk.UserID).
			Where("line_opt_in = TRUE").
			Scan(context.Background(), &lineUserID)
		if lineUserID == "" || bot == nil {
			continue
		}

		msg := fmt.Sprintf(
			"⏰ ใกล้หมดเวลา (เหลือ ~15 นาที)\nเรื่อง: %s\nสิ้นสุด: %s",
			bk.Title,
			time.Unix(bk.EndTime, 0).Format("02/01 15:04"),
		)

		if _, err := bot.PushMessage(lineUserID, linebot.NewTextMessage(msg)).Do(); err == nil {
			// mark ส่งแล้ว (กันส่งซ้ำ)
			_, _ = s.db.NewUpdate().
				Table("bookings").
				Set("expire_warn_sent_at = ?", time.Now().Unix()).
				Where("id = ?", bk.ID).
				Exec(context.Background())
		}
	}

	return nil
}
