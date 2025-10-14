package line

import (
	"context"
	"crypto/rand"
	"database/sql"
	"encoding/hex"
	"errors"
	"strings"
	"time"

	"app/app/model"

	"github.com/gin-gonic/gin"
	"github.com/line/line-bot-sdk-go/v7/linebot"
	"github.com/uptrace/bun"
)

type Service struct {
	db  *bun.DB
	bot *linebot.Client
}

// ====== helpers ======
func randomCode(n int) (string, error) {
	b := make([]byte, n)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return strings.ToUpper(hex.EncodeToString(b)), nil
}

func (s *Service) UserIDFromContext(c *gin.Context) (string, error) {
	v, ok := c.Get("claims")
	if !ok {
		return "", errors.New("no claims in context")
	}
	m, ok := v.(map[string]any)
	if !ok {
		return "", errors.New("invalid claims")
	}
	id, ok := m["user_id"].(string) // ระบบของคุณใช้ string UUID
	if !ok || id == "" {
		return "", errors.New("user_id must be string")
	}
	return id, nil
}

// ====== core ======
func (s *Service) GeneratePairingCode(ctx context.Context, userID string) (string, int64, error) {
	code, err := randomCode(3) // 6 ตัว hex
	if err != nil {
		return "", 0, err
	}
	exp := time.Now().Unix() + 10*60

	rec := &model.LinePairingCode{
		UserID:    userID,
		Code:      code,
		ExpiresAt: exp,
	}
	rec.SetCreatedNow()
	if _, err := s.db.NewInsert().Model(rec).Exec(ctx); err != nil {
		return "", 0, err
	}
	return code, exp, nil
}

func (s *Service) resolvePairing(ctx context.Context, rawCode, lineUserID string) (string, error) {
	clean := func(s string) string {
		s = strings.ReplaceAll(s, "\u200B", "")
		return strings.TrimSpace(s)
	}
	code := strings.ToUpper(clean(rawCode))
	now := time.Now().Unix()

	var userID string
	err := s.db.RunInTx(ctx, nil, func(ctx context.Context, tx bun.Tx) error {
		// 1) หาโค้ด + ล็อกแถว
		var row model.LinePairingCode
		if err := tx.NewSelect().
			Model(&row).
			Where("code = ?", code).
			Where("used_at IS NULL").
			Where("expires_at > ?", now).
			For("UPDATE").
			Scan(ctx); err != nil {
			return errors.New("invalid_or_expired_code")
		}
		userID = row.UserID

		// 2) กัน LINE นี้ไปผูกกับ user อื่นอยู่แล้ว
		var exists int
		if err := tx.NewSelect().
			Table("users").
			Where("line_user_id = ? AND id <> ?", lineUserID, userID).
			ColumnExpr("1").
			Limit(1).
			Scan(ctx, &exists); err != nil && err != sql.ErrNoRows {
			return err
		}
		if exists == 1 {
			return errors.New("line_user_id_already_linked")
		}

		// 3) อัปเดต users
		if _, err := tx.NewUpdate().
			Table("users").
			Set("line_user_id = ?", lineUserID).
			Set("line_opt_in = TRUE").
			Set("line_linked_at = ?", now).
			Where("id = ?", userID).
			Where("deleted_at IS NULL").
			Exec(ctx); err != nil {
			return err
		}

		// 4) ปิดโค้ด
		if _, err := tx.NewUpdate().
			Table("line_pairing_codes").
			Set("used_at = ?", now).
			Set("updated_at = ?", now).
			Where("id = ?", row.ID).
			Exec(ctx); err != nil {
			return err
		}
		return nil
	})
	return userID, err
}

func (s *Service) PushTextToLineUser(lineUserID, text string) error {
	if s.bot == nil || lineUserID == "" || strings.TrimSpace(text) == "" {
		return nil
	}
	_, err := s.bot.PushMessage(lineUserID, linebot.NewTextMessage(text)).Do()
	return err
}
