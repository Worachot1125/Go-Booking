package console

import (
	"github.com/spf13/cobra"

	"app/internal/cmd"
	"app/internal/logger"
)

func helloCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "hello",
		Args: cmd.NotReqArgs,
		Run: func(cmd *cobra.Command, args []string) {
			logger.Infof("Hello, world")
		},
	}
	return cmd
}