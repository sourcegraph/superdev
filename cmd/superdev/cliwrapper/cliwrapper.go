package superdev

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"os/exec"
	"sync"
	"time"
)

// Simple RPC client for Amp CLI worker
type AmpClient struct {
	cmd    *exec.Cmd
	stdin  io.WriteCloser
	stdout *bufio.Scanner
	nextID int
}

// NewAmpClient creates a new client for the Amp CLI worker
func NewAmpClient() (*AmpClient, error) {
	// Create command
	cmd := exec.Command("amp", "worker")

	// Get stdin and stdout pipes
	stdin, err := cmd.StdinPipe()
	if err != nil {
		return nil, fmt.Errorf("failed to get stdin pipe: %w", err)
	}

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return nil, fmt.Errorf("failed to get stdout pipe: %w", err)
	}

	// Set command to use stderr directly
	cmd.Stderr = nil

	// Start the command
	if err := cmd.Start(); err != nil {
		return nil, fmt.Errorf("failed to start worker: %w", err)
	}

	// Create scanner for reading responses
	scanner := bufio.NewScanner(stdout)

	return &AmpClient{
		cmd:    cmd,
		stdin:  stdin,
		stdout: scanner,
		nextID: 1,
	}, nil
}

// Call sends an RPC request to the worker and returns the raw response
func (c *AmpClient) Call(method string, args []interface{}) (string, error) {
	// Create request
	request := map[string]interface{}{
		"streamId": c.nextID,
		"method":   method,
		"args":     args,
	}
	c.nextID++

	// Marshal and send request
	requestJSON, err := json.Marshal(request)
	if err != nil {
		return "", fmt.Errorf("failed to marshal request: %w", err)
	}

	if _, err := c.stdin.Write(append(requestJSON, '\n')); err != nil {
		return "", fmt.Errorf("failed to write request: %w", err)
	}

	// Read response
	if !c.stdout.Scan() {
		if err := c.stdout.Err(); err != nil {
			return "", fmt.Errorf("failed to read response: %w", err)
		}
		return "", fmt.Errorf("unexpected EOF")
	}

	return c.stdout.Text(), nil
}

// Shutdown closes the client
func (c *AmpClient) Shutdown() error {
	// Close stdin
	if err := c.stdin.Close(); err != nil {
		return fmt.Errorf("failed to close stdin: %w", err)
	}

	// Wait for command to exit
	if err := c.cmd.Wait(); err != nil {
		return fmt.Errorf("worker exited with error: %w", err)
	}

	return nil
}

// StartThreadWithPrompt starts a thread with a prompt and returns a channel of AmpItem
func StartThreadWithPrompt(prompt string, duration time.Duration) (<-chan AmpItem, func(), error) {
	// Create channel for items
	ch := make(chan AmpItem, 100) // Buffer to prevent blocking

	// Create client
	client, err := NewAmpClient()
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create client: %w", err)
	}

	// Create thread ID
	threadId := fmt.Sprintf("thread_%d", time.Now().Unix())

	// Start thread worker
	response, err := client.Call("startThreadWorker", []interface{}{threadId})
	if err != nil {
		client.Shutdown()
		return nil, nil, fmt.Errorf("failed to start thread worker: %w", err)
	}

	// Parse the response to check if it was successful
	var respObj AmpWorkerResponse
	if err := json.Unmarshal([]byte(response), &respObj); err != nil {
		client.Shutdown()
		return nil, nil, fmt.Errorf("failed to parse start response: %w", err)
	}

	if respObj.StreamEvent != "next" {
		client.Shutdown()
		return nil, nil, fmt.Errorf("unexpected start response: %s", response)
	}

	// Create user message payload
	delta := map[string]interface{}{
		"type": "user:message",
		"message": map[string]interface{}{
			"content": []map[string]interface{}{
				{
					"type": "text",
					"text": prompt,
				},
			},
		},
	}

	// Send user message
	response, err = client.Call("handleThreadDelta", []interface{}{threadId, delta})
	if err != nil {
		client.Shutdown()
		return nil, nil, fmt.Errorf("failed to send user message: %w", err)
	}

	// Start observing thread
	_, err = client.Call("observeThread", []interface{}{threadId})
	if err != nil {
		client.Shutdown()
		return nil, nil, fmt.Errorf("failed to observe thread: %w", err)
	}

	// Start goroutine to process responses
	go func() {
		defer close(ch)

		startTime := time.Now()
		done := false
		for time.Since(startTime) < duration && !done {
			// Read next response
			if !client.stdout.Scan() {
				break
			}

			response := client.stdout.Text()

			// Parse response as typed AmpWorkerResponse
			var respObj AmpWorkerResponse
			if err := json.Unmarshal([]byte(response), &respObj); err != nil {
				ch <- AmpGenericItem{fmt.Sprintf("Error parsing response: %v", err)}
				continue
			}

			// Check for thread data
			if respObj.StreamEvent == "next" && respObj.Data != nil {
				// First, try to parse as a AmpThreadState
				var threadState AmpThreadState
				dataBytes, _ := json.Marshal(respObj.Data)

				if err := json.Unmarshal(dataBytes, &threadState); err == nil && threadState.State != "" {
					// It's a AmpThreadState
					ch <- threadState

					// If the inference state is idle, we're done
					if threadState.State == "active" && threadState.InferenceState == "idle" {
						// Thread is complete, break from the loop
						done = true
					}
				} else {
					// Try as a AmpThread
					var thread AmpThread
					if err := json.Unmarshal(dataBytes, &thread); err == nil && thread.ID != "" {
						// It's a AmpThread
						ch <- thread

						// Check if the thread has a completed state
						if thread.State == "active" && thread.InferenceState == "idle" {
							// Thread is complete, break from the loop
							done = true
						}
					} else {
						// Use generic for other types
						var prettyData interface{}
						json.Unmarshal(dataBytes, &prettyData)
						ch <- AmpGenericItem{prettyData}
					}
				}
			}
		}
	}()

	// Return cleanup function
	cleanup := func() {
		client.Shutdown()
	}

	return ch, cleanup, nil
}

// ContinueThreadWithPrompt continues a thread with a prompt and existing messages
func ContinueThreadWithPrompt(prompt string, threadID string, messages []AmpMessage, duration time.Duration) (<-chan AmpItem, func(), error) {
	// Create channel for items
	ch := make(chan AmpItem, 100) // Buffer to prevent blocking

	// Create client
	client, err := NewAmpClient()
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create client: %w", err)
	}

	// Start thread worker
	response, err := client.Call("startThreadWorker", []interface{}{threadID})
	if err != nil {
		client.Shutdown()
		return nil, nil, fmt.Errorf("failed to start thread worker: %w", err)
	}

	// Parse the response to check if it was successful
	var respObj AmpWorkerResponse
	if err := json.Unmarshal([]byte(response), &respObj); err != nil {
		client.Shutdown()
		return nil, nil, fmt.Errorf("failed to parse start response: %w", err)
	}

	if respObj.StreamEvent != "next" {
		client.Shutdown()
		return nil, nil, fmt.Errorf("unexpected start response: %s", response)
	}

	// filter out messages that are 'thinking'
	finishedMessages := FilterMessages(messages)

	// Serialize all previous messages into a single string
	history := SerializeMessages(finishedMessages)

	// Create single user message with both history and new prompt
	combinedPrompt := fmt.Sprintf("Previous conversation:\n\n%s\n\nNew question: %s", history, prompt)

	// Create the delta with the combined message
	delta := map[string]interface{}{
		"type": "user:message",
		"message": map[string]interface{}{
			"content": []map[string]interface{}{
				{
					"type": "text",
					"text": combinedPrompt,
				},
			},
		},
	}

	// Send user message
	response, err = client.Call("handleThreadDelta", []interface{}{threadID, delta})
	if err != nil {
		client.Shutdown()
		return nil, nil, fmt.Errorf("failed to send user message: %w", err)
	}

	// Start observing thread
	_, err = client.Call("observeThread", []interface{}{threadID})
	if err != nil {
		client.Shutdown()
		return nil, nil, fmt.Errorf("failed to observe thread: %w", err)
	}

	// Set up a context with cancellation for managing goroutines
	var wg sync.WaitGroup
	wg.Add(1)

	// Start goroutine to process responses
	go func() {
		defer wg.Done()
		defer close(ch)

		startTime := time.Now()
		done := false
		for time.Since(startTime) < duration && !done {
			// Read next response
			if !client.stdout.Scan() {
				break
			}

			response := client.stdout.Text()

			// Parse response as typed AmpWorkerResponse
			var respObj AmpWorkerResponse
			if err := json.Unmarshal([]byte(response), &respObj); err != nil {
				ch <- AmpGenericItem{fmt.Sprintf("Error parsing response: %v", err)}
				continue
			}

			// Check for thread data
			if respObj.StreamEvent == "next" && respObj.Data != nil {
				// First, try to parse as a AmpThreadState
				var threadState AmpThreadState
				dataBytes, _ := json.Marshal(respObj.Data)

				if err := json.Unmarshal(dataBytes, &threadState); err == nil && threadState.State != "" {
					// It's a AmpThreadState
					ch <- threadState

					// If the inference state is idle, we're done
					if threadState.State == "active" && threadState.InferenceState == "idle" {
						// Thread is complete, break from the loop
						done = true
					}
				} else {
					// Try as a AmpThread
					var thread AmpThread
					if err := json.Unmarshal(dataBytes, &thread); err == nil && thread.ID != "" {
						// It's a AmpThread
						ch <- thread

						// Check if the thread has a completed state
						if thread.State == "active" && thread.InferenceState == "idle" {
							// Thread is complete, break from the loop
							done = true
						}
					} else {
						// Use generic for other types
						var prettyData interface{}
						json.Unmarshal(dataBytes, &prettyData)
						ch <- AmpGenericItem{prettyData}
					}
				}
			}
		}
	}()

	// Return cleanup function
	cleanup := func() {
		client.Shutdown()
		wg.Wait()
	}

	return ch, cleanup, nil
}
