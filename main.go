package main

import (
	"app/app/console"
	"app/config"
	"app/internal/cmd"
	"log"

	"github.com/spf13/cobra"
)

func main() {

	config.Init() // Initialize configuration

	// Start the HTTP server
	if err := command(); err != nil {
		log.Fatalf("Error running command: %v", err)
	}
}

// Command setup for running the HTTP server
func command() error {
	cmda := &cobra.Command{
		Use:  "app",
		Args: cobra.NoArgs,
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Usage()
		},
	}

	cmds := &cobra.Command{
		Use:   "cmd",
		Short: "List all commands",
	}

	cmda.AddCommand(cmds)
	cmds.AddCommand(console.Commands()...)
	cmda.AddCommand(cmd.HttpCmd())
	cmda.AddCommand(cmd.Migrate())

	return cmda.Execute()
}
