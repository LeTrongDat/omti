package cmd

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	"github.com/spf13/cobra"
)

var sshKeyPath = filepath.Join(os.Getenv("HOME"), ".ssh", "id_rsa.pub")

// backupCmd represents the command to back up a PostgreSQL database
var backupCmd = &cobra.Command{
	Use:   "backup <remote_user> <remote_host> <remote_port> <remote_db_port> <db_name> <db_user> <db_password> <local_save_path>",
	Short: "Backup a PostgreSQL database on a remote server and save it locally",
	Args:  cobra.ExactArgs(8),
	Run: func(cmd *cobra.Command, args []string) {
		remoteUser := args[0]
		remoteHost := args[1]
		remotePort := args[2]
		remoteDBPort := args[3]
		dbName := args[4]
		dbUser := args[5]
		dbPassword := args[6]
		localSavePath := args[7]

		logger := createCustomLogger()
		logger.Info("🚀 Starting database backup process")

		// Ensure pg_dump is installed locally, or install it
		if err := ensurePgDump(); err != nil {
			logger.Fatalf("❌ pg_dump is required but not installed: %v", err)
		}

		// Set up SSH tunnel and perform the backup
		if err := backupDatabase(remoteUser, remoteHost, remotePort, remoteDBPort, dbName, dbUser, dbPassword, localSavePath); err != nil {
			logger.Fatalf("❌ Database backup failed: %v", err)
		}

		logger.Info("✅ Database backup completed successfully")
	},
}

func init() {
	dbCmd.AddCommand(backupCmd)
}

// backupDatabase performs the database backup over an SSH tunnel and saves it locally
func backupDatabase(remoteUser, remoteHost, remotePort, remoteDBPort, dbName, dbUser, dbPassword, localSavePath string) error {
	timestamp := time.Now().Format("20060102_150405")
	backupFile := filepath.Join(localSavePath, fmt.Sprintf("%s_backup_%s.sql", dbName, timestamp))

	// Define SSH tunnel command assuming SSH key-based authentication
	tunnelCmd := exec.Command("ssh", "-fN", "-L", fmt.Sprintf("5433:localhost:%s", remoteDBPort), "-p", remotePort, fmt.Sprintf("%s@%s", remoteUser, remoteHost))
	tunnelCmd.Stdout = os.Stdout
	tunnelCmd.Stderr = os.Stderr

	// Start SSH tunnel and ensure the process is closed afterward
	if err := tunnelCmd.Run(); err != nil {
		return fmt.Errorf("failed to start SSH tunnel: %w", err)
	}
	// Ensure the SSH tunnel process on port 5433 is killed after backup completes
	defer func() {
		if err := killProcessOnPort("5433"); err != nil {
			fmt.Errorf("❌ Failed to kill SSH tunnel process on port 5433: %v", err)
		} else {
			fmt.Println("✅ SSH tunnel process on port 5433 terminated successfully.")
		}
	}()
	// Define pg_dump command to back up the database through the SSH tunnel
	pgDumpCmd := exec.Command("pg_dump",
		"-h", "localhost",
		"-p", "5433",
		"-U", dbUser,
		"-d", dbName,
		"-F", "c",
		"-f", backupFile,
	)

	// Set the PGPASSWORD environment variable to provide the password for pg_dump
	pgDumpCmd.Env = append(os.Environ(), fmt.Sprintf("PGPASSWORD=%s", dbPassword))
	var stdOut, stdErr bytes.Buffer
	pgDumpCmd.Stdout = &stdOut
	pgDumpCmd.Stderr = &stdErr

	// Run pg_dump to create the backup file and check for errors
	if err := pgDumpCmd.Run(); err != nil {
		return fmt.Errorf("failed to execute pg_dump: %w\nOutput: %s\nError: %s", err, stdOut.String(), stdErr.String())
	}

	fmt.Printf("✅ Backup saved to %s\n", backupFile)
	return nil
}
