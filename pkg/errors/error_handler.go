package errors

import (
	"go.uber.org/zap"
)

func HandleCriticalError(loggerSugar *zap.SugaredLogger, message string, err error) {
	if err != nil {
		loggerSugar.Fatalw(message, "error", err.Error())
	} else {
		loggerSugar.Fatal(message)
	}
}
