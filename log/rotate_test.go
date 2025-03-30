package log

import (
	"testing"
)

func TestRotate_New(t *testing.T) {
	rotateLog := newTee([]TeeOption{
		{
			Output: NewProductionRotateBySize("rotate-by-size.log"),
			LevelEnablerFunc: func(lvl Level) bool {
				return lvl < WarnLevel
			},
		},
		{
			Output: NewProductionRotateByTime("rotate-by-time.log"),
			LevelEnablerFunc: func(lvl Level) bool {
				return lvl >= WarnLevel
			},
		},
	},
		JsonEncoder,
		defaultZapOptions()...,
	)

	defer rotateLog.Sync()
	rotateLog.Debug("Debug TeeAndRotate")
	rotateLog.Info("Info TeeAndRotate")
	rotateLog.Warn("Warn TeeAndRotate")
	rotateLog.Error("Error TeeAndRotate")
}
