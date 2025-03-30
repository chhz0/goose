package log

import "testing"

func TestTee_New(t *testing.T) {
	teelog := newTee([]TeeOption{
		{
			Output: OpenLogFile("tee_test.log"),
			LevelEnablerFunc: func(lvl Level) bool {
				return lvl >= InfoLevel
			},
		},
	},
		JsonEncoder,
		defaultZapOptions()...,
	)

	teelog.Info("tee test")
}
