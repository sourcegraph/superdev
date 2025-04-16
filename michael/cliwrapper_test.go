package main

import (
	"fmt"
	"testing"
	"time"
)

// Test running the Amp worker
func TestRunPrompt(t *testing.T) {
	// Skip the test if in CI or no CLI available

	// fmt.Println("Starting Amp worker with prompt...")

	// Start a thread with a prompt
	prompt := "Generate a Hello World program in Go"
	ch, cleanup, err := StartThreadWithPrompt(prompt, 5*time.Second) // Shorter duration for tests
	if err != nil {
		t.Fatalf("Error starting thread: %v", err)
	}
	defer cleanup()

	// Process items from the channel
	fmt.Println("Processing updates...")
	count := 0
	for item := range ch {
		// Print the rendered item
		fmt.Println(item.Render())
		count++
	}

	fmt.Println("amp Done!")
	// At least we should get some response
	if count == 0 {
		t.Fatal("No items received from the worker")
	}
}
