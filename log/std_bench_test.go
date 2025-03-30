package log

import (
	"context"
	"testing"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"go.uber.org/zap/zaptest/observer"
)

type ctxType struct{}

func BenchmarkLoggers(b *testing.B) {
	// 设置测试环境
	core, _ := observer.New(zapcore.InfoLevel)
	infoL := zap.NewAtomicLevelAt(zapcore.InfoLevel)
	logger := &zapLogger{
		l:  zap.New(core),
		al: &infoL,
	}
	ReplaceDefault(logger)

	// 基准测试组
	b.Run("Info", func(b *testing.B) {
		b.ReportAllocs()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				Info("benchmark message")
			}
		})
	})

	b.Run("Infow", func(b *testing.B) {
		b.ReportAllocs()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				Infow("benchmark message",
					"string", "test",
					"int", 123,
					"float", 3.14,
					"bool", true)
			}
		})
	})

	b.Run("WithValues", func(b *testing.B) {
		b.ReportAllocs()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				l := WithValues(
					"request_id", "12345",
					"user_id", "67890")
				l.Info("message with values")
			}
		})
	})

	b.Run("Context", func(b *testing.B) {
		ctx := context.WithValue(context.Background(), "trace_id", "abc123")
		b.ReportAllocs()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				l := L(ctx, "trace_id")
				l.Info("message from context")
			}
		})
	})
}

func BenchmarkLevels(b *testing.B) {
	core, _ := observer.New(zapcore.DebugLevel)
	debugL := zap.NewAtomicLevelAt(zapcore.DebugLevel)
	logger := &zapLogger{
		l:  zap.New(core),
		al: &debugL,
	}
	ReplaceDefault(logger)

	levels := []struct {
		name string
		fn   func(string, ...Field)
	}{
		{"Debug", Debug},
		{"Info", Info},
		{"Warn", Warn},
		{"Error", Error},
	}

	for _, lvl := range levels {
		b.Run(lvl.name, func(b *testing.B) {
			b.ReportAllocs()
			for i := 0; i < b.N; i++ {
				lvl.fn("level test message")
			}
		})
	}
}

func BenchmarkStdLogger_Info(b *testing.B) {
	core, _ := observer.New(zapcore.InfoLevel)
	infoL := zap.NewAtomicLevelAt(zapcore.InfoLevel)
	logger := &zapLogger{
		l:  zap.New(core),
		al: &infoL,
	}
	ReplaceDefault(logger)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		Info("benchmark message")
	}
}

func BenchmarkStdLogger_Infow(b *testing.B) {
	core, _ := observer.New(zapcore.InfoLevel)
	infoL := zap.NewAtomicLevelAt(zapcore.InfoLevel)
	logger := &zapLogger{
		l:  zap.New(core),
		al: &infoL,
	}
	ReplaceDefault(logger)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		Infow("benchmark message", "key1", "value1", "key2", 123, "hello", "world")
	}
}
