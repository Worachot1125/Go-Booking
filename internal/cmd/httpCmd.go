package cmd

import (
	"app/app/routes"
	"app/config"
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/spf13/cobra"

	bookingService "app/app/controller/booking"
	linectrl "app/app/controller/line"
)

func HttpCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "http",
		Short: "Run server on HTTP protocol",
		Run: func(cmd *cobra.Command, args []string) {

			config.Database()
			db := config.GetDB()

			if err := linectrl.InitLine(db); err != nil {
                fmt.Printf("LINE init error: %v\n", err)
            }
			bookingSvc := bookingService.NewBookingService(db)

			// === เปลี่ยนสถานะเป็น Finished หลังหมดเวลา (ทุก 5 ชม.) ===
			go func() {
				for {
					err := bookingSvc.AutoExpiredBookings()
					if err != nil {
						fmt.Printf("Auto-status error: %v\n", err)
					}
					time.Sleep(5 * time.Hour)
				}
			}()

			// === แจ้งเตือนเหลือ 15 นาที (ทุก 60 วิ) — แยก loop ไม่ทับของเดิม ===
			go func() {
				ticker := time.NewTicker(60 * time.Second)
				defer ticker.Stop()
				for range ticker.C {
					if err := bookingSvc.WarnExpiringBookings(); err != nil {
						fmt.Printf("Warn-15m error: %v\n", err)
					}
				}
			}()

			r := gin.New()

			// Logger
			r.Use(gin.Logger())
			r.Use(gin.Recovery())

			routes.Router(r)
			r.Run(":8080")
		},
	}
}
