package log

import (
	"io"
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type LevelEnablerFunc func(lvl Level) bool

type TeeOption struct {
	Output io.Writer
	LevelEnablerFunc
}

// newTee creates a new tee logger.
// 默认使用的 日志级别为 InfoLevel，
func newTee(tees []TeeOption, encoder LogEncoder, opts ...ZapOption) *zapLogger {
	cores := make([]zapcore.Core, 0, len(tees))

	for _, tee := range tees {
		if tee.Output == nil {
			tee.Output = os.Stdout
		}

		encoderConfig := defaultEncoderConfig()
		core := zapcore.NewCore(
			coreEncoder(encoder, encoderConfig),
			zapcore.AddSync(tee.Output),
			zap.LevelEnablerFunc(tee.LevelEnablerFunc),
		)
		cores = append(cores, core)
	}

	atomicl := zap.NewAtomicLevelAt(InfoLevel)
	return &zapLogger{
		l:  zap.New(zapcore.NewTee(cores...), opts...),
		al: &atomicl,
	}
}

func OpenLogFile(file string) io.Writer {
	logf, err := os.OpenFile(file, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		panic("open log file failed")
	}
	return logf
}

func NewTeeLogger(tees []TeeOption, encoder LogEncoder, opts ...ZapOption) Logger {
	return newTee(tees, encoder, opts...)
}
