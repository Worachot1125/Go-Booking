package line

import (
	"os"

	"github.com/line/line-bot-sdk-go/v7/linebot"
	"github.com/uptrace/bun"
)

// ===== constructors เดิม =====
func NewService(db *bun.DB) *Service {
	b, _ := linebot.New(os.Getenv("LINE_CHANNEL_SECRET"), os.Getenv("LINE_CHANNEL_TOKEN"))
	return &Service{db: db, bot: b}
}

func NewController(db *bun.DB) *Controller {
	return &Controller{Svc: NewService(db)}
}

// ====== SINGLETON + WRAPPERS (เพื่อรองรับโค้ดเก่า) ======
var defaultSvc *Service

// Export ชื่อเหมือนที่ที่อื่นเรียกใช้อยู่: linectrl.InitLine(...)
func InitLine(db *bun.DB) error {
	b, err := linebot.New(os.Getenv("LINE_CHANNEL_SECRET"), os.Getenv("LINE_CHANNEL_TOKEN"))
	if err != nil {
		return err
	}
	defaultSvc = &Service{db: db, bot: b}
	return nil
}

// ถ้าต้องเข้าถึง service ภายนอก (optional)
func Svc() *Service { return defaultSvc }

// Export ชื่อเหมือนที่ที่อื่นเรียกใช้อยู่: linectrl.PushTextToLineUser(...)
func PushTextToLineUser(lineUserID, text string) error {
	if defaultSvc == nil {
		return nil
	}
	return defaultSvc.PushTextToLineUser(lineUserID, text)
}
