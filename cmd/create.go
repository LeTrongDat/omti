package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"

	"github.com/spf13/cobra"
)

// createCmd represents the create command
var createCmd = &cobra.Command{
	Use:   "create <repo_name> <folder_path>",
	Short: "Create a new GitHub repository and push local folder as the first commit",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		repoName := args[0]
		folderPath := args[1]

		logger := createCustomLogger()
		logger.Info("🚀 Starting the repository creation process")

		// Run all necessary environment checks before proceeding
		if err := checkEnvironment(); err != nil {
			logger.Fatalf("❌ Environment checks failed: %v", err)
		}

		// Check if the repository already exists
		exists, err := repoExists(repoName)
		if err != nil {
			logger.Fatalf("❌ Failed to check if repository exists: %v", err)
		}
		if exists {
			logger.Infof("✅ Repository %s already exists, skipping creation", repoName)
		} else {
			// Create new GitHub repository if it doesn't exist
			logger.Infof("Creating new GitHub repository: %s", repoName)
			if err := createRepo(repoName); err != nil {
				logger.Fatalf("❌ Failed to create GitHub repository: %v", err)
			}
			logger.Infof("✅ Repository %s created successfully", repoName)
		}

		// Check if the repository is empty (has no commits) before pushing
		isEmpty, err := repoIsEmpty(repoName)
		if err != nil {
			logger.Fatalf("❌ Failed to check if repository is empty: %v", err)
		}
		if isEmpty {
			// Push folder to GitHub as the first commit if the repo is empty
			logger.Infof("Pushing %s to GitHub", folderPath)
			if err := pushToRepo(repoName, folderPath); err != nil {
				logger.Fatalf("❌ Failed to push folder to repository: %v", err)
			}
			logger.Info("✅ Folder pushed successfully")
		} else {
			logger.Infof("✅ Repository %s already has commits, skipping push", repoName)
		}

		logger.Info("🎉 Repository creation process completed successfully!")
	},
}

func init() {
	repoCmd.AddCommand(createCmd)
}

// repoExists checks if the GitHub repository already exists
func repoExists(repoName string) (bool, error) {
	cmd := exec.Command("gh", "repo", "view", repoName)
	if err := cmd.Run(); err != nil {
		// If the command fails with an exit code, the repository likely doesn't exist
		if exitErr, ok := err.(*exec.ExitError); ok && exitErr.ExitCode() != 0 {
			return false, nil
		}
		return false, err
	}
	return true, nil // Repo exists if the command succeeded
}

// repoIsEmpty checks if the GitHub repository has no commits (is empty)
func repoIsEmpty(repoName string) (bool, error) {
	// Run the `gh` command to retrieve repository information
	cmd := exec.Command("gh", "repo", "view", repoName, "--json", "defaultBranchRef")
	output, err := cmd.Output()
	if err != nil {
		return false, fmt.Errorf("failed to retrieve repository info: %w", err)
	}

	// Define a local struct to match the JSON output structure
	var result struct {
		DefaultBranchRef struct {
			Name string `json:"name"`
		} `json:"defaultBranchRef"`
	}

	// Unmarshal the JSON output into the local struct
	if err := json.Unmarshal(output, &result); err != nil {
		return false, fmt.Errorf("failed to parse JSON output: %w", err)
	}

	// Check if the name of the default branch is empty, indicating no commits
	return result.DefaultBranchRef.Name == "", nil
}

// createRepo creates a new repository on GitHub using the GitHub CLI
func createRepo(repoName string) error {
	return runCommand("gh", "repo", "create", repoName, "--public", "--source=.", "--remote=origin")
}

// pushToRepo initializes and pushes a local folder to the GitHub repository
func pushToRepo(repoName, folderPath string) error {
	if err := os.Chdir(folderPath); err != nil {
		return fmt.Errorf("failed to change directory: %w", err)
	}

	if err := runCommand("git", "init"); err != nil {
		return err
	}
	if err := runCommand("git", "add", "."); err != nil {
		return err
	}
	if err := runCommand("git", "commit", "-m", "Initial commit"); err != nil {
		return err
	}
	if err := runCommand("git", "branch", "-M", "main"); err != nil {
		return err
	}
	return runCommand("git", "push", "-u", "origin", "main")
}
