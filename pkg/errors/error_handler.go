package errors

import (
	"github.com/lechitz/AionApi/pkg/contextkeys"
	"go.uber.org/zap"
)

func HandleCriticalError(loggerSugar *zap.SugaredLogger, message string, err error) {
	if err != nil {
		loggerSugar.Fatalw(message, contextkeys.Error, err.Error())
	} else {
		loggerSugar.Fatal(message)
	}
}
