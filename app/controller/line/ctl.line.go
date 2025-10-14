package line

import (
	"strings"

	"app/app/response"
	"app/internal/logger"

	"github.com/gin-gonic/gin"
	"github.com/line/line-bot-sdk-go/v7/linebot"
)

type Controller struct {
	Svc *Service // <- ใช้ S ตัวใหญ่ และทุกที่เรียก c.Svc
}

func (c *Controller) IssuePairingCode(ctx *gin.Context) {
	uid, err := c.Svc.UserIDFromContext(ctx)
	if err != nil {
		response.Unauthorized(ctx, err.Error())
		return
	}
	code, exp, err := c.Svc.GeneratePairingCode(ctx, uid)
	if err != nil {
		response.InternalError(ctx, "cannot generate pairing code")
		return
	}
	response.Success(ctx, gin.H{"code": "PAIR-" + code, "expires_at": exp})
}

func (c *Controller) Webhook(ctx *gin.Context) {
	events, err := c.Svc.bot.ParseRequest(ctx.Request)
	if err != nil {
		logger.Errf("line webhook parse: %v", err)
		ctx.Status(400)
		return
	}

	clean := func(s string) string {
		s = strings.ReplaceAll(s, "\u200B", "")
		s = strings.TrimSpace(s)
		return s
	}

	for _, ev := range events {
		switch ev.Type {
		case linebot.EventTypeFollow:
			_, _ = c.Svc.bot.ReplyMessage(ev.ReplyToken,
				linebot.NewTextMessage("ขอบคุณที่เพิ่มบอท! พิมพ์ PAIR-<โค้ด> เพื่อเชื่อมบัญชี"),
			).Do()

		case linebot.EventTypeMessage:
			msg, ok := ev.Message.(*linebot.TextMessage)
			if !ok {
				break
			}
			txt := clean(msg.Text)
			upper := strings.ToUpper(txt)

			// ต้องมี UserID (เฉพาะแชต 1:1) — กันเคสจาก group/room
			if ev.Source == nil || ev.Source.UserID == "" {
				_, _ = c.Svc.bot.ReplyMessage(ev.ReplyToken,
					linebot.NewTextMessage("คำสั่งนี้ใช้ได้เฉพาะห้องแชตกับบอทแบบ 1:1"),
				).Do()
				break
			}

			// แยก code
			var code string
			switch {
			case strings.HasPrefix(upper, "PAIR-"):
				code = clean(txt[5:])
			case strings.HasPrefix(upper, "PAIR "):
				code = clean(txt[5:])
			default:
				// ไม่ใช่คำสั่ง pairing → echo แล้วจบ
				_, _ = c.Svc.bot.ReplyMessage(ev.ReplyToken, linebot.NewTextMessage("echo: "+msg.Text)).Do()
				continue // <<< สำคัญ: ต้อง continue ไม่งั้นจะ reply ซ้ำ
			}

			if code == "" {
				_, _ = c.Svc.bot.ReplyMessage(ev.ReplyToken,
					linebot.NewTextMessage("รูปแบบโค้ดไม่ถูกต้อง ลองพิมพ์เป็น PAIR-<โค้ด>"),
				).Do()
				break
			}

			// ใช้เมธอดใหม่ที่ทำทุกอย่างใน transaction เดียว
			_, err := c.Svc.resolvePairing(ctx, code, ev.Source.UserID)
			if err != nil {
				msg := "เกิดข้อผิดพลาด กรุณาลองใหม่"
				if err.Error() == "invalid_or_expired_code" {
					msg = "โค้ดไม่ถูกต้องหรือหมดอายุ"
				} else if err.Error() == "line_user_id_already_linked" {
					msg = "LINE นี้ถูกเชื่อมกับบัญชีอื่นแล้ว"
				}
				_, _ = c.Svc.bot.ReplyMessage(ev.ReplyToken, linebot.NewTextMessage(msg)).Do()
				break
			}

			_, _ = c.Svc.bot.ReplyMessage(ev.ReplyToken,
				linebot.NewTextMessage("เชื่อมบัญชีเรียบร้อย ✅"),
			).Do()
		}
	}
	ctx.Status(200)
}
