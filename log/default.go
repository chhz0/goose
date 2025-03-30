package log

import (
	"os"

	"go.uber.org/zap/zapcore"
)

func defaultZapOptions() []ZapOption {
	return []ZapOption{
		WithCaller(true),
		Development(),
		ErrorOutput(zapcore.AddSync(os.Stderr)),
		AddStacktrace(zapcore.PanicLevel),
		AddCallerSkip(1),
		WithFatalHook(zapcore.WriteThenNoop),
	}
}

func defaultEncoderConfig() *zapcore.EncoderConfig {
	return &zapcore.EncoderConfig{
		TimeKey:        "timestamp",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		FunctionKey:    zapcore.OmitKey,
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder,
		EncodeTime:     zapcore.RFC3339TimeEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}
}
