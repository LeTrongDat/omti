package cmd

import (
	"github.com/spf13/cobra"
)

var dbCmd = &cobra.Command{
	Use:   "db",
	Short: "Database-related operations, including backup and restore",
	Long:  `The "db" command contains subcommands for managing and interacting with databases.`,
}

func init() {
	rootCmd.AddCommand(dbCmd)
}
