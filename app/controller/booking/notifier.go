package booking

import (
	linectrl "app/app/controller/line"
	"app/app/enum"
	"context"
	"fmt"
	"time"
)

func (s *Service) StartExpiryWarningWorker(ctx context.Context) {
	ticker := time.NewTicker(60 * time.Second)
	go func() {
		for {
			select {
			case <-ctx.Done():
				ticker.Stop()
				return
			case <-ticker.C:
				_ = s.warnExpiringBookings(ctx)
			}
		}
	}()
}

func (s *Service) warnExpiringBookings(ctx context.Context) error {
	now := time.Now().Unix()
	from := now + 15*60 // อีก 15 นาที
	to := from + 90     // กันรอบ tick

	type row struct {
		ID           string
		UserID       string
		Title        string
		EndTime      int64
		LineUserID   string
		RoomName     string
		BuildingName string
	}

	var rows []row
	err := s.db.NewSelect().
		TableExpr("bookings AS b").
		Join("JOIN users  AS u ON u.id = b.user_id AND u.deleted_at IS NULL AND u.line_opt_in = TRUE").
		Join("JOIN rooms  AS r ON r.id = b.room_id AND r.deleted_at IS NULL").
		Join("JOIN buildings AS bd ON bd.id = r.building_id AND bd.deleted_at IS NULL").
		ColumnExpr("b.id, b.user_id, b.title, b.end_time, u.line_user_id, r.name AS room_name").
		Where("b.status = ?", enum.BookingApproved).
		Where("b.deleted_at IS NULL").
		Where("b.expire_warn_sent_at IS NULL").
		Where("b.end_time BETWEEN ? AND ?", from, to).
		Scan(ctx, &rows)
	if err != nil {
		return err
	}
	if len(rows) == 0 {
		return nil
	}

	// เวลาให้เป็นโซนไทย
	loc, _ := time.LoadLocation("Asia/Bangkok")

	for _, bk := range rows {
		if bk.LineUserID == "" {
			continue
		}

		msg := fmt.Sprintf(
			"⏰ ใกล้หมดเวลา (เหลือ ~15 นาที)\nเรื่อง: %s\nห้อง: %s\nสิ้นสุด: %s",
			bk.Title,
			bk.RoomName,
			time.Unix(bk.EndTime, 0).In(loc).Format("02/01 15:04"),
		)

		if err := linectrl.PushTextToLineUser(bk.LineUserID, msg); err == nil {
			_, _ = s.db.NewUpdate().
				Table("bookings").
				Set("expire_warn_sent_at = ?", time.Now().Unix()).
				Where("id = ?", bk.ID).
				Exec(ctx)
		}
	}
	return nil
}
