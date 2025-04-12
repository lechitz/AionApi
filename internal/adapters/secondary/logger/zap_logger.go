package logger

import "go.uber.org/zap"

type ZapLoggerAdapter struct {
	sugar *zap.SugaredLogger
}

func NewZapLoggerAdapter(sugar *zap.SugaredLogger) *ZapLoggerAdapter {
	return &ZapLoggerAdapter{sugar: sugar}
}

func (l *ZapLoggerAdapter) Infof(format string, args ...any) {
	l.sugar.Infof(format, args...)
}

func (l *ZapLoggerAdapter) Errorf(format string, args ...any) {
	l.sugar.Errorf(format, args...)
}

func (l *ZapLoggerAdapter) Warnf(format string, args ...any) {
	l.sugar.Warnf(format, args...)
}

func (l *ZapLoggerAdapter) Debugf(format string, args ...any) {
	l.sugar.Debugf(format, args...)
}

func (l *ZapLoggerAdapter) Infow(msg string, keysAndValues ...any) {
	l.sugar.Infow(msg, keysAndValues...)
}

func (l *ZapLoggerAdapter) Errorw(msg string, keysAndValues ...any) {
	l.sugar.Errorw(msg, keysAndValues...)
}

func (l *ZapLoggerAdapter) Warnw(msg string, keysAndValues ...any) {
	l.sugar.Warnw(msg, keysAndValues...)
}

func (l *ZapLoggerAdapter) Debugw(msg string, keysAndValues ...any) {
	l.sugar.Debugw(msg, keysAndValues...)
}
