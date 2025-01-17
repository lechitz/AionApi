package errors

import (
	"fmt"
	"os"
	"time"
)

// ReportCriticalError saves critical errors to a local file for manual inspection
func ReportCriticalError(err error, message string) {
	file, err := os.OpenFile("/app/logs/error.log", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0666)
	if err != nil {
		fmt.Printf("failed to write critical error: %v\n", err)
		return
	}
	defer file.Close()

	logMessage := fmt.Sprintf("[%s] %s: %v\n", time.Now().Format(time.RFC3339), message, err)
	file.WriteString(logMessage)
}
