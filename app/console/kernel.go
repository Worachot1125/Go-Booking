package console

import "github.com/spf13/cobra"

// Commands Foe AddCommand
func Commands() []*cobra.Command {
	return []*cobra.Command{
		helloCmd(),
	}
}
