package cmd

import (
	"app/app/routes"

	"github.com/gin-gonic/gin"
	"github.com/spf13/cobra"
)

// HTTP is serve http ot https
func HttpCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "http",
		Short: "Run server on HTTP protocol",
		Run: func(cmd *cobra.Command, args []string) {
			r := gin.Default()
			routes.Router(r)
			r.Run(":8080") // Start server on port 8080
		},
	}
}
