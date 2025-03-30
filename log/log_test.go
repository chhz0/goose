package log

import (
	"os"
	"testing"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"go.uber.org/zap/zaptest/observer"
)

func TestZapLogger_Info(t *testing.T) {
	tests := []struct {
		name     string
		msg      string
		fields   []Field
		wantLogs int
	}{
		{
			name:     "simple message",
			msg:      "test message",
			fields:   nil,
			wantLogs: 1,
		},
		{
			name:     "with fields",
			msg:      "test with fields",
			fields:   []Field{zap.String("key", "value")},
			wantLogs: 1,
		},
		{
			name:     "empty message",
			msg:      "",
			fields:   nil,
			wantLogs: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup observer
			core, recorded := observer.New(zapcore.InfoLevel)
			level := zap.NewAtomicLevelAt(zapcore.InfoLevel)
			logger := &zapLogger{
				l:  zap.New(core),
				al: &level,
			}

			// Call method
			logger.Info(tt.msg, tt.fields...)

			// Verify
			logs := recorded.All()
			if len(logs) != tt.wantLogs {
				t.Errorf("Expected %d logs, got %d", tt.wantLogs, len(logs))
			}

			if len(logs) > 0 {
				if logs[0].Message != tt.msg {
					t.Errorf("Expected message %q, got %q", tt.msg, logs[0].Message)
				}
			}
		})
	}
}

func TestLog_ErrorOutput(t *testing.T) {
	errorlog := newZaplogger(os.Stdout, InfoLevel, JsonEncoder,
		AddCaller(),
		ErrorOutput(zapcore.AddSync(OpenLogFile("error.log"))),
	)

	defer errorlog.Sync()
	errorlog.Error("error test log")
	errorlog.Error("fatal")
	errorlog.Info("info test log")
	errorlog.Warn("warn")
	errorlog.Debug("debug")
}
