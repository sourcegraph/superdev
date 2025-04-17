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
		// fmt.Println(item.Render())
		fmt.Println(item)
		count++
	}

	fmt.Println("amp Done!")
	// At least we should get some response
	if count == 0 {
		t.Fatal("No items received from the worker")
	}
}

// Test running the Amp worker with conversation continuity
func TestContinueConversation(t *testing.T) {
	// Skip the test if in CI or no CLI available

	// =====================================================
	// PART 1: Start the first conversation
	// =====================================================
	fmt.Println("\n\n==== STARTING FIRST CONVERSATION ====")

	// First prompt to start the conversation
	firstPrompt := "Create a 2-line poem about programming"
	ch1, cleanup1, err := StartThreadWithPrompt(firstPrompt, 5*time.Second)
	if err != nil {
		t.Fatalf("Error starting first thread: %v", err)
	}
	defer cleanup1()

	// Collect all messages from the first conversation
	fmt.Println("Processing first conversation responses...")
	var messages []AmpMessage
	var threadID string
	count1 := 0
	for item := range ch1 {
		// Print the item
		fmt.Println(item.Render())

		// Store thread messages and ID if available
		if thread, ok := item.(AmpThread); ok {
			// Get thread ID and messages for continuation
			threadID = thread.ID
			messages = append(messages, thread.Messages...)
		}
		count1++
	}

	// Verify we got a response
	if count1 == 0 {
		t.Fatal("No items received from the first conversation")
	}

	fmt.Println("\n\n==== FIRST CONVERSATION COMPLETED ====")
	fmt.Println("Thread ID captured for continuation:", threadID)
	fmt.Println("Number of messages captured:", len(messages))
	fmt.Println()

	// =====================================================
	// PART 2: Continue with a second conversation
	// =====================================================
	fmt.Println("\n\n==== STARTING SECOND CONVERSATION (CONTINUATION) ====")

	// Start a second conversation that depends on the first one
	secondPrompt := "What are the first and last words of the poem you just created?"
	// Convert messages to the expected type
	ch2, cleanup2, err := ContinueThreadWithPrompt(secondPrompt, threadID, messages, 5*time.Second)
	if err != nil {
		t.Fatalf("Error starting second thread: %v", err)
	}
	defer cleanup2()

	// Process items from the second conversation
	fmt.Println("Processing second conversation responses...")
	count2 := 0
	for item := range ch2 {
		// Print the rendered item
		fmt.Println(item.Render())
		count2++
	}

	fmt.Println("\n\n==== SECOND CONVERSATION COMPLETED ====")
	// Verify we got a response from the second conversation too
	if count2 == 0 {
		t.Fatal("No items received from the second conversation")
	}
}
