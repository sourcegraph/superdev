package superdev

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	superdev "superdev/cmd/superdev/cliwrapper"
	"time"
)

// runThreadCmd creates a command to run the AMP client with specified parameters
func runThreadCmd() error {
	// Define command-line flags
	threadCmd := flag.NewFlagSet("thread", flag.ExitOnError)
	prompt := threadCmd.String("prompt", "", "The prompt to send to the model (required)")
	outputPath := threadCmd.String("output", "", "Path to output file for thread messages (required)")
	previousPath := threadCmd.String("previous", "", "Path to file containing previous messages (optional)")
	timeoutSecs := threadCmd.Int("timeout", 60, "Timeout in seconds for the thread (default: 60s)")

	// Parse command-line flags
	threadCmd.Parse(os.Args[2:])

	// Check required flags
	if *prompt == "" {
		return fmt.Errorf("--prompt is required")
	}

	if *outputPath == "" {
		return fmt.Errorf("--output is required")
	}

	// Set timeout
	duration := time.Duration(*timeoutSecs) * time.Second

	var ch <-chan superdev.AmpItem
	var cleanup func()
	var err error

	// Check if we have previous messages
	if *previousPath != "" {
		// Read and parse previous messages
		previousData, err := ioutil.ReadFile(*previousPath)
		if err != nil {
			return fmt.Errorf("failed to read previous messages file: %w", err)
		}

		// Unmarshal messages
		var messages []superdev.AmpMessage
		if err := json.Unmarshal(previousData, &messages); err != nil {
			return fmt.Errorf("failed to parse previous messages: %w", err)
		}

		// Get thread ID from the first message if available
		threadID := fmt.Sprintf("thread_%d", time.Now().Unix()) // Default thread ID

		// Continue thread with prompt
		fmt.Println("Continuing thread with prompt...")
		ch, cleanup, err = superdev.ContinueThreadWithPrompt(*prompt, threadID, messages, duration)
	} else {
		// Start a new thread with prompt
		fmt.Println("Starting new thread with prompt...")
		ch, cleanup, err = superdev.StartThreadWithPrompt(*prompt, duration)
	}

	if err != nil {
		return fmt.Errorf("failed to start thread: %w", err)
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
		return fmt.Errorf("failed to marshal messages: %w", err)
	}

	if err := ioutil.WriteFile(*outputPath, messagesData, 0644); err != nil {
		return fmt.Errorf("failed to write output file: %w", err)
	}

	fmt.Println("Messages saved to", *outputPath)
	return nil
}
