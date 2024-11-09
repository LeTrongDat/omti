package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"

	"github.com/sirupsen/logrus"
)

// checkEnvironment performs all necessary checks for `gh`, `git`, and SSH key setup
func checkEnvironment() error {
	logger := createCustomLogger()

	logger.Info("🔍 Starting environment checks")

	// Check for GitHub CLI
	if err := checkGH(); err != nil {
		logger.Fatalf("❌ GitHub CLI setup failed: %v", err)
	}

	// Check for Git
	if err := checkGit(); err != nil {
		logger.Fatalf("❌ Git setup failed: %v", err)
	}

	// Ensure SSH key is set up
	if err := ensureSSHKey(); err != nil {
		logger.Fatalf("❌ SSH key setup failed: %v", err)
	}

	logger.Info("✅ All environment checks passed")
	return nil
}

// checkGH ensures the GitHub CLI is installed, installing it if necessary
func checkGH() error {
	if !isCommandAvailable("gh") {
		fmt.Println("Installing GitHub CLI...")
		if err := installGH(); err != nil {
			return fmt.Errorf("failed to install GitHub CLI: %w", err)
		}
		fmt.Println("✅ GitHub CLI installed")
	} else {
		fmt.Println("✅ GitHub CLI is already installed")
	}
	return nil
}

// checkGit ensures Git is installed, installing it if necessary
func checkGit() error {
	if !isCommandAvailable("git") {
		fmt.Println("Installing Git...")
		if err := installGit(); err != nil {
			return fmt.Errorf("failed to install Git: %w", err)
		}
		fmt.Println("✅ Git installed")
	} else {
		fmt.Println("✅ Git is already installed")
	}
	return nil
}

// ensureSSHKey generates an SSH key if it doesn't exist and adds it to GitHub
func ensureSSHKey() error {
	sshKeyPath := filepath.Join(os.Getenv("HOME"), ".ssh", "id_rsa.pub")
	if _, err := os.Stat(sshKeyPath); os.IsNotExist(err) {
		fmt.Println("Generating SSH key...")
		if err := generateSSHKey(); err != nil {
			return fmt.Errorf("failed to generate SSH key: %w", err)
		}
		fmt.Println("✅ SSH key generated")
	}
	return addSSHKeyToGitHub(sshKeyPath)
}

// generateSSHKey creates a new SSH key pair
func generateSSHKey() error {
	sshDir := filepath.Join(os.Getenv("HOME"), ".ssh")
	os.MkdirAll(sshDir, 0700)

	cmd := exec.Command("ssh-keygen", "-t", "rsa", "-b", "4096", "-f", filepath.Join(sshDir, "id_rsa"), "-N", "")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// addSSHKeyToGitHub adds the SSH key to GitHub using the GitHub CLI
func addSSHKeyToGitHub(sshKeyPath string) error {
	_, err := os.ReadFile(sshKeyPath)
	if err != nil {
		return fmt.Errorf("failed to read SSH key: %w", err)
	}

	fmt.Println("Adding SSH key to GitHub...")
	if err := runCommand("gh", "ssh-key", "add", sshKeyPath, "--title", "omti-cli-key"); err != nil {
		return fmt.Errorf("failed to add SSH key to GitHub: %w", err)
	}
	fmt.Println("✅ SSH key added to GitHub")
	return nil
}

// installGH installs GitHub CLI based on the operating system
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

// installGit installs Git based on the operating system
func installGit() error {
	switch runtime.GOOS {
	case "darwin":
		return runCommand("brew", "install", "git")
	case "linux":
		if err := runCommand("sudo", "apt", "update"); err != nil {
			return err
		}
		return runCommand("sudo", "apt", "install", "-y", "git")
	default:
		return fmt.Errorf("unsupported OS: %s", runtime.GOOS)
	}
}

// isCommandAvailable checks if a command is available on the system
func isCommandAvailable(name string) bool {
	_, err := exec.LookPath(name)
	return err == nil
}

// createCustomLogger creates a custom logger with symbols for user-friendly output
func createCustomLogger() *logrus.Logger {
	logger := logrus.New()
	logger.SetFormatter(&logrus.TextFormatter{
		DisableTimestamp: true,
		ForceColors:      true,
	})
	return logger
}

// runCommand is a helper function to run shell commands
func runCommand(name string, args ...string) error {
	cmd := exec.Command(name, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
