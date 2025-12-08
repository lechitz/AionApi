// Package dto contains data transfer objects for chat HTTP requests/responses.
package dto

// ChatRequest represents the incoming chat message from the client.
type ChatRequest struct {
	Message string `json:"message" validate:"required,min=1,max=2000" example:"Quanto de água eu bebi hoje?"`
}

// ChatResponse represents the response returned to the client.
type ChatResponse struct {
	Response string                   `json:"response"          example:"Você bebeu 2.5 litros de água hoje..."`
	Sources  []map[string]interface{} `json:"sources,omitempty"`
	Usage    *TokenUsage              `json:"usage,omitempty"`
}

// TokenUsage represents LLM token consumption statistics.
type TokenUsage struct {
	PromptTokens     int `json:"prompt_tokens,omitempty"     example:"50"`
	CompletionTokens int `json:"completion_tokens,omitempty" example:"100"`
	TotalTokens      int `json:"total_tokens"                example:"150"`
}

// InternalChatRequest represents the request sent to the Aion-Chat service (Python).
type InternalChatRequest struct {
	UserID  uint64                 `json:"user_id"`
	Message string                 `json:"message"`
	Context map[string]interface{} `json:"context,omitempty"`
}

// InternalChatResponse represents the response from the Aion-Chat service.
type InternalChatResponse struct {
	Response      string                   `json:"response"`
	FunctionCalls []FunctionCall           `json:"function_calls,omitempty"`
	TokensUsed    int                      `json:"tokens_used,omitempty"`
	Sources       []map[string]interface{} `json:"sources,omitempty"`
}

// FunctionCall represents a function that was called by the LLM.
type FunctionCall struct {
	Name string                 `json:"name"`
	Args map[string]interface{} `json:"args"`
}
