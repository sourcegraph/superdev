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

// ThreadOutput stores the output for each thread
type ThreadOutput struct {
	ID        string
	Output    string
	Status    string    // "processing", "completed", or "error"
	CreatedAt time.Time // For cleanup purposes
	Error     string    // Error message if status is "error"
}

// In-memory storage for thread outputs
var (
	threadOutputs = make(map[string]*ThreadOutput)
	outputMutex   = &sync.Mutex{}
	maxOutputAge  = 24 * time.Hour // Outputs older than this will be cleaned up
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
		http.HandleFunc("/run", corsMiddleware(handleRunRequest))
		http.HandleFunc("/output", corsMiddleware(handleOutputRequest))
		http.HandleFunc("/threads", corsMiddleware(handleThreadsRequest))

		fmt.Printf("Server started on :%s\n", port)
		err := http.ListenAndServe(":"+port, nil)
		if err != nil {
			fmt.Printf("Error starting server: %v\n", err)
			os.Exit(1)
		}
	},
}

// handleRunRequest processes Docker build requests
func handleRunRequest(w http.ResponseWriter, r *http.Request) {
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
		Prompt         string   `json:"prompt,omitempty"`
		DockerImage    string   `json:"docker_image,omitempty"`
	}

	if err := json.Unmarshal(body, &req); err != nil {
		http.Error(w, "Error parsing JSON", http.StatusBadRequest)
		return
	}

	// Validate required fields
	if req.DockerImage == "" {
		// Default to the standard image if not specified
		req.DockerImage = "superdev-wrapped-image"
	}

	if req.RepositoryLink == "" {
		http.Error(w, "Repository link is required", http.StatusBadRequest)
		return
	}

	// Process the request
	fmt.Printf("Received request: Docker image: %s, Repo: %s, Context files count: %d\n",
		req.DockerImage, req.RepositoryLink, len(req.ContextFiles))

	if req.Prompt != "" {
		fmt.Printf("Prompt: %s\n", req.Prompt)
	}

	// We'll respond with thread ID later

	// Execute Docker command and store the output
	// Generate a unique thread ID
	threadID, err := generateThreadID()
	if err != nil {
		http.Error(w, "Error generating thread ID", http.StatusInternalServerError)
		return
	}

	// Create new thread output entry
	outputMutex.Lock()
	threadOutputs[threadID] = &ThreadOutput{
		ID:        threadID,
		Status:    "processing",
		CreatedAt: time.Now(),
	}
	outputMutex.Unlock()

	// Add thread ID to response
	response := map[string]string{
		"status":    "success",
		"message":   "Request processed successfully",
		"thread_id": threadID,
	}

	// Respond to client immediately with thread ID
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)

	fmt.Printf("===== Starting thread %s processing =====\n", threadID)
	os.Stdout.Sync()
	// Process Docker execution in a goroutine
	// Set up Docker run command
	outputChan, err := runDockerContainer(threadID, req.RepositoryLink, req.ContextFiles, req.Prompt, req.DockerImage)
	if err != nil {
		// Update thread output in memory with error
		outputMutex.Lock()
		threadOutput, exists := threadOutputs[threadID]
		if exists {
			threadOutput.Status = "error"
			threadOutput.Error = err.Error()
		}
		outputMutex.Unlock()
		fmt.Printf("Error setting up Docker container for thread %s: %v\n", threadID, err)
		return
	}

	go func() {
		// Process messages from channel
		for msg := range outputChan {
			// Check if it's an error message
			if strings.HasPrefix(msg, "ERROR: ") {
				outputMutex.Lock()
				threadOutput, exists := threadOutputs[threadID]
				if exists {
					threadOutput.Status = "error"
					threadOutput.Error = strings.TrimPrefix(msg, "ERROR: ")
					// Append error to output as well
					threadOutput.Output += "\nERROR: " + threadOutput.Error
				}
				outputMutex.Unlock()
				continue
			}

			// Update thread output in memory with new message
			outputMutex.Lock()
			threadOutput, exists := threadOutputs[threadID]
			if exists {
				threadOutput.Output += msg + "\n"
			}
			outputMutex.Unlock()
		}

		// When channel is closed, mark as completed
		outputMutex.Lock()
		threadOutput, exists := threadOutputs[threadID]
		if exists && threadOutput.Status != "error" {
			threadOutput.Status = "completed"
			fmt.Printf("Thread %s completed\n", threadID)
		}
		outputMutex.Unlock()
	}()
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
	threadOutput, exists := threadOutputs[threadID]
	outputMutex.Unlock()

	if !exists {
		// Thread ID not found or processing not completed yet
		http.Error(w, "Output not found for thread ID", http.StatusNotFound)
		return
	}

	// Respond with the stored output
	w.Header().Set("Content-Type", "application/json")
	response := map[string]string{
		"thread_id": threadID,
		"status":    threadOutput.Status,
	}

	// Add output or error depending on status
	if threadOutput.Status != "error" {
		response["output"] = threadOutput.Output
	} else if threadOutput.Status == "error" {
		response["error"] = threadOutput.Error
	}

	json.NewEncoder(w).Encode(response)
}

// runDockerContainer handles running a command in a Docker container and captures output
// handleThreadsRequest returns all active thread IDs
func handleThreadsRequest(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Lock for thread-safe access to output map
	outputMutex.Lock()

	// Get all thread IDs
	threadIDs := make([]string, 0, len(threadOutputs))
	threadData := make([]map[string]interface{}, 0, len(threadOutputs))

	for id, output := range threadOutputs {
		threadIDs = append(threadIDs, id)
		threadData = append(threadData, map[string]interface{}{
			"thread_id":  id,
			"status":     output.Status,
			"created_at": output.CreatedAt,
		})
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

func runDockerContainer(threadID, repoLink string, contextFiles [][]byte, prompt string, dockerImage string) (chan string, error) {
	// Create a buffer to store the output
	var output bytes.Buffer

	// Create output channel
	outputChan := make(chan string)

	// Create temporary directory for this execution
	tempDir, err := os.MkdirTemp("", "superdev-"+threadID)
	if err != nil {
		return nil, fmt.Errorf("failed to create temp directory: %w", err)
	}
	defer os.RemoveAll(tempDir) // Clean up after execution

	// Create repo directory for volume mounting
	repoDir := tempDir + "/repo"
	if err := os.Mkdir(repoDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create repo directory: %w", err)
	}

	// Create guidance directory for context files
	guidanceDir := tempDir + "/guidance"
	if err := os.Mkdir(guidanceDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create guidance directory: %w", err)
	}

	// Write context files to guidance directory
	for i, fileContent := range contextFiles {
		filePath := fmt.Sprintf("%s/context_%d.txt", guidanceDir, i)
		if err := os.WriteFile(filePath, fileContent, 0644); err != nil {
			return nil, fmt.Errorf("failed to write context file %d: %w", i, err)
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
		return nil, fmt.Errorf("failed to create stdout pipe for git clone: %w", err)
	}
	cloneStderrPipe, err := cloneCmd.StderrPipe()
	if err != nil {
		return nil, fmt.Errorf("failed to create stderr pipe for git clone: %w", err)
	}

	// Log the command being executed
	fmt.Printf("Executing git clone: %s\n", strings.Join(cloneCmd.Args, " "))
	os.Stdout.Sync()

	// Start the command
	if err := cloneCmd.Start(); err != nil {
		return nil, fmt.Errorf("failed to start git clone: %w", err)
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
		return nil, fmt.Errorf("failed to clone repository: %w, output: %s", err, cloneOutput.String())
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
		return nil, fmt.Errorf("failed to create stdout pipe for git pull: %w", err)
	}
	pullStderrPipe, err := pullCmd.StderrPipe()
	if err != nil {
		return nil, fmt.Errorf("failed to create stderr pipe for git pull: %w", err)
	}

	// Log the command being executed
	fmt.Printf("Executing git pull: %s (in %s)\n", strings.Join(pullCmd.Args, " "), repoDir)
	os.Stdout.Sync()

	// Start the command
	if err := pullCmd.Start(); err != nil {
		return nil, fmt.Errorf("failed to start git pull: %w", err)
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
		return nil, fmt.Errorf("failed to pull from main branch: %w, output: %s", err, pullOutput.String())
	}

	fmt.Printf("===== Repository pull completed for thread %s =====\n", threadID)
	os.Stdout.Sync()

	// Get ANTHROPIC_API_KEY from environment
	anthropicKey := os.Getenv("ANTHROPIC_API_KEY")

	// Prepare Docker run command
	dockerArgs := []string{
		"run",
		"--rm",
		"-v", repoDir + ":/workdir/repo",
		"-v", guidanceDir + ":/workdir/guidance",
	}

	// Add ANTHROPIC_API_KEY as environment variable if available
	if anthropicKey != "" {
		dockerArgs = append(dockerArgs, "-e", "ANTHROPIC_API_KEY="+anthropicKey)
	}

	// Add image and command
	dockerArgs = append(dockerArgs,
		dockerImage,
		"sh", "-c", "echo '"+prompt+"' | amp")

	// Create command
	runCmd := exec.Command("docker", dockerArgs...)

	// Log the command being executed directly to stdout for visibility
	fmt.Printf("Executing Docker command: %s\n", strings.Join(runCmd.Args, " "))
	os.Stdout.Sync()

	// Set up pipes for real-time streaming
	stdoutPipe, err := runCmd.StdoutPipe()
	if err != nil {
		return nil, fmt.Errorf("failed to create stdout pipe: %w", err)
	}
	stderrPipe, err := runCmd.StderrPipe()
	if err != nil {
		return nil, fmt.Errorf("failed to create stderr pipe: %w", err)
	}

	// Start the command
	if err := runCmd.Start(); err != nil {
		return nil, fmt.Errorf("failed to start docker command: %w", err)
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
			outputChan <- line
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
			outputChan <- line
		}
	}()

	// Wait for command to complete in a separate goroutine
	go func() {
		if err := runCmd.Wait(); err != nil {
			fmt.Printf("===== Docker output for thread %s END (with error) =====\n", threadID)
			os.Stdout.Sync()
			// Send error message to channel
			outputChan <- "ERROR: " + err.Error()
		} else {
			fmt.Printf("===== Docker output for thread %s END (success) =====\n", threadID)
			os.Stdout.Sync()
		}
		// Close the channel when done
		close(outputChan)
	}()

	// Return the channel immediately
	return outputChan, nil
}
