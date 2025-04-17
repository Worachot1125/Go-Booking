package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// NotReqArgs Not required areuments
func NotReqArgs(cmd *cobra.Command, args []string) error {
	if len(args) != 0 {
		return fmt.Errorf("not required arguments")
	}
	return nil
}
