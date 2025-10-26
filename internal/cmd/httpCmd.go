package cmd

import (
	"app/app/routes"
	"app/config"
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
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
		RunE: func(cmd *cobra.Command, args []string) error {

			// ====== ENV & PORT ======
			port := os.Getenv("PORT")
			if port == "" {
				port = "8080"
			}
			addr := "0.0.0.0:" + port

			// ====== DB & Services ======
			config.Database()
			db := config.GetDB()

			if err := linectrl.InitLine(db); err != nil {
				log.Printf("LINE init error: %v\n", err)
			}
			bookingSvc := bookingService.NewBookingService(db)

			// ====== Background jobs ======
			// 1) เปลี่ยนสถานะเป็น Finished หลังหมดเวลา (ทุก 5 ชม.)
			go func() {
				for {
					if err := bookingSvc.AutoExpiredBookings(); err != nil {
						log.Printf("Auto-status error: %v\n", err)
					}
					time.Sleep(5 * time.Hour)
				}
			}()

			// 2) แจ้งเตือนเหลือ 15 นาที
			//    NOTE: ถ้าอยาก "ทุก 60 วิ" ให้ปรับเป็น time.NewTicker(1 * time.Minute)
			warnInterval := 5 * time.Minute // เปลี่ยนเป็น 1 * time.Minute ถ้าต้องการทุก 60 วินาที
			go func() {
				ticker := time.NewTicker(warnInterval)
				defer ticker.Stop()
				for range ticker.C {
					if err := bookingSvc.WarnExpiringBookings(); err != nil {
						log.Printf("Warn-15m error: %v\n", err)
					}
				}
			}()

			// ====== HTTP Router (Gin) ======
			if os.Getenv("GIN_MODE") == "" {
				gin.SetMode(gin.ReleaseMode) // default ให้เบา log ใน prod; ลบออกได้หากอยากเห็น log
			}
			r := gin.New()
			r.Use(gin.Logger(), gin.Recovery())

			// หน้าแรกแบบ Node.js: “API is running…”
			r.GET("/", func(c *gin.Context) {
				c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(`<h1>API is running...</h1>`))
			})

			// Health endpoint สำหรับเช็ค readiness
			r.GET("/health", func(c *gin.Context) {
				c.JSON(http.StatusOK, gin.H{"status": "ok", "time": time.Now().UTC()})
			})

			// เส้นทางหลักของระบบ (เดิมของคุณ)
			routes.Router(r)

			// ====== Graceful server ======
			srv := &http.Server{
				Addr:              addr,
				Handler:           r,
				ReadHeaderTimeout: 10 * time.Second,
			}

			// start server
			go func() {
				log.Printf("🚀 Server running on http://localhost:%s\n", port)
				if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
					log.Fatalf("server error: %v", err)
				}
			}()

			// รอสัญญาณปิด (CTRL+C / SIGTERM)
			quit := make(chan os.Signal, 1)
			signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
			<-quit
			log.Println("🛑 Shutting down...")

			ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
			defer cancel()
			if err := srv.Shutdown(ctx); err != nil {
				return fmt.Errorf("server forced to shutdown: %w", err)
			}
			log.Println("✅ Server exited gracefully")
			return nil
		},
	}
}
