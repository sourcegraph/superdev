package superdev

import (
	"bufio"
	"bytes"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"sync"
	"time"

	"github.com/spf13/cobra"
)

// ThreadMessage stores the output for each thread
type ThreadMessage struct {
	ID        string
	Output    string
	Direction string
	Status    string    // "processing", "completed", or "error"
	CreatedAt time.Time // For cleanup purposes
	Error     string    // Error message if status is "error"
}

// In-memory storage for thread outputs
var (
	threads          = make(map[string][]*ThreadMessage)
	threadContainers = make(map[string]string)
	outputMutex      = &sync.Mutex{}
	maxOutputAge     = 24 * time.Hour // Outputs older than this will be cleaned up
)

// generateThreadID creates a unique thread ID
func generateThreadID() (string, error) {
	bytes := make([]byte, 16)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

// serverCmd handles server operations
// corsMiddleware handles CORS preflight requests and adds CORS headers
func corsMiddleware(handler http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Set CORS headers
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		// Handle preflight requests
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		// Call the actual handler
		handler(w, r)
	}
}

var serverCmd = &cobra.Command{
	Use:   "server [port]",
	Short: "Start a server that accepts Docker build requests",
	Args:  cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		port := "8080"
		if len(args) > 0 {
			port = args[0]
		}

		// Check for ANTHROPIC_API_KEY environment variable
		if os.Getenv("ANTHROPIC_API_KEY") == "" {
			fmt.Println("WARNING: ANTHROPIC_API_KEY environment variable is not set.")
			fmt.Println("AMP execution will fail without an API key. Set the environment variable before starting the server.")
		}

		// Run the server command
		fmt.Printf("Starting server on port %s...\n", port)

		// Setup HTTP server
		// Start a thread for a new conversation
		http.HandleFunc("/start", corsMiddleware(handleStartContainerRequest))
		// Write a human message
		http.HandleFunc("/storeMessage", corsMiddleware(handleStoreMessageRequest))
		// Worker pulls message
		http.HandleFunc("/pullMessages", corsMiddleware(handlePullMessagesRequest))
		// Worker sends message response
		http.HandleFunc("/answerMessage", corsMiddleware(handleAnswerMessageRequest))

		http.HandleFunc("/output", corsMiddleware(handleOutputRequest))
		http.HandleFunc("/threads", corsMiddleware(handleThreadsRequest))

		fmt.Printf("Server started on :%s\n", port)
		err := http.ListenAndServe(":"+port, nil)
		if err != nil {
			fmt.Printf("Error starting server: %v\n", err)
			os.Exit(1)
		}

		// todo: clean up running containers when we exit
	},
}

type Message struct {
	ID      string
	Content string
}

func handlePullMessagesRequest(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Get thread ID from query parameter
	threadID := r.URL.Query().Get("thread_id")
	if threadID == "" {
		http.Error(w, "Missing thread_id parameter", http.StatusBadRequest)
		return
	}

	// Get last message ID from query parameter
	lastMessageID := r.URL.Query().Get("last_message_id")

	// Check if thread exists
	outputMutex.Lock()
	defer outputMutex.Unlock()

	_, exists := threads[threadID]
	if !exists {
		var response []Message
		json.NewEncoder(w).Encode(response)
		return
	}

	// Filter messages: direction = "input" and after lastMessageID
	var messages []*ThreadMessage
	for _, msg := range threads[threadID] {
		// Filter by direction
		if msg.Direction != "input" {
			continue
		}

		// Filter by lastMessageID if provided
		if lastMessageID != "" && msg.ID <= lastMessageID {
			continue
		}

		messages = append(messages, msg)
	}

	// Prepare response
	w.Header().Set("Content-Type", "application/json")

	var response []Message
	for _, msg := range messages {
		response = append(response, Message{
			ID:      msg.ID,
			Content: msg.Output,
		})
	}

	json.NewEncoder(w).Encode(response)
}

func handleAnswerMessageRequest(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Error reading request body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	var req struct {
		ThreadId string `json:"thread_id"`
		Payload  string `json:"payload"`
	}

	if err := json.Unmarshal(body, &req); err != nil {
		http.Error(w, "Error parsing JSON", http.StatusBadRequest)
		return
	}

	if req.Payload == "" {
		http.Error(w, "Payload is required", http.StatusBadRequest)
		return
	}

	if req.ThreadId == "" {
		http.Error(w, "ThreadId is required", http.StatusBadRequest)
		return
	}

	if threads[req.ThreadId] == nil {
		http.Error(w, "Thread history for threadId not found", http.StatusNotFound)
		return
	}

	messageId := time.Now().String()
	threads[req.ThreadId] = append(threads[req.ThreadId], &ThreadMessage{
		ID:        messageId,
		Direction: "output",
		Output:    req.Payload,
		CreatedAt: time.Now(),
	})

	w.Header().Set("Content-Type", "application/json")
	response := map[string]string{
		"message_id": messageId,
	}

	json.NewEncoder(w).Encode(response)

	// Create new thread output entry
	//outputMutex.Lock()
	//threads[threadID] = &ThreadMessage{
	//	ID:        threadID,
	//	Status:    "processing",
	//	CreatedAt: time.Now(),
	//}
	//outputMutex.Unlock()

	// Add thread ID to response
	//response := map[string]string{
	//	"status":    "success",
	//	"message":   "Request processed successfully",
	//	"thread_id": threadID,
	//}
	//
	//// Respond to client immediately with thread ID
	//w.Header().Set("Content-Type", "application/json")
	//json.NewEncoder(w).Encode(response)
	//
	//fmt.Printf("===== Starting thread %s processing =====\n", threadID)
	//os.Stdout.Sync()
	//// Process Docker execution in a goroutine
	//go func() {
	//	// Set up Docker run command
	//	output, err := startDockerContainer(threadID, req.RepositoryLink, req.ContextFiles, req.Prompt, req.DockerImage)
	//
	//	// Update thread output in memory
	//	outputMutex.Lock()
	//	defer outputMutex.Unlock()
	//
	//	// Check if thread output still exists (might have been cleaned up)
	//	threadOutput, exists := threads[threadID]
	//	if !exists {
	//		fmt.Printf("Warning: Thread %s was cleaned up before processing completed\n", threadID)
	//		return
	//	}
	//
	//	if err != nil {
	//		threadOutput.Status = "error"
	//		threadOutput.Error = err.Error()
	//		fmt.Printf("Error running Docker container for thread %s: %v\n", threadID, err)
	//	} else {
	//		threadOutput.Status = "completed"
	//		threadOutput.Output = output
	//		// Log output to console for now
	//		fmt.Printf("Thread %s completed with output:\n%s\n", threadID, output)
	//	}
	//}()
}

func handleStoreMessageRequest(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Error reading request body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	var req struct {
		ThreadId string `json:"thread_id"`
		Prompt   string `json:"prompt"`
	}

	if err := json.Unmarshal(body, &req); err != nil {
		http.Error(w, "Error parsing JSON", http.StatusBadRequest)
		return
	}

	if req.Prompt == "" {
		http.Error(w, "Prompt is required", http.StatusBadRequest)
		return
	}

	if req.ThreadId == "" {
		http.Error(w, "ThreadId is required", http.StatusBadRequest)
		return
	}

	if threadContainers[req.ThreadId] == "" {
		http.Error(w, "Container for threadId not found", http.StatusNotFound)
		return
	}
	if threads[req.ThreadId] == nil {
		threads[req.ThreadId] = make([]*ThreadMessage, 0)
	}

	threads[req.ThreadId] = append(threads[req.ThreadId], &ThreadMessage{
		ID:        time.Now().String(),
		Direction: "input",
		Output:    req.Prompt,
		CreatedAt: time.Now(),
	})

	// Create new thread output entry
	//outputMutex.Lock()
	//threads[threadID] = &ThreadMessage{
	//	ID:        threadID,
	//	Status:    "processing",
	//	CreatedAt: time.Now(),
	//}
	//outputMutex.Unlock()

	// Add thread ID to response
	//response := map[string]string{
	//	"status":    "success",
	//	"message":   "Request processed successfully",
	//	"thread_id": threadID,
	//}
	//
	//// Respond to client immediately with thread ID
	//w.Header().Set("Content-Type", "application/json")
	//json.NewEncoder(w).Encode(response)
	//
	//fmt.Printf("===== Starting thread %s processing =====\n", threadID)
	//os.Stdout.Sync()
	//// Process Docker execution in a goroutine
	//go func() {
	//	// Set up Docker run command
	//	output, err := startDockerContainer(threadID, req.RepositoryLink, req.ContextFiles, req.Prompt, req.DockerImage)
	//
	//	// Update thread output in memory
	//	outputMutex.Lock()
	//	defer outputMutex.Unlock()
	//
	//	// Check if thread output still exists (might have been cleaned up)
	//	threadOutput, exists := threads[threadID]
	//	if !exists {
	//		fmt.Printf("Warning: Thread %s was cleaned up before processing completed\n", threadID)
	//		return
	//	}
	//
	//	if err != nil {
	//		threadOutput.Status = "error"
	//		threadOutput.Error = err.Error()
	//		fmt.Printf("Error running Docker container for thread %s: %v\n", threadID, err)
	//	} else {
	//		threadOutput.Status = "completed"
	//		threadOutput.Output = output
	//		// Log output to console for now
	//		fmt.Printf("Thread %s completed with output:\n%s\n", threadID, output)
	//	}
	//}()
}

// handleStartContainerRequest processes Docker build requests
func handleStartContainerRequest(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Error reading request body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	var req struct {
		RepositoryLink string   `json:"repository_link"`
		ContextFiles   [][]byte `json:"contextFiles,omitempty"`
		DockerImage    string   `json:"docker_image,omitempty"`
		ServerUrl      string   `json:"server_url,omitempty"`
		Prompt         string   `json:"prompt,omitempty"`
	}

	if err := json.Unmarshal(body, &req); err != nil {
		http.Error(w, "Error parsing JSON", http.StatusBadRequest)
		return
	}

	// Validate required fields
	if req.DockerImage == "" {
		// Default to the standard image if not specified
		http.Error(w, "Docker image is required", http.StatusBadRequest)
		return
	}

	if req.RepositoryLink == "" {
		http.Error(w, "Repository link is required", http.StatusBadRequest)
		return
	}

	if req.ServerUrl == "" {
		req.ServerUrl = "http://localhost:8080"
	}

	// Process the request
	fmt.Printf("Received request: Docker image: %s, Repo: %s, Context files count: %d\n",
		req.DockerImage, req.RepositoryLink, len(req.ContextFiles))

	// We'll respond with thread ID later

	// Execute Docker command and store the output
	// Generate a unique thread ID
	threadID, err := generateThreadID()
	if err != nil {
		http.Error(w, "Error generating thread ID", http.StatusInternalServerError)
		return
	}

	dockerContainerId, err := startDockerContainer(threadID, req.RepositoryLink, req.ContextFiles, req.DockerImage, req.ServerUrl)
	threadContainers[threadID] = strings.ReplaceAll(dockerContainerId, "\n", "")

	fmt.Println(dockerContainerId)

	w.Header().Set("Content-Type", "application/json")
	response := map[string]interface{}{
		"thread_id": threadID,
	}

	if threads[threadID] == nil {
		threads[threadID] = make([]*ThreadMessage, 0)
	}

	threads[threadID] = append(threads[threadID], &ThreadMessage{
		ID:        time.Now().String(),
		Direction: "input",
		Output:    req.Prompt,
		CreatedAt: time.Now(),
	})

	json.NewEncoder(w).Encode(response)

	// Create new thread output entry
	//outputMutex.Lock()
	//threads[threadID] = &ThreadMessage{
	//	ID:        threadID,
	//	Status:    "processing",
	//	CreatedAt: time.Now(),
	//}
	//outputMutex.Unlock()

	// Add thread ID to response
	//response := map[string]string{
	//	"status":    "success",
	//	"message":   "Request processed successfully",
	//	"thread_id": threadID,
	//}
	//
	//// Respond to client immediately with thread ID
	//w.Header().Set("Content-Type", "application/json")
	//json.NewEncoder(w).Encode(response)
	//
	//fmt.Printf("===== Starting thread %s processing =====\n", threadID)
	//os.Stdout.Sync()
	//// Process Docker execution in a goroutine
	//go func() {
	//	// Set up Docker run command
	//	output, err := startDockerContainer(threadID, req.RepositoryLink, req.ContextFiles, req.Prompt, req.DockerImage)
	//
	//	// Update thread output in memory
	//	outputMutex.Lock()
	//	defer outputMutex.Unlock()
	//
	//	// Check if thread output still exists (might have been cleaned up)
	//	threadOutput, exists := threads[threadID]
	//	if !exists {
	//		fmt.Printf("Warning: Thread %s was cleaned up before processing completed\n", threadID)
	//		return
	//	}
	//
	//	if err != nil {
	//		threadOutput.Status = "error"
	//		threadOutput.Error = err.Error()
	//		fmt.Printf("Error running Docker container for thread %s: %v\n", threadID, err)
	//	} else {
	//		threadOutput.Status = "completed"
	//		threadOutput.Output = output
	//		// Log output to console for now
	//		fmt.Printf("Thread %s completed with output:\n%s\n", threadID, output)
	//	}
	//}()
}

// handleOutputRequest retrieves output for a specific thread ID
func handleOutputRequest(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Get thread ID from query parameter
	threadID := r.URL.Query().Get("thread_id")
	if threadID == "" {
		http.Error(w, "Missing thread_id parameter", http.StatusBadRequest)
		return
	}

	// Lock for thread-safe access to output map
	outputMutex.Lock()
	threadOutput, exists := threads[threadID]
	outputMutex.Unlock()

	if !exists {
		// Thread ID not found or processing not completed yet
		http.Error(w, "Output not found for thread ID", http.StatusNotFound)
		return
	}

	// Respond with the stored output
	w.Header().Set("Content-Type", "application/json")
	b, err := json.Marshal(threadOutput)
	if err != nil {
		http.Error(w, "Error marshaling JSON", http.StatusInternalServerError)
		return
	}
	response := map[string]string{
		"thread_id": threadID,
		"thread":    string(b),
	}

	// Add output or error depending on status
	//if threadOutput.Status == "completed" {
	//	response["output"] = threadOutput.Output
	//} else if threadOutput.Status == "error" {
	//	response["error"] = threadOutput.Error
	//}

	json.NewEncoder(w).Encode(response)
}

// startDockerContainer handles running a command in a Docker container and captures output
// handleThreadsRequest returns all active thread IDs
func handleThreadsRequest(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Lock for thread-safe access to output map
	outputMutex.Lock()

	// Get all thread IDs
	threadIDs := make([]string, 0, len(threads))
	threadData := make([]map[string]interface{}, 0, len(threads))

	for id, messages := range threads {
		threadIDs = append(threadIDs, id)
		if len(messages) > 0 {
			threadData = append(threadData, map[string]interface{}{
				"thread_id":  id,
				"status":     messages[len(messages)-1].Status,
				"created_at": time.Now(),
			})
		}
	}
	outputMutex.Unlock()

	w.Header().Set("Content-Type", "application/json")

	// Return the list of thread IDs
	response := map[string]interface{}{
		"thread_ids": threadIDs,
		"threads":    threadData,
	}

	json.NewEncoder(w).Encode(response)
}

func startDockerContainer(threadID, repoLink string, contextFiles [][]byte, dockerImage, serverUrl string) (string, error) {
	// Create a buffer to store the output
	var output bytes.Buffer

	// Create temporary directory for this execution
	tempDir, err := os.MkdirTemp("", "superdev-"+threadID)
	if err != nil {
		return "", fmt.Errorf("failed to create temp directory: %w", err)
	}
	//defer os.RemoveAll(tempDir) // Clean up after execution

	// Create repo directory for volume mounting
	repoDir := tempDir + "/repo"
	if err := os.Mkdir(repoDir, 0755); err != nil {
		return "", fmt.Errorf("failed to create repo directory: %w", err)
	}

	// Create context directory for context files
	contextDir := tempDir + "/context"
	if err := os.Mkdir(contextDir, 0755); err != nil {
		return "", fmt.Errorf("failed to create guidance directory: %w", err)
	}

	// Write context files to context directory
	for i, fileContent := range contextFiles {
		filePath := fmt.Sprintf("%s/context_%d.txt", contextDir, i)
		if err := os.WriteFile(filePath, fileContent, 0644); err != nil {
			return "", fmt.Errorf("failed to write context file %d: %w", i, err)
		}
	}

	// Log that we're using a pre-built Docker image
	fmt.Printf("===== Using Docker image %s for thread %s =====\n", dockerImage, threadID)
	os.Stdout.Sync()

	// Clone repository
	fmt.Printf("===== Cloning repository for thread %s =====\n", threadID)
	os.Stdout.Sync()

	cloneCmd := exec.Command("git", "clone", repoLink, repoDir)

	// Set up pipes for real-time output
	cloneStdoutPipe, err := cloneCmd.StdoutPipe()
	if err != nil {
		return "", fmt.Errorf("failed to create stdout pipe for git clone: %w", err)
	}
	cloneStderrPipe, err := cloneCmd.StderrPipe()
	if err != nil {
		return "", fmt.Errorf("failed to create stderr pipe for git clone: %w", err)
	}

	// Log the command being executed
	fmt.Printf("Executing git clone: %s\n", strings.Join(cloneCmd.Args, " "))
	os.Stdout.Sync()

	// Start the command
	if err := cloneCmd.Start(); err != nil {
		return "", fmt.Errorf("failed to start git clone: %w", err)
	}

	// Capture clone output
	var cloneOutput bytes.Buffer

	// Process stdout
	go func() {
		scanner := bufio.NewScanner(cloneStdoutPipe)
		for scanner.Scan() {
			line := scanner.Text()
			fmt.Println("[CLONE]", line)
			os.Stdout.Sync()
			cloneOutput.WriteString(line + "\n")
		}
	}()

	// Process stderr
	go func() {
		scanner := bufio.NewScanner(cloneStderrPipe)
		for scanner.Scan() {
			line := scanner.Text()
			fmt.Println("[CLONE ERR]", line)
			os.Stdout.Sync()
			cloneOutput.WriteString(line + "\n")
		}
	}()

	// Wait for clone to complete
	if err := cloneCmd.Wait(); err != nil {
		fmt.Printf("===== Git clone failed for thread %s =====\n", threadID)
		return "", fmt.Errorf("failed to clone repository: %w, output: %s", err, cloneOutput.String())
	}

	fmt.Printf("===== Repository cloned successfully for thread %s =====\n", threadID)
	os.Stdout.Sync()

	// Pull latest from main branch
	fmt.Printf("===== Pulling latest changes for thread %s =====\n", threadID)
	os.Stdout.Sync()

	pullCmd := exec.Command("git", "pull", "origin", "main")
	pullCmd.Dir = repoDir

	// Set up pipes for real-time output
	pullStdoutPipe, err := pullCmd.StdoutPipe()
	if err != nil {
		return "", fmt.Errorf("failed to create stdout pipe for git pull: %w", err)
	}
	pullStderrPipe, err := pullCmd.StderrPipe()
	if err != nil {
		return "", fmt.Errorf("failed to create stderr pipe for git pull: %w", err)
	}

	// Log the command being executed
	fmt.Printf("Executing git pull: %s (in %s)\n", strings.Join(pullCmd.Args, " "), repoDir)
	os.Stdout.Sync()

	// Start the command
	if err := pullCmd.Start(); err != nil {
		return "", fmt.Errorf("failed to start git pull: %w", err)
	}

	// Capture pull output
	var pullOutput bytes.Buffer

	// Process stdout
	go func() {
		scanner := bufio.NewScanner(pullStdoutPipe)
		for scanner.Scan() {
			line := scanner.Text()
			fmt.Println("[PULL]", line)
			os.Stdout.Sync()
			pullOutput.WriteString(line + "\n")
		}
	}()

	// Process stderr
	go func() {
		scanner := bufio.NewScanner(pullStderrPipe)
		for scanner.Scan() {
			line := scanner.Text()
			fmt.Println("[PULL ERR]", line)
			os.Stdout.Sync()
			pullOutput.WriteString(line + "\n")
		}
	}()

	// Wait for pull to complete
	if err := pullCmd.Wait(); err != nil {
		fmt.Printf("===== Git pull failed for thread %s =====\n", threadID)
		return "", fmt.Errorf("failed to pull from main branch: %w, output: %s", err, pullOutput.String())
	}

	fmt.Printf("===== Repository pull completed for thread %s =====\n", threadID)
	os.Stdout.Sync()

	// Get ANTHROPIC_API_KEY from environment
	anthropicKey := os.Getenv("ANTHROPIC_API_KEY")

	// Prepare Docker run command
	dockerArgs := []string{
		"run",
		"--rm",
		"-d",
		"-e SERVER_URL=" + serverUrl,
		"-e THREAD_ID=" + threadID,
		"-v", repoDir + ":/workdir/repo",
		"-v", contextDir + ":/workdir/context",
	}

	// Add ANTHROPIC_API_KEY as environment variable if available
	if anthropicKey != "" {
		dockerArgs = append(dockerArgs, "-e", "ANTHROPIC_API_KEY="+anthropicKey)
	}

	// Add image
	dockerArgs = append(dockerArgs, dockerImage)

	// Create command
	runCmd := exec.Command("docker", dockerArgs...)

	// Log the command being executed directly to stdout for visibility
	fmt.Printf("Executing Docker command: %s\n", strings.Join(runCmd.Args, " "))
	os.Stdout.Sync()

	// Set up pipes for real-time streaming
	stdoutPipe, err := runCmd.StdoutPipe()
	if err != nil {
		return "", fmt.Errorf("failed to create stdout pipe: %w", err)
	}
	stderrPipe, err := runCmd.StderrPipe()
	if err != nil {
		return "", fmt.Errorf("failed to create stderr pipe: %w", err)
	}

	// Start the command
	if err := runCmd.Start(); err != nil {
		return "", fmt.Errorf("failed to start docker command: %w", err)
	}

	fmt.Printf("===== Docker output for thread %s BEGIN =====\n", threadID)
	os.Stdout.Sync()

	// Capture stdout
	go func() {
		scanner := bufio.NewScanner(stdoutPipe)
		for scanner.Scan() {
			line := scanner.Text()
			fmt.Println(line)
			os.Stdout.Sync()
			output.WriteString(line + "\n")
		}
	}()

	// Capture stderr
	go func() {
		scanner := bufio.NewScanner(stderrPipe)
		for scanner.Scan() {
			line := scanner.Text()
			fmt.Println(line)
			os.Stdout.Sync()
			output.WriteString(line + "\n")
		}
	}()

	// Wait for command to complete
	if err := runCmd.Wait(); err != nil {
		fmt.Printf("===== Docker output for thread %s END (with error) =====\n", threadID)
		os.Stdout.Sync()
		return output.String(), fmt.Errorf("error running Docker container: %w", err)
	}

	fmt.Printf("===== Docker output for thread %s END (success) =====\n", threadID)
	os.Stdout.Sync()

	return output.String(), nil
}
