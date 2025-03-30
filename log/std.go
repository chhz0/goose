package log

import (
	"context"
	"os"
	"sync/atomic"
)

var stdLogger atomic.Pointer[zapLogger]

func init() {
	stdLogger.Store(newZaplogger(os.Stdout, InfoLevel, JsonEncoder, defaultZapOptions()...))
}
func ZapLogger() *zapLogger {
	return stdLogger.Load()
}

func Info(msg string, fields ...Field) {
	stdLogger.Load().Info(msg, fields...)
}

func Infof(format string, args ...any) {
	stdLogger.Load().Infof(format, args...)
}

func Infow(msg string, keysAndValues ...any) {
	stdLogger.Load().Infow(msg, keysAndValues...)
}

func Enabled() bool {
	return stdLogger.Load().Enabled()
}

func Debug(msg string, fields ...Field) {
	stdLogger.Load().Debug(msg, fields...)
}

func Debugf(format string, args ...any) {
	stdLogger.Load().Debugf(format, args...)
}

func Debugw(msg string, keysAndValues ...any) {
	stdLogger.Load().Debugw(msg, keysAndValues...)
}

func Warn(msg string, fields ...Field) {
	stdLogger.Load().Warn(msg, fields...)
}

func Warnf(format string, args ...any) {
	stdLogger.Load().Warnf(format, args...)
}

func Warnw(msg string, keysAndValues ...any) {
	stdLogger.Load().Warnw(msg, keysAndValues...)
}

func Error(msg string, fields ...Field) {
	stdLogger.Load().Error(msg, fields...)
}

func Errorf(format string, args ...any) {
	stdLogger.Load().Errorf(format, args...)
}

func Errorw(msg string, keysAndValues ...any) {
	stdLogger.Load().Errorw(msg, keysAndValues...)
}

func Panic(msg string, fields ...Field) {
	stdLogger.Load().Panic(msg, fields...)
}

func Panicf(format string, args ...any) {
	stdLogger.Load().Panicf(format, args...)
}

func Panicw(msg string, keysAndValues ...any) {
	stdLogger.Load().Panicw(msg, keysAndValues...)
}

func Fatal(msg string, fields ...Field) {
	stdLogger.Load().Fatal(msg, fields...)
}

func Fatalf(format string, args ...any) {
	stdLogger.Load().Fatalf(format, args...)
}

func Fatalw(msg string, keysAndValues ...any) {
	stdLogger.Load().Fatalw(msg, keysAndValues...)
}

func V(level Level) InfoLogger {
	return stdLogger.Load().V(level)
}

func WithValues(keysAndvalues ...any) Logger {
	return stdLogger.Load().WithValues(keysAndvalues...)
}

func WithName(name string) Logger {
	return stdLogger.Load().WithName(name)
}

func WithContext(ctx context.Context) context.Context {
	return stdLogger.Load().WithContext(ctx)
}

func L(ctx context.Context, keys ...string) Logger {
	return stdLogger.Load().L(ctx, keys...)
}

func Sync() {
	stdLogger.Load().Sync()
}

func SetLevel(level Level) {
	stdLogger.Load().SetLevel(level)
}

func ReplaceDefault(l Logger) {
	stdLogger.Store(l.(*zapLogger))
}

func Named(name string) {
	nl := stdLogger.Load().WithName(name)
	stdLogger.Store(nl.(*zapLogger))
}

func FromContext(ctx context.Context) Logger {
	if ctx != nil {
		logger := ctx.Value(loggerKey)
		if logger != nil {
			return logger.(Logger)
		}
	}

	return nil
}
