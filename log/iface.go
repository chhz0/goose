package log

import "context"

// InfoLogger 作为基础的日志接口，仅提供info级别的日志接口，
type InfoLogger interface {
	Info(msg string, fields ...Field)
	Infof(format string, args ...any)
	Infow(msg string, keysAndValues ...any)

	Enabled() bool
}

// Logger 完整的日志接口
type Logger interface {
	InfoLogger
	Debug(msg string, fields ...Field)
	Debugf(format string, args ...any)
	Debugw(msg string, keysAndValues ...any)
	Warn(msg string, fields ...Field)
	Warnf(format string, args ...any)
	Warnw(msg string, keysAndValues ...any)
	Error(msg string, fields ...Field)
	Errorf(format string, args ...any)
	Errorw(msg string, keysAndValues ...any)
	Panic(msg string, fields ...Field)
	Panicf(format string, args ...any)
	Panicw(msg string, keysAndValues ...any)
	Fatal(msg string, fields ...Field)
	Fatalf(format string, args ...any)
	Fatalw(msg string, keysAndValues ...any)

	// V 返回一个 InfoLogger 允许通过传入的日志级别灵活指定 InfoLogger.
	V(level Level) InfoLogger

	// WithValues 添加一个 key-value 键值对到 logger 中.
	WithValues(keysAndvalues ...any) Logger

	// WithName 为 logger 添加一个名称，方便在日志中区分日志来源.
	WithName(name string) Logger

	// WithContext
	WithContext(ctx context.Context) context.Context

	// L 从 context 中获取对应的值添加到 logger 中
	L(ctx context.Context, keys ...string) Logger

	Sync()
}
