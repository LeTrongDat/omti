package cmd

import (
	"github.com/spf13/cobra"
)

// repoCmd represents the repo command
var repoCmd = &cobra.Command{
	Use:   "repo",
	Short: "Manage GitHub repositories",
	Long: `The "repo" command provides useful tools to quickly spin up and manage
your GitHub repositories. It allows you to create new repositories, push
local projects to GitHub, and perform other common GitHub operations.`,
}

func init() {
	rootCmd.AddCommand(repoCmd)
}
