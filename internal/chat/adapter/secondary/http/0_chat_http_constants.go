package http

import "time"

const (
	// TracerAionChatClient is the tracer name for Aion-Chat client.
	TracerAionChatClient = "chat.client.aion_chat"
	// SpanSendMessage is the span name for sending a message.
	SpanSendMessage = "SendMessage"
	// DefaultTimeout is the default timeout for HTTP requests.
	DefaultTimeout = 30 * time.Second
)

// HTTP path(s) used by Aion-Chat service.
const (
	// PathProcess is the endpoint for processing chat messages.
	PathProcess = "/internal/process"
)

// HTTP header and content-type constants.
const (
	// HeaderContentType is the HTTP header for Content-Type.
	HeaderContentType = "Content-Type"
	// HeaderAccept is the HTTP header for Accept.
	HeaderAccept = "Accept"
	// ContentTypeJSON is the content type for JSON.
	ContentTypeJSON = "application/json"
)

// Attribute keys used for tracing and span attributes.
const (
	// AttrHTTPURL is the attribute key for the HTTP URL.
	AttrHTTPURL = "http.url"
	// AttrHTTPMethod is the attribute key for the HTTP method.
	AttrHTTPMethod = "http.method"
	// AttrUserID is the attribute key for the user ID.
	AttrUserID = "user_id"
	// AttrHTTPStatusCode is the attribute key for the HTTP status code.
	AttrHTTPStatusCode = "http.status_code"
	// AttrTokensUsed is the attribute key for the number of tokens used.
	AttrTokensUsed = "tokens_used"
	// AttrResponseLength is the attribute key for the response length.
	AttrResponseLength = "response_length"
)

// Log messages used by the client adapter.
const (
	// MsgCallingAionChatService indicates that a call to Aion-Chat service is being made.
	MsgCallingAionChatService = "Calling Aion-Chat service"
	// MsgAionChatResponseReceived indicates that a response from Aion-Chat service has been received.
	MsgAionChatResponseReceived = "Aion-Chat response received"
)

// Error messages used for tracing, logging and httpresponse wrapping.
const (
	// ErrFailedMarshal indicates a failure to marshal the request.
	ErrFailedMarshal = "failed to marshal request"
	// ErrFailedCreateRequest indicates a failure to create the HTTP request.
	ErrFailedCreateRequest = "failed to create request"
	// ErrHTTPRequestFailed indicates that the HTTP request failed.
	ErrHTTPRequestFailed = "http request failed"
	// ErrFailedReadResponse indicates a failure to read the HTTP response.
	ErrFailedReadResponse = "failed to read response"
	// ErrAionChatNonOK indicates that Aion-Chat returned a non-200 status code.
	ErrAionChatNonOK = "aion-chat returned non-200 status"
	// ErrFailedUnmarshal indicates a failure to unmarshal the response.
	ErrFailedUnmarshal = "failed to unmarshal response"
	// ErrAionChatRequestFailed indicates that the request to Aion-Chat service failed.
	ErrAionChatRequestFailed = "aion-chat request failed"
)

// Status names for semantic span states.
const (
	// StatusMessageSent indicates that the message was sent successfully.
	StatusMessageSent = "message_sent"
)
