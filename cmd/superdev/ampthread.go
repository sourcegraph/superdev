package superdev

import (
	"encoding/json"
	"fmt"
	"os"
	superdev "superdev/cmd/superdev/cliwrapper"
	"time"

	"github.com/spf13/cobra"
)

// Command variables
var (
	promptText   string
	outputPath   string
	previousPath string
	timeoutSecs  int
)

var threadCmd = &cobra.Command{
	Use:   "thread",
	Short: "Start or continue an AMP thread with a prompt",
	Run: func(cmd *cobra.Command, args []string) {
		// Check required flags
		if promptText == "" {
			fmt.Println("Error: --prompt is required")
			os.Exit(1)
		}

		if outputPath == "" {
			fmt.Println("Error: --output is required")
			os.Exit(1)
		}

		// Set timeout
		duration := time.Duration(timeoutSecs) * time.Second

		var ch <-chan superdev.AmpItem
		var cleanup func()
		var err error

		// Check if we have previous messages
		if previousPath != "" {
			// Read and parse previous messages
			previousData, err := os.ReadFile(previousPath)
			if err != nil {
				fmt.Printf("Error: failed to read previous messages file: %v\n", err)
				os.Exit(1)
			}

			// Unmarshal messages
			var messages []superdev.AmpMessage
			if err := json.Unmarshal(previousData, &messages); err != nil {
				fmt.Printf("Error: failed to parse previous messages: %v\n", err)
				os.Exit(1)
			}

			// Get thread ID from the first message if available
			threadID := fmt.Sprintf("thread_%d", time.Now().Unix()) // Default thread ID

			// Continue thread with prompt
			fmt.Println("Continuing thread with prompt...")
			ch, cleanup, err = superdev.ContinueThreadWithPrompt(promptText, threadID, messages, duration)
		} else {
			// Start a new thread with prompt
			fmt.Println("Starting new thread with prompt...")
			ch, cleanup, err = superdev.StartThreadWithPrompt(promptText, duration)
		}

		if err != nil {
			fmt.Printf("Error: failed to start thread: %v\n", err)
			os.Exit(1)
		}

		defer cleanup()

		// Collect all messages
		var allMessages []superdev.AmpMessage
		var threadID string

		// Process items from the channel
		fmt.Println("Processing responses...")
		for item := range ch {
			// Print update
			fmt.Println(item.Render())

			// Store thread messages and ID if available
			if thread, ok := item.(superdev.AmpThread); ok {
				// Get thread ID and messages for continuation
				threadID = thread.ID
				allMessages = append(allMessages, thread.Messages...)
			}
		}

		fmt.Println("Thread completed!")
		fmt.Println("Thread ID:", threadID)
		fmt.Println("Total messages:", len(allMessages))

		// Save messages to output file
		messagesData, err := json.MarshalIndent(allMessages, "", "  ")
		if err != nil {
			fmt.Printf("Error: failed to marshal messages: %v\n", err)
			os.Exit(1)
		}

		if err := os.WriteFile(outputPath, messagesData, 0644); err != nil {
			fmt.Printf("Error: failed to write output file: %v\n", err)
			os.Exit(1)
		}

		fmt.Println("Messages saved to", outputPath)
	},
}
