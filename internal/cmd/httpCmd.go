package cmd

import (
	"app/app/routes"

	"github.com/gin-gonic/gin"
	"github.com/spf13/cobra"
)

func HttpCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "http",
		Short: "Run server on HTTP protocol",
		Run: func(cmd *cobra.Command, args []string) {
			r := gin.New()
			
			// Logger
			r.Use(gin.Logger())
			r.Use(gin.Recovery())

			// Router ที่กำหนด CORS
			routes.Router(r)

			r.Run(":8080")
		},
	}
}

