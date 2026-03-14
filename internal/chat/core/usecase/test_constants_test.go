// Package usecase_test contains test constants for chat usecase tests.
package usecase_test

// Test user messages.
const (
	TestMessageHello           = "Hello, how are you?"
	TestMessageWhatAskedBefore = "What did I ask before?"
	TestMessageFirst           = "First message"
	TestMessageEmpty           = ""
	TestMessageWaterIntake     = "Check my water intake and set a reminder"
)

// Test AI responses.
const (
	TestResponseWellThanks    = "I'm doing well, thank you for asking!"
	TestResponseAIInventor    = "You asked about AI and its inventor"
	TestResponseHelloHelp     = "Hello! How can I help?"
	TestResponseNoMessage     = "I didn't receive a message. Please try again."
	TestResponseWaterReminder = "You drank 1.5L today. I've set a reminder for 3pm."
)

// Test conversation history content.
const (
	TestHistoryQuestionAI       = "What is AI?"
	TestHistoryAnswerAI         = "AI stands for Artificial Intelligence"
	TestHistoryQuestionInventor = "Who invented it?"
	TestHistoryAnswerInventor   = "The term was coined by John McCarthy in 1956"
	TestHistoryQuestionWater    = "What is my water intake?"
	TestHistoryAnswerWater      = "You drank 2.5L today"
	TestHistoryQuestionReminder = "Set a reminder"
	TestHistoryAnswerReminder   = "Reminder set for 3pm"
)

// Test role constants.
const (
	TestRoleUser      = "user"
	TestRoleAssistant = "assistant"
)

// Test function names.
const (
	TestFunctionGetWeather     = "get_weather"
	TestFunctionGetWaterIntake = "get_water_intake"
	TestFunctionSetReminder    = "set_reminder"
)

// Test error messages.
const (
	TestErrorCacheUnavailable   = "cache unavailable"
	TestErrorServiceUnavailable = "aion-chat service unavailable"
	TestErrorDatabaseTimeout    = "database connection timeout"
)

// Test metadata.
const (
	TestSourceType = "knowledge_base"
	TestSourceID   = "kb_123"
)
