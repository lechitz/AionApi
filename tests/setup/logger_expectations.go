package setup

import (
	"github.com/lechitz/AionApi/tests/mocks"
	"go.uber.org/mock/gomock"
)

// ExpectLoggerDefaultBehavior sets relaxed expectations for ContextLogger methods with ctx.
// It registers multiple arities to tolerate 0, 1, or 2 KV pairs commonly used in tests.
func ExpectLoggerDefaultBehavior(logger *mocks.MockContextLogger) {
	// InfowCtx
	logger.EXPECT().InfowCtx(gomock.Any(), gomock.Any()).AnyTimes()
	logger.EXPECT().InfowCtx(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes()
	logger.EXPECT().InfowCtx(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes()

	// ErrorwCtx
	logger.EXPECT().ErrorwCtx(gomock.Any(), gomock.Any()).AnyTimes()
	logger.EXPECT().ErrorwCtx(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes()
	logger.EXPECT().ErrorwCtx(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes()

	// WarnwCtx
	logger.EXPECT().WarnwCtx(gomock.Any(), gomock.Any()).AnyTimes()
	logger.EXPECT().WarnwCtx(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes()
	logger.EXPECT().WarnwCtx(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes()

	// DebugwCtx
	logger.EXPECT().DebugwCtx(gomock.Any(), gomock.Any()).AnyTimes()
	logger.EXPECT().DebugwCtx(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes()
	logger.EXPECT().DebugwCtx(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes()
}
