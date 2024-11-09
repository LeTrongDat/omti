package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"

	"github.com/sirupsen/logrus"
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

		logger := logrus.New()
		logger.SetFormatter(&logrus.TextFormatter{
			DisableTimestamp: true,
		})

		logger.Info("Starting the repository creation process")

		// Step 1: Check if GitHub CLI (gh) is installed
		logger.Info("Checking if GitHub CLI is installed")
		if !isCommandAvailable("gh") {
			logger.Warn("GitHub CLI not found. Installing it now...")
			if err := installGH(); err != nil {
				logger.Fatalf("Failed to install GitHub CLI: %v", err)
			}
		} else {
			logger.Info("GitHub CLI is already installed")
		}

		// Step 2: Authenticate with GitHub if needed
		logger.Info("Checking GitHub authentication")
		if !isAuthenticated() {
			logger.Warn("GitHub authentication not detected. Logging in...")
			if err := authenticateGH(); err != nil {
				logger.Fatalf("Failed to authenticate with GitHub: %v", err)
			}
		} else {
			logger.Info("GitHub authentication is already set up")
		}

		// Step 3: Create new GitHub repository
		logger.Infof("Creating new GitHub repository: %s", repoName)
		if err := createRepo(repoName); err != nil {
			logger.Fatalf("Failed to create GitHub repository: %v", err)
		}
		logger.Infof("Repository %s created successfully", repoName)

		// Step 4: Push folder to GitHub as the first commit
		logger.Infof("Pushing %s to GitHub", folderPath)
		if err := pushToRepo(repoName, folderPath); err != nil {
			logger.Fatalf("Failed to push folder to repository: %v", err)
		}
		logger.Info("Folder pushed successfully")

		logger.Info("Repository creation process completed successfully!")
	},
}

func init() {
	repoCmd.AddCommand(createCmd)
}

// Check if a command is available on the system
func isCommandAvailable(name string) bool {
	_, err := exec.LookPath(name)
	return err == nil
}

// Install GitHub CLI based on the operating system
func installGH() error {
	switch runtime.GOOS {
	case "darwin":
		return runCommand("brew", "install", "gh")
	case "linux":
		if err := runCommand("sudo", "apt", "update"); err != nil {
			return err
		}
		return runCommand("sudo", "apt", "install", "-y", "gh")
	default:
		return fmt.Errorf("unsupported OS: %s", runtime.GOOS)
	}
}

// Check if the user is authenticated with GitHub
func isAuthenticated() bool {
	output, err := exec.Command("gh", "auth", "status").CombinedOutput()
	return err == nil && strings.Contains(string(output), "Logged in")
}

// Authenticate with GitHub
func authenticateGH() error {
	return runCommand("gh", "auth", "login")
}

// Create a new repository with the given name on GitHub
func createRepo(repoName string) error {
	return runCommand("gh", "repo", "create", repoName, "--public", "--source=.", "--remote=origin")
}

// Push the specified folder to GitHub as the first commit
func pushToRepo(repoName, folderPath string) error {
	os.Chdir(folderPath)
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

// Helper function to run a shell command
func runCommand(name string, args ...string) error {
	cmd := exec.Command(name, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
