package superdev

import (
	"encoding/json"
	"testing"
)

func TestThreadDeltaSerialization(t *testing.T) {
	// Test user:message
	userMessage := ThreadDelta{
		Type: ThreadDeltaUserMessage,
		Message: &ThreadUserMessage{
			Content: []map[string]interface{}{
				{
					"type": "text",
					"text": "Hello, world!",
				},
			},
		},
	}

	// Marshal to JSON
	userMessageJSON, err := json.Marshal(userMessage)
	if err != nil {
		t.Fatalf("Failed to marshal user message: %v", err)
	}

	// Unmarshal back to verify
	var decoded ThreadDelta
	if err := json.Unmarshal(userMessageJSON, &decoded); err != nil {
		t.Fatalf("Failed to unmarshal user message: %v", err)
	}

	if decoded.Type != ThreadDeltaUserMessage {
		t.Errorf("Expected type 'user:message', got '%s'", decoded.Type)
	}

	// Test title delta
	titleValue := "Test Thread"
	titleDelta := ThreadDelta{
		Type:  ThreadDeltaTitle,
		Title: &titleValue,
	}

	// Marshal to JSON
	titleJSON, err := json.Marshal(titleDelta)
	if err != nil {
		t.Fatalf("Failed to marshal title delta: %v", err)
	}

	// Unmarshal back to verify
	var decodedTitle ThreadDelta
	if err := json.Unmarshal(titleJSON, &decodedTitle); err != nil {
		t.Fatalf("Failed to unmarshal title delta: %v", err)
	}

	if decodedTitle.Type != ThreadDeltaTitle {
		t.Errorf("Expected type 'title', got '%s'", decodedTitle.Type)
	}

	if *decodedTitle.Title != titleValue {
		t.Errorf("Expected title '%s', got '%s'", titleValue, *decodedTitle.Title)
	}
}

func TestFilterMessages(t *testing.T) {
	// Test messages with thinking content
	messages := []AmpMessage{
		{
			Role: "user",
			Content: []AmpContent{
				{
					Type: "text",
					Text: "Hello",
				},
			},
		},
		{
			Role: "assistant",
			Content: []AmpContent{
				{
					Type: "thinking",
					Thinking: "Thinking about the response",
				},
			},
		},
		{
			Role: "assistant",
			Content: []AmpContent{
				{
					Type: "thinking",
					Thinking: "Still thinking",
				},
				{
					Type: "text",
					Text: "Hello, how can I help?",
				},
			},
		},
	}

	// Filter messages
	filtered := FilterMessages(messages)

	// Expected results
	if len(filtered) != 2 {
		t.Errorf("Expected 2 messages after filtering, got %d", len(filtered))
	}

	// Check first message (user message should be unchanged)
	if filtered[0].Role != "user" || len(filtered[0].Content) != 1 || filtered[0].Content[0].Type != "text" {
		t.Errorf("First message not filtered correctly")
	}

	// Check second message (assistant message with mixed content should only have text content)
	if filtered[1].Role != "assistant" || len(filtered[1].Content) != 1 || filtered[1].Content[0].Type != "text" {
		t.Errorf("Second message not filtered correctly")
	}

	// Check that the message with only thinking content was removed
	for _, msg := range filtered {
		hasOnlyThinking := true
		for _, content := range msg.Content {
			if content.Type != "thinking" {
				hasOnlyThinking = false
				break
			}
		}
		if hasOnlyThinking {
			t.Errorf("Found a message with only thinking content after filtering")
		}
	}
}

func TestAmpMessageToDelta(t *testing.T) {
	tests := []struct {
		name     string
		message  AmpMessage
		expected ThreadDelta
	}{
		{
			name: "User message",
			message: AmpMessage{
				Role: "user",
				Content: []AmpContent{
					{
						Type: "text",
						Text: "Hello, world!",
					},
				},
			},
			expected: ThreadDelta{
				Type: ThreadDeltaUserMessage,
				Message: &ThreadUserMessage{
					Role: "user",
					Content: []map[string]interface{}{
						{
							"type": "text",
							"text": "Hello, world!",
						},
					},
				},
			},
		},
		{
			name: "Assistant message",
			message: AmpMessage{
				Role: "assistant",
				Content: []AmpContent{
					{
						Type: "text",
						Text: "I can help with that.",
					},
				},
			},
			expected: ThreadDelta{
				Type: ThreadDeltaAssistantMessage,
				AssistantMessage: &ThreadAssistantMessage{
					Role: "assistant",
					Content: []map[string]interface{}{
						{
							"type": "text",
							"text": "I can help with that.",
						},
					},
				},
			},
		},
		{
			name: "Tool call message",
			message: AmpMessage{
				Role: "assistant",
				Content: []AmpContent{
					{
						Type: "tool_use",
						Name: "search",
						Input: map[string]interface{}{
							"query": "golang examples",
						},
					},
				},
			},
			expected: ThreadDelta{
				Type: ThreadDeltaAssistantMessage,
				AssistantMessage: &ThreadAssistantMessage{
					Role: "assistant",
					Content: []map[string]interface{}{
						{
							"type": "tool_use",
							"name": "search",
							"input": map[string]interface{}{
								"query": "golang examples",
							},
						},
					},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.message.ToDelta()

			// Serialize both to JSON for comparison
			expectedJSON, err := json.Marshal(tt.expected)
			if err != nil {
				t.Fatalf("Failed to marshal expected: %v", err)
			}

			resultJSON, err := json.Marshal(result)
			if err != nil {
				t.Fatalf("Failed to marshal result: %v", err)
			}

			// Compare JSON strings
			if string(expectedJSON) != string(resultJSON) {
				t.Errorf("\nExpected: %s\nGot: %s", string(expectedJSON), string(resultJSON))
			}
		})
	}
}
