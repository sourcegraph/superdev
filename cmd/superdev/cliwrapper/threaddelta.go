package superdev

// ThreadDeltaType defines the possible types of thread deltas
type ThreadDeltaType string

// Constants for ThreadDeltaType
const (
	ThreadDeltaUserMessage       ThreadDeltaType = "user:message"
	ThreadDeltaUserToolInput     ThreadDeltaType = "user:tool-input"
	ThreadDeltaToolData          ThreadDeltaType = "tool:data"
	ThreadDeltaAssistantMessage  ThreadDeltaType = "assistant:message"
	ThreadDeltaAssistantDelta    ThreadDeltaType = "assistant:message-delta"
	ThreadDeltaInferenceComplete ThreadDeltaType = "inference:completed"
	ThreadDeltaCancelled         ThreadDeltaType = "cancelled"
	ThreadDeltaTitle             ThreadDeltaType = "title"
	ThreadDeltaEnvironment       ThreadDeltaType = "environment"
	ThreadDeltaSummaryCreated    ThreadDeltaType = "summary:created"
)

// ToolRunUserInput represents user input for a tool
type ToolRunUserInput struct {
	Value string `json:"value,omitempty"`
}

// ThreadToolUseID is the ID of a tool use
type ThreadToolUseID string

// ThreadID represents a thread identifier
type ThreadID string

// ThreadEnvironment represents the environment for a thread
type ThreadEnvironment struct {
	Interactive bool   `json:"interactive,omitempty"`
	Platform    string `json:"platform,omitempty"`
}

// DebugUsage contains usage information from inference
type DebugUsage struct {
	PromptTokens     int `json:"promptTokens,omitempty"`
	CompletionTokens int `json:"completionTokens,omitempty"`
	TotalTokens      int `json:"totalTokens,omitempty"`
}

// ToolRun represents data from a tool invocation
type ToolRun struct {
	Status string      `json:"status,omitempty"`
	Result interface{} `json:"result,omitempty"`
}

// ThreadUserMessage represents a user message in a thread
type ThreadUserMessage struct {
	Role    string                   `json:"role,omitempty"`
	Content []map[string]interface{} `json:"content,omitempty"`
}

// ThreadAssistantMessage represents an assistant message in a thread
type ThreadAssistantMessage struct {
	Role    string                   `json:"role,omitempty"`
	Content []map[string]interface{} `json:"content,omitempty"`
}

// AnthropicMessageStreamEvent represents a message stream event from Anthropic
type AnthropicMessageStreamEvent struct {
	Type    string      `json:"type,omitempty"`
	Message interface{} `json:"message,omitempty"`
	Delta   interface{} `json:"delta,omitempty"`
}

// ThreadDelta represents an atomic update to a thread
type ThreadDelta struct {
	// Common fields
	Type ThreadDeltaType `json:"type"`

	// For user:message
	Index   *int               `json:"index,omitempty"`
	Message *ThreadUserMessage `json:"message,omitempty"`

	// For user:tool-input
	ToolUse ThreadToolUseID   `json:"toolUse,omitempty"`
	Value   *ToolRunUserInput `json:"value,omitempty"`

	// For tool:data
	Data *ToolRun `json:"data,omitempty"`

	// For assistant:message
	AssistantMessage *ThreadAssistantMessage `json:"assistantMessage,omitempty"`

	// For assistant:message-delta
	Event *AnthropicMessageStreamEvent `json:"event,omitempty"`

	// For inference:completed
	Usage  *DebugUsage `json:"usage,omitempty"`
	Params interface{} `json:"params,omitempty"`

	// For title
	Title *string `json:"title,omitempty"`

	// For environment
	Env *struct {
		Initial ThreadEnvironment `json:"initial"`
	} `json:"env,omitempty"`

	// For summary:created
	SummaryThreadID ThreadID `json:"summaryThreadID,omitempty"`
}
