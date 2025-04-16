package superdev

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"text/template"

	"github.com/spf13/cobra"
)

// Command variables
var (
	rootCmd = &cobra.Command{
		Use:   "superdev",
		Short: "SuperDev is a simple CLI tool",
	}
)

// GetRootCmd returns the root command
func GetRootCmd() *cobra.Command {
	return rootCmd
}

// InitCommands initializes all CLI commands
func InitCommands() {
	// Add flags to run command
	runCmd.Flags().StringVar(&serverURL, "server", "http://localhost:8080", "Server URL to send the Docker image to")
	runCmd.Flags().StringVar(&prompt, "prompt", "Hello from the CLI", "Prompt to send to the server")

	rootCmd.AddCommand(runCmd)
	rootCmd.AddCommand(serverCmd)
}

// sendImageToServer sends a request to the server with the Docker image and prompt
func sendImageToServer(serverURL, dockerImage, prompt string) (string, error) {
	// Create request body
	requestBody := map[string]string{
		"docker_image":    dockerImage,
		"repository_link": "https://github.com/sourcegraph/amp.git", // Hardcoded for now
		"prompt":          prompt,
	}

	// Convert request body to JSON
	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		return "", fmt.Errorf("failed to marshal JSON: %w", err)
	}

	// Create HTTP request
	req, err := http.NewRequest("POST", serverURL+"/start", bytes.NewBuffer(jsonData))
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	// Set Content-Type header
	req.Header.Set("Content-Type", "application/json")

	// Send request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	// Read response body
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response: %w", err)
	}

	// Check status code
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("server returned non-OK status: %d, body: %s", resp.StatusCode, string(respBody))
	}

	// Parse response JSON
	var responseData struct {
		Status   string `json:"status"`
		Message  string `json:"message"`
		ThreadID string `json:"thread_id"`
	}

	err = json.Unmarshal(respBody, &responseData)
	if err != nil {
		return "", fmt.Errorf("failed to parse response JSON: %w", err)
	}

	// Return thread ID
	return responseData.ThreadID, nil
}

var (
	serverURL string
	prompt    string
)

var runCmd = &cobra.Command{
	Use:   "run [dockerfile]",
	Short: "Build a Docker image from the specified Dockerfile and send it to the server",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		dockerfilePath := args[0]

		// Check if file exists
		if _, err := os.Stat(dockerfilePath); os.IsNotExist(err) {
			fmt.Printf("Error: Dockerfile not found at %s\n", dockerfilePath)
			os.Exit(1)
		}

		// Build Docker image
		fmt.Printf("Building Docker image from %s...\n", dockerfilePath)

		// Execute docker build command
		dockerCmd := exec.Command("docker", "build", "-f", dockerfilePath, "-t", "superdev-image", ".")
		dockerCmd.Stdout = os.Stdout
		dockerCmd.Stderr = os.Stderr

		err := dockerCmd.Run()
		if err != nil {
			fmt.Printf("Error building Docker image: %v\n", err)
			os.Exit(1)
		}

		fmt.Printf("Successfully built Docker image from %s\n", dockerfilePath)

		// Create a wrapper Dockerfile using a template
		wrapperTemplate := `FROM ubuntu:latest as wrapper

# Additional wrapper configuration
LABEL wrapped.by="superdev"
LABEL original.dockerfile="{{.OriginalFile}}"

# Install Node.js 22 and npm
RUN apt-get update && \
    apt-get install -y curl gnupg2 && \
    mkdir -p /etc/apt/keyrings && \
    curl -fsSL https://deb.nodesource.com/gpgkey/nodesource-repo.gpg.key | gpg --dearmor -o /etc/apt/keyrings/nodesource.gpg && \
    echo "deb [signed-by=/etc/apt/keyrings/nodesource.gpg] https://deb.nodesource.com/node_22.x nodistro main" | tee /etc/apt/sources.list.d/nodesource.list && \
    apt-get update && \
    apt-get install -y nodejs && \
    apt-get clean && \
    rm -rf /var/lib/apt/lists/*

# Install pnpm via npm
RUN npm install -g pnpm

# Install ripgrep
RUN apt-get update && \
    apt-get install -y ripgrep && \
    apt-get clean && \
    rm -rf /var/lib/apt/lists/*

# Setup PNPM environment
RUN mkdir -p /usr/local/pnpm-global
ENV PNPM_HOME=/usr/local/pnpm-global
ENV PATH="$PNPM_HOME:$PNPM_HOME/node_modules/.bin:$PATH"
ENV SHELL=/bin/bash

# Install Sourcegraph AMP
RUN pnpm setup && pnpm add -g @sourcegraph/amp

# Create final image with original content included
FROM superdev-image

# Copy all installed tools from wrapper
COPY --from=wrapper /usr/bin /usr/bin
COPY --from=wrapper /usr/local /usr/local
COPY --from=wrapper /usr/lib /usr/lib
COPY --from=wrapper /usr/share /usr/share
COPY --from=wrapper /lib /lib
COPY --from=wrapper /etc/apt /etc/apt

# Set environment variables
ENV PNPM_HOME=/usr/local/pnpm-global
ENV PATH="$PNPM_HOME:$PNPM_HOME/node_modules/.bin:$PATH"
ENV SHELL=/bin/bash

CMD ["echo", "This image was wrapped by SuperDev with all requested tools installed"]
`

		// Create a temporary directory for the wrapper Dockerfile
		tempDir, err := os.MkdirTemp("", "superdev-wrapper")
		if err != nil {
			fmt.Printf("Error creating temp directory: %v\n", err)
			os.Exit(1)
		}
		defer os.RemoveAll(tempDir)

		// Create wrapper Dockerfile
		wrapperPath := filepath.Join(tempDir, "Dockerfile.wrapper")

		// Parse and execute the template
		tmpl, err := template.New("wrapper").Parse(wrapperTemplate)
		if err != nil {
			fmt.Printf("Error parsing template: %v\n", err)
			os.Exit(1)
		}

		// Create the wrapper Dockerfile
		wrapperFile, err := os.Create(wrapperPath)
		if err != nil {
			fmt.Printf("Error creating wrapper Dockerfile: %v\n", err)
			os.Exit(1)
		}
		defer wrapperFile.Close()

		// Execute the template with data
		templateData := struct {
			OriginalFile string
		}{
			OriginalFile: dockerfilePath,
		}

		err = tmpl.Execute(wrapperFile, templateData)
		if err != nil {
			fmt.Printf("Error writing wrapper Dockerfile: %v\n", err)
			os.Exit(1)
		}
		wrapperFile.Close()

		// Build the wrapped image
		fmt.Println("Building wrapped Docker image...")
		wrapperCmd := exec.Command("docker", "build", "-f", wrapperPath, "-t", "superdev-wrapped-image", ".")
		wrapperCmd.Stdout = os.Stdout
		wrapperCmd.Stderr = os.Stderr

		err = wrapperCmd.Run()
		if err != nil {
			fmt.Printf("Error building wrapped Docker image: %v\n", err)
			os.Exit(1)
		}

		fmt.Printf("Successfully built wrapped Docker image from %s\n", dockerfilePath)

		// Use the server URL and prompt provided via flags

		fmt.Printf("Sending wrapped Docker image to the server...\n")
		threadID, err := sendImageToServer(serverURL, "superdev-wrapped-image", prompt)
		if err != nil {
			fmt.Printf("Error sending image to server: %v\n", err)
			os.Exit(1)
		}

		fmt.Printf("Successfully sent image to server. Thread ID: %s\n", threadID)
		fmt.Printf("To check status, use: curl -X GET \"http://localhost:8080/output?thread_id=%s\"\n", threadID)
	},
}
