package superdev

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/spf13/cobra"
)

// serverCmd handles server operations
var serverCmd = &cobra.Command{
	Use:   "server [port]",
	Short: "Start a server that accepts Docker build requests",
	Args:  cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		port := "8080"
		if len(args) > 0 {
			port = args[0]
		}

		// Run the server command
		fmt.Printf("Starting server on port %s...\n", port)

		// Setup HTTP server
		http.HandleFunc("/run", handleRunRequest)

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
		Dockerfile     string   `json:"dockerfile"`
		RepositoryLink string   `json:"repository_link"`
		ContextFiles   [][]byte `json:"contextFiles,omitempty"`
		Prompt         string   `json:"prompt,omitempty"`
	}

	if err := json.Unmarshal(body, &req); err != nil {
		http.Error(w, "Error parsing JSON", http.StatusBadRequest)
		return
	}

	// Validate required fields
	if req.Dockerfile == "" {
		http.Error(w, "Dockerfile is required", http.StatusBadRequest)
		return
	}

	if req.RepositoryLink == "" {
		http.Error(w, "Repository link is required", http.StatusBadRequest)
		return
	}

	// Process the request
	fmt.Printf("Received request: Dockerfile length: %d, Repo: %s, Context files count: %d\n",
		len(req.Dockerfile), req.RepositoryLink, len(req.ContextFiles))

	if req.Prompt != "" {
		fmt.Printf("Prompt: %s\n", req.Prompt)
	}

	// Respond with success
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"status":  "success",
		"message": "Request processed successfully",
	})

	// put them together and run amp

	// then stream the output to a file on the server or into memory
}
