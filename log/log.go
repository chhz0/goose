// 基于zap zapcore进行二次封装，扩展的log日志库，使用支持的弱类型的sugar操作
package log

import (
	"context"
	"errors"
	"io"
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Level = zapcore.Level

const (
	DebugLevel = zapcore.DebugLevel
	InfoLevel  = zapcore.InfoLevel
	WarnLevel  = zapcore.WarnLevel
	ErrorLevel = zapcore.ErrorLevel
	PanicLevel = zapcore.PanicLevel
	FatalLevel = zapcore.FatalLevel
)

type zapLogger struct {
	l  *zap.Logger
	al *zap.AtomicLevel
}

func (l *zapLogger) Info(msg string, fields ...Field) {
	l.check(msg, fields...)
}

func (l *zapLogger) Infof(format string, args ...any) {
	l.checkf(format, args...)
}

func (l *zapLogger) Infow(msg string, keysAndValues ...any) {
	l.checkw(msg, keysAndValues...)
}

func (l *zapLogger) Enabled() bool {
	return l.l.Core().Enabled(l.al.Level())
}

func (l *zapLogger) check(msg string, fields ...Field) {
	if len(msg) > 1024 {
		msg = msg[:1024]
	}
	if l.al.Enabled(zapcore.Level(l.al.Level())) {
		l.l.Info(msg, fields...)
	}
}

func (l *zapLogger) checkf(format string, args ...any) {
	if len(format) > 1024 {
		format = format[:1024]
	}
	if l.al.Enabled(zapcore.Level(l.al.Level())) {
		l.l.Sugar().Infof(format, args...)
	}
}

func (l *zapLogger) checkw(msg string, kaysAndvalues ...any) {
	if l.al.Enabled(zapcore.Level(l.al.Level())) {
		l.l.Sugar().Infow(msg, kaysAndvalues...)
	}
}

func handleFields(zl *zap.Logger, keysAndValues ...any) []zap.Field {
	fields := make([]zap.Field, 0, len(keysAndValues)/2)

	for i := 0; i < len(keysAndValues); i += 2 {
		if i == len(keysAndValues)-1 {
			zl.Error("Ignored key without a value.", zap.Any("ignored", keysAndValues[i]))
			break
		}

		key, val := keysAndValues[i], keysAndValues[i+1]
		keyStr, isString := key.(string)
		if !isString {
			zl.DPanic("non-string key supplied", zap.Any("invalid key", key))
			continue
		}
		// 根据字段类型选择更合适的 zap.Field 构造函数
		switch v := val.(type) {
		case string:
			fields = append(fields, zap.String(keyStr, v))
		case int:
			fields = append(fields, zap.Int(keyStr, v))
		case float64:
			fields = append(fields, zap.Float64(keyStr, v))
		case bool:
			fields = append(fields, zap.Bool(keyStr, v))
		default:
			fields = append(fields, zap.Any(keyStr, v))
		}
	}

	return fields
}

func (l *zapLogger) Debug(msg string, fields ...Field) {
	l.l.Debug(msg, fields...)
}

func (l *zapLogger) Debugf(format string, args ...any) {
	l.l.Sugar().Debugf(format, args...)
}

func (l *zapLogger) Debugw(msg string, keysAndValues ...any) {
	l.l.Sugar().Debugw(msg, keysAndValues...)
}

func (l *zapLogger) Warn(msg string, fields ...Field) {
	l.l.Warn(msg, fields...)
}

func (l *zapLogger) Warnf(format string, args ...any) {
	l.l.Sugar().Warnf(format, args...)
}

func (l *zapLogger) Warnw(msg string, keysAndValues ...any) {
	l.l.Sugar().Warnw(msg, keysAndValues...)
}

func (l *zapLogger) Error(msg string, fields ...Field) {
	l.l.Error(msg, fields...)
}

func (l *zapLogger) Errorf(format string, args ...any) {
	l.l.Sugar().Errorf(format, args...)
}

func (l *zapLogger) Errorw(msg string, keysAndValues ...any) {
	l.l.Sugar().Errorw(msg, keysAndValues...)
}

func (l *zapLogger) Panic(msg string, fields ...Field) {
	l.l.Panic(msg, fields...)
}

func (l *zapLogger) Panicf(format string, args ...any) {
	l.l.Sugar().Panicf(format, args...)
}

func (l *zapLogger) Panicw(msg string, keysAndValues ...any) {
	l.l.Sugar().Panicw(msg, keysAndValues...)
}

func (l *zapLogger) Fatal(msg string, fields ...Field) {
	l.l.Fatal(msg, fields...)
}

func (l *zapLogger) Fatalf(format string, args ...any) {
	l.l.Sugar().Fatalf(format, args...)
}

func (l *zapLogger) Fatalw(msg string, keysAndValues ...any) {
	l.l.Sugar().Fatalw(msg, keysAndValues...)
}

func (l *zapLogger) V(level Level) InfoLogger {
	vlog := l.clone()
	if vlog.al != nil {
		vlog.al.SetLevel(level)
	}
	return vlog
}

func (l *zapLogger) WithValues(keysAndvalues ...any) Logger {
	valuesLog := l.l.With(handleFields(l.l, keysAndvalues...)...)

	return &zapLogger{
		l:  valuesLog,
		al: l.al,
	}
}

func (l *zapLogger) WithName(name string) Logger {
	newLogger := l.l.Named(name)
	return &zapLogger{
		l:  newLogger,
		al: l.al,
	}
}

type loggerContext int

const loggerKey loggerContext = iota

func (l *zapLogger) WithContext(ctx context.Context) context.Context {
	return context.WithValue(ctx, loggerKey, l)
}

func (l *zapLogger) L(ctx context.Context, keys ...string) Logger {
	if len(keys) == 0 || ctx == nil {
		return l
	}

	// 预检查避免不必要的clone
	hasValues := false
	for _, key := range keys {
		if ctx.Value(key) != nil {
			hasValues = true
			break
		}
	}
	if !hasValues {
		return l
	}

	llog := l.clone()
	fields := make([]zap.Field, 0, len(keys))

	for _, key := range keys {
		if value := ctx.Value(key); value != nil {
			fields = append(fields, zap.Any(key, value))
		}
	}

	if len(fields) > 0 {
		llog.l = llog.l.With(fields...)
	}

	return llog
}

func (l *zapLogger) Sync() {
	if err := l.l.Sync(); err != nil && !errors.Is(err, os.ErrInvalid) {
		l.l.Error("Failed to sync logger", zap.Error(err))
	}
}

func (l *zapLogger) clone() *zapLogger {
	clone := *l
	return &clone
}

func (l *zapLogger) SetLevel(level Level) {
	if l.al != nil {
		l.al.SetLevel(level)
	}
}

type Output func() io.Writer

func NewLogger(out Output, lvl Level, encoder LogEncoder, zOpts ...ZapOption) Logger {
	return newZaplogger(out(), lvl, encoder, zOpts...)
}

type LogEncoder string

const (
	ConsoleEncoder LogEncoder = "Console"
	JsonEncoder    LogEncoder = "Json"
)

func newZaplogger(out io.Writer, level Level, encoder LogEncoder, opts ...ZapOption) *zapLogger {
	if out == nil {
		out = os.Stdout
	}

	atomicl := zap.NewAtomicLevelAt(level)
	encoderConfig := defaultEncoderConfig()
	core := zapcore.NewCore(
		coreEncoder(encoder, encoderConfig),
		zapcore.AddSync(out),
		atomicl,
	)

	return &zapLogger{
		l:  zap.New(core, opts...),
		al: &atomicl,
	}
}

func coreEncoder(encoder LogEncoder, encoderConfig *zapcore.EncoderConfig) zapcore.Encoder {
	switch encoder {
	case ConsoleEncoder:
		return zapcore.NewConsoleEncoder(*encoderConfig)
	case JsonEncoder:
		return zapcore.NewJSONEncoder(*encoderConfig)
	default:
		return zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig())
	}
}
