package helper

import "time"

// โหลดโซนเวลา BKK ครั้งเดียว
var bkk *time.Location

func init() {
	loc, err := time.LoadLocation("Asia/Bangkok")
	if err != nil {
		bkk = time.Local
	} else {
		bkk = loc
	}
}

// Export function
func FormatTS(ts int64) string {
	if ts <= 0 {
		return "-"
	}
	return time.Unix(ts, 0).In(bkk).Format("02/01 15:04")
}