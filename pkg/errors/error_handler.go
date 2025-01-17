package errors

import (
	"go.uber.org/zap"
)

// HandleCriticalError logs the error as fatal and stops the application.
func HandleCriticalError(loggerSugar *zap.SugaredLogger, message string, err error) {
	if err != nil {
		loggerSugar.Fatalw(message, "error", err.Error())
	} else {
		loggerSugar.Fatal(message)
	}
}

// HandleError logs the error but allows the application to continue running.
func HandleError(loggerSugar *zap.SugaredLogger, message string, err error) {
	if err != nil {
		loggerSugar.Errorw(message, "error", err.Error())
	} else {
		loggerSugar.Error(message)
	}
}

// HandleWarning logs a warning but allows the application to continue running.
func HandleWarning(loggerSugar *zap.SugaredLogger, message string, err error) {
	if err != nil {
		loggerSugar.Warnw(message, "error", err.Error())
	} else {
		loggerSugar.Warn(message)
	}
}
