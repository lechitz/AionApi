package logger

type Logger interface {
	Infof(format string, args ...any)
	Errorf(format string, args ...any)
	Debugf(format string, args ...any)
	Warnf(format string, args ...any)

	Infow(msg string, keysAndValues ...any)
	Errorw(msg string, keysAndValues ...any)
	Debugw(msg string, keysAndValues ...any)
	Warnw(msg string, keysAndValues ...any)
}
