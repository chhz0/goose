package log

import (
	"io"
	"strings"
	"time"

	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"gopkg.in/natefinch/lumberjack.v2"
)

// 日志轮转功能 借用file-rotatelogs lumberjack实现
type RotateConfig struct {
	// 共用配置
	Filename string // 完整配置名
	MaxAge   int    // 保留旧日志文件最大时间

	// 按时间轮转配置
	RotationTime time.Duration

	// 按大小轮转配置
	MaxSize    int  // 日志文件最大大小
	MaxBackups int  // 保留日志文件的最大数量
	Compress   bool // 是否对日志文件进行压缩
	LocalTime  bool // 是否使用本地时间
}

func NewProductionRotateByTime(filename string) io.Writer {
	return NewRotateByTime(NewProductionRotateConfig(filename))
}

func NewProductionRotateBySize(filename string) io.Writer {
	return NewRotateBySize(NewProductionRotateConfig(filename))
}

func NewProductionRotateConfig(filename string) *RotateConfig {
	return &RotateConfig{
		Filename: filename,
		MaxAge:   30,

		RotationTime: time.Hour * 24,

		MaxSize:    100,
		MaxBackups: 100,
		Compress:   true,
		LocalTime:  false,
	}
}

func NewRotateByTime(cfg *RotateConfig) io.Writer {
	opts := []rotatelogs.Option{
		rotatelogs.WithMaxAge(time.Duration(cfg.MaxAge) * time.Hour * 24),
		rotatelogs.WithRotationTime(cfg.RotationTime),
	}
	if !cfg.LocalTime {
		rotatelogs.WithClock(rotatelogs.UTC)
	}
	filename := strings.SplitN(cfg.Filename, ".", 2)
	logs, err := rotatelogs.New(
		filename[0]+".%Y-%m-%d-%H-%M-%S."+filename[1],
		opts...,
	)
	if err != nil {
		return nil
	}
	return logs
}

func NewRotateBySize(cfg *RotateConfig) io.Writer {
	return &lumberjack.Logger{
		Filename:   cfg.Filename,
		MaxSize:    cfg.MaxSize,
		MaxAge:     cfg.MaxAge,
		MaxBackups: cfg.MaxBackups,
		LocalTime:  cfg.LocalTime,
		Compress:   cfg.Compress,
	}
}
