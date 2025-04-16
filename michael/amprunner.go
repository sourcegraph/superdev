package main

import (
	"bufio"
	"database/sql"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/spf13/cobra"
	_ "modernc.org/sqlite" // Import for side effects
)

var rootCmd = &cobra.Command{
	Use:   "superdev-amprunner",
	Short: "Runs Amp with interactivity",
}

var dbPath string

var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Runs Amp reading from and writing to an SQLite database",
	Run: func(cmd *cobra.Command, args []string) {
		if err := runAmpWithSQLite(dbPath); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
	},
}

func init() {
	runCmd.Flags().StringVar(&dbPath, "db", "amp.db", "Path to the SQLite database file")
	rootCmd.AddCommand(runCmd)
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

// runAmpWithSQLite reads from SQLite DB, sends content to amp CLI,
// and writes output back to the DB
func runAmpWithSQLite(dbPath string) error {
	if _, err := os.Stat(dbPath); os.IsNotExist(err) {
		fmt.Println("DB doesn't exist yet. Creating...")
	}

	// Ensure the database exists and has the correct schema
	db, err := initializeDatabase(dbPath)
	if err != nil {
		return fmt.Errorf("failed to initialize database: %w", err)
	}
	defer db.Close()

	// Get the last ID to use as a starting point
	var lastID int64
	//var count int
	//row := db.QueryRow("SELECT COUNT(*), COALESCE(MAX(id), 0) FROM messages")
	//if err := row.Scan(&count, &lastID); err != nil {
	//	return fmt.Errorf("failed to get last message ID: %w", err)
	//}

	// Main processing loop
	for {
		// Check for new input messages
		newInputs, err := collectNewMessages(db, lastID)
		if err != nil {
			return err
		}

		// Process each new input message
		for _, input := range newInputs {
			fmt.Printf("Processing input #%d: %s\n", input.id, input.content)
			lastID = input.id

			// Create and set up the amp command
			cmd := exec.Command("amp")
			cmd.Stdin = bufio.NewReader(strings.NewReader(input.content))

			// Capture stdout
			stdout, err := cmd.StdoutPipe()
			if err != nil {
				return fmt.Errorf("failed to create stdout pipe: %w", err)
			}

			// Start the command
			if err := cmd.Start(); err != nil {
				return fmt.Errorf("failed to start amp command: %w", err)
			}

			// Read all output
			outputBytes, err := io.ReadAll(stdout)
			if err != nil {
				return fmt.Errorf("failed to read amp output: %w", err)
			}
			output := string(outputBytes)

			// Wait for the command to finish
			if err := cmd.Wait(); err != nil {
				return fmt.Errorf("amp command failed: %w", err)
			}

			// Save output to database with fresh connection
			_, err = db.Exec("INSERT INTO messages (direction, content, thread_id) VALUES (?, ?, ?)", "output", output, "1")
			if err != nil {
				return fmt.Errorf("failed to save output to database: %w", err)
			}

			fmt.Printf("Saved output to database: %s\n", output)
		}

		// Sleep before next check
		time.Sleep(1000 * time.Millisecond)
	}
}

func collectNewMessages(db *sql.DB, lastID int64) ([]struct {
	id      int64
	content string
}, error) {
	fmt.Println("Checking for new messages at ", time.Now().Format("2006-01-02 15:04:05"))

	rows, err := db.Query("SELECT id, content FROM messages WHERE id > ? AND direction = 'input' AND thread_id = '1' ORDER BY id ASC", lastID)
	if err != nil {
		return nil, fmt.Errorf("failed to query for new messages: %w", err)
	}
	defer rows.Close()

	var newInputs []struct {
		id      int64
		content string
	}

	for rows.Next() {
		var id int64
		var content string
		if err := rows.Scan(&id, &content); err != nil {
			return nil, fmt.Errorf("failed to scan message row: %w", err)
		}
		newInputs = append(newInputs, struct {
			id      int64
			content string
		}{id, content})
	}

	fmt.Printf("Found %d new messages\n", len(newInputs))
	return newInputs, nil
}

// initializeDatabase ensures the database exists and has the correct schema
func initializeDatabase(dbPath string) (*sql.DB, error) {
	// Open database connection
	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	// Create the messages table if it doesn't exist
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS messages (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			thread_id TEXT NOT NULL,
			direction TEXT NOT NULL,
			content TEXT NOT NULL,
			timestamp DATETIME DEFAULT CURRENT_TIMESTAMP
		)
	`)
	if err != nil {
		db.Close()
		return nil, fmt.Errorf("failed to create messages table: %w", err)
	}

	return db, nil
}
