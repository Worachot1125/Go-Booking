package cmd

import (
	provider "app/app/provider/database"
	"app/config"
	"app/internal/logger"
	"os"

	"github.com/spf13/cobra"
)

// Migrate Command
func Migrate() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "migrate",
		Args: NotReqArgs,
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			return provider.Open(cmd.Context())
		},
		PersistentPostRunE: func(cmd *cobra.Command, args []string) error {
			return provider.Close(cmd.Context())
		},
		Run: func(cmd *cobra.Command, args []string) {
			migrateUp().Run(cmd, args)
		},
	}
	cmd.AddCommand(migrateUp())
	cmd.AddCommand(migrateDown())
	cmd.AddCommand(migrateSeed())
	cmd.AddCommand(migrateRefresh())
	return cmd
}

func migrateUp() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "up",
		Args: NotReqArgs,
		Run: func(cmd *cobra.Command, args []string) {
			db := config.GetDB()
			if err := modelUp(db); err != nil {
				logger.Errf("%s", err)
				os.Exit(1)
			}
		},
	}
	return cmd
}

func migrateDown() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "down",
		Args: NotReqArgs,
		Run: func(cmd *cobra.Command, args []string) {
			db := config.GetDB()
			if err := modelDown(db); err != nil {
				logger.Errf("%s", err)
				os.Exit(1)
			}
		},
	}
	return cmd
}

func migrateRefresh() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "refresh",
		Args: NotReqArgs,
		Run: func(cmd *cobra.Command, args []string) {
			db := config.GetDB()
			if err := modelDown(db); err != nil {
				logger.Errf("%s", err)
				os.Exit(1)
			}
			if err := modelUp(db); err != nil {
				logger.Errf("%s", err)
				os.Exit(1)
			}
		},
	}
	return cmd
}

func migrateSeed() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "seed",
		Args: NotReqArgs,
		Run: func(cmd *cobra.Command, args []string) {
			db := config.GetDB()
			if err := modelSeed(db); err != nil {
				logger.Errf("%s", err)
				os.Exit(1)
			}
		},
	}
	return cmd
}
