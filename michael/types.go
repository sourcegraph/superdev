package main

import (
	"encoding/json"
	"fmt"
	"strings"
)

// AmpItem is an interface for any item that can be rendered
type AmpItem interface {
	Render() string
}

// AmpWorkerResponse represents the main response structure from the worker
type AmpWorkerResponse struct {
	StreamID    int         `json:"streamId,omitempty"`
	StreamEvent string      `json:"streamEvent,omitempty"`
	Data        interface{} `json:"data,omitempty"`
}

// Thread information
type AmpThread struct {
	Created        int64           `json:"created"`
	ID             string          `json:"id"`
	Title          string          `json:"title,omitempty"`
	Messages       []AmpMessage    `json:"messages"`
	Env            *AmpEnv         `json:"env,omitempty"`
	V              int             `json:"v"`
	FileChanges    *AmpFileChanges `json:"fileChanges,omitempty"`
	State          string          `json:"state,omitempty"`
	InferenceState string          `json:"inferenceState,omitempty"`
}

// AmpGenericItem allows unknown types to implement AmpItem
type AmpGenericItem struct {
	Data interface{}
}

// Implement the AmpItem interface for our types
func (t AmpThread) Render() string {
	var sb strings.Builder

	sb.WriteString(fmt.Sprintf("Thread [ID: %s] (v%d)\n", t.ID, t.V))
	if t.Title != "" {
		sb.WriteString(fmt.Sprintf("Title: %s\n", t.Title))
	}

	// Show messages in a more readable format
	if len(t.Messages) > 0 {
		sb.WriteString("\nMessages:\n")
		for i, msg := range t.Messages {
			sb.WriteString(fmt.Sprintf("  [%d] %s: ", i+1, msg.Role))

			if len(msg.Content) > 0 {
				for _, content := range msg.Content {
					switch content.Type {
					case "text":
						if content.Text != "" {
							sb.WriteString(fmt.Sprintf("%s\n", content.Text))
						}
					case "thinking":
						if content.Thinking != "" {
							sb.WriteString(fmt.Sprintf("(thinking: %s)\n", summarizeThinking(content.Thinking)))
						}
					case "tool_use":
						sb.WriteString(fmt.Sprintf("(using tool: %s)\n", content.Name))
					default:
						sb.WriteString(fmt.Sprintf("(%s content)\n", content.Type))
					}
				}
			} else {
				sb.WriteString("(empty)\n")
			}
		}
	}

	// Include the state if present
	if t.State != "" || t.InferenceState != "" {
		sb.WriteString(fmt.Sprintf("\nState: %s, Inference: %s\n", t.State, t.InferenceState))
	}

	// Include file changes if present
	if t.FileChanges != nil && len(t.FileChanges.Files) > 0 {
		sb.WriteString("\nFile Changes:\n")
		for _, file := range t.FileChanges.Files {
			sb.WriteString(fmt.Sprintf("  - %s\n", file.Path))
		}
	}

	return sb.String()
}

func (g AmpGenericItem) Render() string {
	bytes, _ := json.MarshalIndent(g.Data, "", "  ")
	return fmt.Sprintf("Unknown: %s", string(bytes))
}

// Helper function to summarize long thinking content
func summarizeThinking(thinking string) string {
	if len(thinking) <= 60 {
		return thinking
	}
	return thinking[:57] + "..."
}

// SerializeMessages converts an array of AmpMessages into a single string representation
// of the conversation history
func SerializeMessages(messages []AmpMessage) string {
	var sb strings.Builder

	for i, msg := range messages {
		// Add a separator between messages
		if i > 0 {
			sb.WriteString("\n---\n")
		}

		// Add the role
		sb.WriteString(fmt.Sprintf("%s: ", msg.Role))

		// Add the content
		for j, content := range msg.Content {
			if j > 0 {
				sb.WriteString("\n")
			}

			switch content.Type {
			case "text":
				sb.WriteString(content.Text)
			case "thinking":
				sb.WriteString(fmt.Sprintf("(thinking: %s)", content.Thinking))
			case "tool_use":
				toolData, _ := json.Marshal(content.Input)
				sb.WriteString(fmt.Sprintf("(using tool: %s with input: %s)", content.Name, string(toolData)))
			default:
				sb.WriteString(fmt.Sprintf("(%s content)", content.Type))
			}
		}
	}

	return sb.String()
}

// Environment details
type AmpEnv struct {
	Initial AmpEnvInitial `json:"initial"`
}

type AmpEnvInitial struct {
	Interactive bool        `json:"interactive"`
	Platform    AmpPlatform `json:"platform"`
	Tags        []string    `json:"tags"`
	Trees       []AmpTree   `json:"trees"`
}

type AmpPlatform struct {
	OS         string `json:"os"`
	WebBrowser bool   `json:"webBrowser"`
}

type AmpTree struct {
	DisplayName string        `json:"displayName"`
	Repository  AmpRepository `json:"repository"`
}

type AmpRepository struct {
	Ref  string `json:"ref"`
	SHA  string `json:"sha"`
	Type string `json:"type"`
	URL  string `json:"url"`
}

// Message content
type AmpMessage struct {
	Content []AmpContent `json:"content"`
	Role    string       `json:"role"`
	State   *AmpState    `json:"state,omitempty"`
}

// ToDelta converts an AmpMessage to a ThreadDelta
func (m AmpMessage) ToDelta() ThreadDelta {
	// Convert AmpContent to map[string]interface{}
	contents := make([]map[string]interface{}, len(m.Content))
	
	for i, content := range m.Content {
		contentMap := map[string]interface{}{
			"type": content.Type,
		}
		
		// Add fields based on content type
		if content.Text != "" {
			contentMap["text"] = content.Text
		}
		if content.Thinking != "" {
			contentMap["thinking"] = content.Thinking
		}
		if content.Name != "" {
			contentMap["name"] = content.Name
		}
		if content.Input != nil {
			contentMap["input"] = content.Input
		}
		if content.InputPartialJSON != nil {
			contentMap["inputPartialJSON"] = content.InputPartialJSON
		}
		
		contents[i] = contentMap
	}
	
	// Create ThreadDelta based on role
	if m.Role == "user" {
		return ThreadDelta{
			Type: ThreadDeltaUserMessage,
			Message: &ThreadUserMessage{
				Role:    m.Role,
				Content: contents,
			},
		}
	} else if m.Role == "assistant" {
		return ThreadDelta{
			Type: ThreadDeltaAssistantMessage,
			AssistantMessage: &ThreadAssistantMessage{
				Role:    m.Role,
				Content: contents,
			},
		}
	}
	
	// Default case (should not happen in normal use)
	return ThreadDelta{
		Type: ThreadDeltaUserMessage,
		Message: &ThreadUserMessage{
			Role:    m.Role,
			Content: contents,
		},
	}
}

type AmpContent struct {
	Text             string                 `json:"text,omitempty"`
	Type             string                 `json:"type"`
	Thinking         string                 `json:"thinking,omitempty"`
	Signature        string                 `json:"signature,omitempty"`
	ID               string                 `json:"id,omitempty"`
	Name             string                 `json:"name,omitempty"`
	Input            map[string]interface{} `json:"input,omitempty"`
	InputPartialJSON *AmpPartialJSON        `json:"inputPartialJSON,omitempty"`
}

type AmpPartialJSON struct {
	JSON string `json:"json"`
}

type AmpState struct {
	Type string `json:"type"`
}

// File changes tracking
type AmpFileChanges struct {
	Files []AmpFile `json:"files"`
}

type AmpFile struct {
	Path    string `json:"path"`
	Content string `json:"content,omitempty"`
}

// Thread state
type AmpThreadState struct {
	State          string         `json:"state"`
	InferenceState string         `json:"inferenceState"`
	FileChanges    AmpFileChanges `json:"fileChanges"`
}

func (ts AmpThreadState) Render() string {
	var sb strings.Builder

	sb.WriteString("ThreadState:\n")
	sb.WriteString(fmt.Sprintf("  State: %s\n", ts.State))
	sb.WriteString(fmt.Sprintf("  Inference: %s\n", ts.InferenceState))

	// Include file changes if present
	if len(ts.FileChanges.Files) > 0 {
		sb.WriteString("  File Changes:\n")
		for _, file := range ts.FileChanges.Files {
			sb.WriteString(fmt.Sprintf("    - %s\n", file.Path))
		}
	} else {
		sb.WriteString("  No file changes\n")
	}

	return sb.String()
}
