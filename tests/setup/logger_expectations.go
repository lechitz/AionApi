package setup

import (
	mockLogger "github.com/lechitz/AionApi/tests/mocks/logger"
	"go.uber.org/mock/gomock"
)

// ExpectLoggerDefaultBehavior sets up default expectations for Infow, Errorw, Warnw, and Debugw calls on the mockLogger.MockLogger instance.
func ExpectLoggerDefaultBehavior(logger *mockLogger.MockLogger) {
	logger.EXPECT().Infow(gomock.Any(), gomock.Any()).AnyTimes()
	logger.EXPECT().Errorw(gomock.Any(), gomock.Any()).AnyTimes()
	logger.EXPECT().Warnw(gomock.Any(), gomock.Any()).AnyTimes()
	logger.EXPECT().Debugw(gomock.Any(), gomock.Any()).AnyTimes()
}
