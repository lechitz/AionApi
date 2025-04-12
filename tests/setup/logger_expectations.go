package setup

import (
	"github.com/golang/mock/gomock"
	mockLogger "github.com/lechitz/AionApi/tests/mocks/logger"
)

func ExpectLoggerDefaultBehavior(logger *mockLogger.MockLogger) {
	logger.EXPECT().Infow(gomock.Any(), gomock.Any()).AnyTimes()
	logger.EXPECT().Errorw(gomock.Any(), gomock.Any()).AnyTimes()
	logger.EXPECT().Warnw(gomock.Any(), gomock.Any()).AnyTimes()
	logger.EXPECT().Debugw(gomock.Any(), gomock.Any()).AnyTimes()
}
