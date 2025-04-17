package superdev

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"strings"
	"sync"
	"time"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "superdev-amprunner",
	Short: "Runs Amp with interactivity",
}

var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Runs Amp reading from and writing to a remote server",
	Run: func(cmd *cobra.Command, args []string) {
		if err := runAmpWithServer(); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(runCmd)
}

// runAmpWithServer reads from a remote server, sends content to amp CLI,
// and writes output back to the server
func runAmpWithServer() error {
	serverURL := os.Getenv("SERVER_URL")
	if serverURL == "" {
		return fmt.Errorf("SERVER_URL environment variable is not set")
	}

	threadID := os.Getenv("THREAD_ID")
	if threadID == "" {
		return fmt.Errorf("THREAD_ID environment variable is not set")
	}

	// Variable to track the last message ID we've processed
	var lastMessageID string

	var lock sync.Mutex

	// Main processing loop
	for {
		// Check for new input messages
		lock.Lock()
		newMessages, err := pullMessages(serverURL, threadID, lastMessageID)
		if err != nil {
			return err
		}
		lock.Unlock()

		// Process each new input message
		for _, input := range newMessages {
			fmt.Printf("Processing input: %s\n", input.Content)
			lastMessageID = input.ID

			// todo: we need to set that up earlier and then pipe input and output
			// Create and set up the amp command
			cmd := exec.Command("amp")
			cmd.Stdin = bufio.NewReader(strings.NewReader(input.Content))

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

			// Send output to server
			newLastMessageID, err := answerMessage(serverURL, threadID, output)
			if err != nil {
				return fmt.Errorf("failed to send output to server: %w", err)
			}

			// Update last message ID
			lastMessageID = newLastMessageID

			fmt.Printf("Sent output to server: %s\n", output)
		}

		// Sleep before next check
		time.Sleep(time.Second)
	}
}

// Message represents a message from the server
type Message struct {
	ID      string
	Content string
}

// pullMessages fetches new messages from the server
func pullMessages(serverURL, threadID, lastMessageID string) ([]Message, error) {
	fmt.Println("Checking server for new messages at", time.Now().Format("2006-01-02 15:04:05"))

	// Build the URL with query parameters
	baseURL, err := url.Parse(serverURL + "/pullMessages")
	if err != nil {
		return nil, fmt.Errorf("invalid server URL: %w", err)
	}

	// Add query parameters
	params := url.Values{}
	params.Add("thread_id", threadID)
	if lastMessageID != "" {
		params.Add("last_message_id", lastMessageID)
	}
	baseURL.RawQuery = params.Encode()

	// Make the request
	resp, err := http.Get(baseURL.String())
	if err != nil {
		return nil, fmt.Errorf("failed to fetch messages: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("server returned error: %s", resp.Status)
	}

	// Parse the response
	var messages []Message
	if err := json.NewDecoder(resp.Body).Decode(&messages); err != nil {
		return nil, fmt.Errorf("failed to decode server response: %w", err)
	}

	fmt.Printf("Found %d new messages\n", len(messages))
	return messages, nil
}

// answerMessage sends the amp output back to the server
func answerMessage(serverURL, threadID, output string) (string, error) {
	// Prepare the request payload
	payload := struct {
		ThreadID string `json:"thread_id"`
		Payload  string `json:"payload"`
	}{
		ThreadID: threadID,
		Payload:  output,
	}

	// Convert payload to JSON
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return "", fmt.Errorf("failed to marshal payload: %w", err)
	}

	// Make POST request
	resp, err := http.Post(
		serverURL+"/answerMessage",
		"application/json",
		bytes.NewBuffer(payloadBytes),
	)
	if err != nil {
		return "", fmt.Errorf("failed to send answer to server: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("server returned error: %s", resp.Status)
	}

	// Parse response to get lastMessageId
	var response struct {
		MessageID string `json:"message_id"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return "", fmt.Errorf("failed to decode server response: %w", err)
	}

	return response.MessageID, nil
}
