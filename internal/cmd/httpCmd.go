package cmd

import (
	"app/app/routes"
	"app/config"
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/spf13/cobra"

	bookingService "app/app/controller/booking"
)

func HttpCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "http",
		Short: "Run server on HTTP protocol",
		Run: func(cmd *cobra.Command, args []string) {

			config.Database()
			db := config.GetDB()
			bookingSvc := bookingService.NewBookingService(db)

			go func() {
				for {
					err := bookingSvc.AutoDeleteExpiredBookings()
					if err != nil {
						fmt.Printf("Auto-delete error: %v\n", err)
					}
					time.Sleep(5 * time.Minute)
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
