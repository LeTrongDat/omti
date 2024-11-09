package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/spf13/cobra"
)

// tagCmd represents the command to tag the latest commit of a specified branch
var tagCmd = &cobra.Command{
	Use:   "tag <tag_name> <branch_name> <folder_path>",
	Short: "Create a new tag for the latest commit of a branch in the repository",
	Args:  cobra.ExactArgs(3),
	Run: func(cmd *cobra.Command, args []string) {
		tagName := args[0]
		branchName := args[1]
		folderPath := args[2]

		logger := createCustomLogger()
		logger.Info("🚀 Starting the tagging process")

		// Run environment checks before proceeding
		if err := checkEnvironment(); err != nil {
			logger.Fatalf("❌ Environment checks failed: %v", err)
		}

		// Change directory to the repository folder path
		if err := os.Chdir(folderPath); err != nil {
			logger.Fatalf("❌ Failed to change directory to %s: %v", folderPath, err)
		}

		// Check if the tag already exists
		exists, err := tagExists(tagName)
		if err != nil {
			logger.Fatalf("❌ Failed to check if tag exists: %v", err)
		}
		if exists {
			logger.Infof("✅ Tag '%s' already exists, skipping creation", tagName)
			return
		}

		// Create the new tag
		if err := createTag(tagName, branchName); err != nil {
			logger.Fatalf("❌ Failed to create tag '%s': %v", tagName, err)
		}
		logger.Infof("✅ Tag '%s' created successfully", tagName)
	},
}

func init() {
	repoCmd.AddCommand(tagCmd)
}

// tagExists checks if a tag already exists in the repository
func tagExists(tagName string) (bool, error) {
	cmd := exec.Command("git", "tag", "--list", tagName)
	output, err := cmd.Output()
	if err != nil {
		return false, fmt.Errorf("failed to list tags: %w", err)
	}
	// If the output contains the tag name, it means the tag exists
	return strings.TrimSpace(string(output)) == tagName, nil
}

// createTag creates a new tag on the latest commit of the specified branch
func createTag(tagName, branchName string) error {
	// Fetch the latest updates from the branch
	if err := runCommand("git", "fetch", "origin", branchName); err != nil {
		return fmt.Errorf("failed to fetch latest changes from branch %s: %w", branchName, err)
	}

	// Create the tag pointing to the latest commit of the branch
	if err := runCommand("git", "tag", tagName, fmt.Sprintf("origin/%s", branchName)); err != nil {
		return fmt.Errorf("failed to create tag: %w", err)
	}

	// Push the tag to the remote repository
	if err := runCommand("git", "push", "origin", tagName); err != nil {
		return fmt.Errorf("failed to push tag to remote: %w", err)
	}
	return nil
}
