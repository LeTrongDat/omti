package cmd

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/spf13/cobra"
)

// backupCmd represents the command to back up a PostgreSQL database
var backupCmd = &cobra.Command{
	Use: "backup <db_config> <local_save_path>",
	Short: `Backup a PostgreSQL database, locally or over SSH.
	
db_config: <username>:<password>@<host>:<port>/<dbname>
e.g., postgres:v8hlDV0yMAHHlIurYupj@10.1.0.54:15432/golang

--remote: <user>@<host>:<remote-db-port>
e.g., --remote admin@192.168.1.10:5432`,
	Args: cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		dbConfig := args[0]
		localSavePath := args[1]

		logger := createCustomLogger()
		logger.Info("🚀 Starting database backup process")

		var err error
		var dbUser, dbPassword, dbHost, dbPort, dbName string

		if remoteFlag != "" {
			dbUser, dbPassword, dbHost, dbPort, dbName, err = parseDBConfig(dbConfig)
			if err != nil {
				logger.Fatalf("❌ Invalid database configuration format: %v", err)
			}
			remoteUser, remoteHost, remoteDBPort, err := parseRemoteFlag(remoteFlag)
			if err != nil {
				logger.Fatalf("❌ Invalid --remote format: %v", err)
			}
			err = backupDatabaseRemote(dbUser, dbPassword, dbHost, "5433", dbName, localSavePath, remoteUser, remoteHost, remoteDBPort)
		} else {
			dbUser, dbPassword, dbHost, dbPort, dbName, err = parseDBConfig(dbConfig)
			if err != nil {
				logger.Fatalf("❌ Invalid database configuration format: %v", err)
			}
			err = backupDatabaseLocal(dbUser, dbPassword, dbHost, dbPort, dbName, localSavePath)
		}

		if err != nil {
			logger.Fatalf("❌ Database backup failed: %v", err)
		}

		logger.Info("✅ Database backup completed successfully")
	},
}

var (
	remoteFlag string
)

func init() {
	dbCmd.AddCommand(backupCmd)
	backupCmd.Flags().StringVar(&remoteFlag, "remote", "", "Specify remote connection in format <user>@<host>:<db_port>")
}

// parseDBConfig parses the local database configuration in the format <username>:<password>@<host>:<port>/<dbname>
func parseDBConfig(config string) (user, password, host, port, dbName string, err error) {
	// Expected format: <username>:<password>@<host>:<port>/<dbname>
	parts := strings.Split(config, "@")
	if len(parts) != 2 {
		return "", "", "", "", "", fmt.Errorf("missing or invalid user and host format")
	}

	userPass := strings.Split(parts[0], ":")
	if len(userPass) != 2 {
		return "", "", "", "", "", fmt.Errorf("missing or invalid username and password format")
	}

	user, password = userPass[0], userPass[1]
	hostPortDB := strings.Split(parts[1], ":")
	if len(hostPortDB) != 2 {
		return "", "", "", "", "", fmt.Errorf("missing or invalid host and port format")
	}

	host, portDB := hostPortDB[0], hostPortDB[1]
	portDBSplit := strings.Split(portDB, "/")
	if len(portDBSplit) != 2 {
		return "", "", "", "", "", fmt.Errorf("missing or invalid database name format")
	}

	port, dbName = portDBSplit[0], portDBSplit[1]
	return user, password, host, port, dbName, nil
}

// parseRemoteFlag parses the remote flag string in the format <user>@<host>:<db_port>
func parseRemoteFlag(remote string) (user, host, dbPort string, err error) {
	parts := strings.Split(remote, "@")
	if len(parts) != 2 {
		return "", "", "", fmt.Errorf("missing or invalid remote user and host format")
	}

	user = parts[0]
	hostPort := strings.Split(parts[1], ":")
	if len(hostPort) != 2 {
		return "", "", "", fmt.Errorf("missing or invalid host and port format")
	}

	host = hostPort[0]
	dbPort = hostPort[1]
	return user, host, dbPort, nil
}

// backupDatabaseLocal performs the database backup locally without SSH tunnel
func backupDatabaseLocal(dbUser, dbPassword, dbHost, dbPort, dbName, localSavePath string) error {
	timestamp := time.Now().Format("20060102_150405")
	backupFile := filepath.Join(localSavePath, fmt.Sprintf("%s_backup_%s.sql", dbName, timestamp))

	pgDumpCmd := exec.Command("pg_dump",
		"-h", dbHost,
		"-p", dbPort,
		"-U", dbUser,
		"-d", dbName,
		"-F", "c",
		"-f", backupFile,
	)

	pgDumpCmd.Env = append(os.Environ(), fmt.Sprintf("PGPASSWORD=%s", dbPassword))
	var stdOut, stdErr bytes.Buffer
	pgDumpCmd.Stdout = &stdOut
	pgDumpCmd.Stderr = &stdErr

	if err := pgDumpCmd.Run(); err != nil {
		return fmt.Errorf("failed to execute pg_dump: %w\nOutput: %s\nError: %s", err, stdOut.String(), stdErr.String())
	}

	fmt.Printf("✅ Backup saved to %s\n", backupFile)
	return nil
}

// backupDatabaseRemote performs the database backup over an SSH tunnel and saves it locally
func backupDatabaseRemote(dbUser, dbPassword, dbHost, dbPort, dbName, localSavePath, remoteUser, remoteHost, remoteDBPort string) error {
	timestamp := time.Now().Format("20060102_150405")
	backupFile := filepath.Join(localSavePath, fmt.Sprintf("%s_backup_%s.sql", dbName, timestamp))

	// Start SSH tunnel to forward to specified local dbPort
	tunnelCmd := exec.Command("ssh", "-fN", "-L", fmt.Sprintf("%s:%s:%s", dbPort, dbHost, remoteDBPort), fmt.Sprintf("%s@%s", remoteUser, remoteHost))
	tunnelCmd.Stdout = os.Stdout
	tunnelCmd.Stderr = os.Stderr

	if err := tunnelCmd.Run(); err != nil {
		return fmt.Errorf("failed to start SSH tunnel: %w", err)
	}

	defer func() {
		if err := killProcessOnPort(dbPort); err != nil {
			fmt.Printf("❌ Failed to kill SSH tunnel process on port %s: %v\n", dbPort, err)
		} else {
			fmt.Println("✅ SSH tunnel process on port", dbPort, "terminated successfully.")
		}
	}()

	pgDumpCmd := exec.Command("pg_dump",
		"-h", "localhost",
		"-p", dbPort,
		"-U", dbUser,
		"-d", dbName,
		"-F", "c",
		"-f", backupFile,
	)

	pgDumpCmd.Env = append(os.Environ(), fmt.Sprintf("PGPASSWORD=%s", dbPassword))
	var stdOut, stdErr bytes.Buffer
	pgDumpCmd.Stdout = &stdOut
	pgDumpCmd.Stderr = &stdErr

	if err := pgDumpCmd.Run(); err != nil {
		return fmt.Errorf("failed to execute pg_dump: %w\nOutput: %s\nError: %s", err, stdOut.String(), stdErr.String())
	}

	fmt.Printf("✅ Backup saved to %s\n", backupFile)
	return nil
}
