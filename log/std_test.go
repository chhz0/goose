package log

import (
	"context"
	"testing"
)

func TestStd_Info(t *testing.T) {
	Info("std info test ", String("val", "std info test "))
	Infof("std info test %s", "")
	Infow("std info test ", "val", "std info test ", "nokey")
	t.Log(Enabled())
}

func TestStd_Debug(t *testing.T) {
	SetLevel(DebugLevel)
	Debug("std debug test std debug test", String("val", "std debug test "))
	Debugf("std debug test %s", "hello")
	Debugw("std debug test ", "val", "std debug test ", "nokey")
}

func TestStd_Warn(t *testing.T) {
	Warn("std warn test std warn test", String("val", "std warn test "))
	Warnf("std warn test %s", "hello")
	Warnw("std warn test ", "val", "std warn test ")
}

func TestStd_Error(t *testing.T) {
	Error("std error test std error test", String("val", "std error test "))
	Errorf("std error test %s", "hello")
	Errorw("std error test ", "val", "std error test ")
}

func TestStd_Panic(t *testing.T) {
	Panic("std panic test std panic test", String("val", "std panic test "))
	Panicf("std panic test %s", "hello")
	Panicw("std panic test ", "val", "std panic test ")
}

func TestStd_Fatal(t *testing.T) {
	Fatal("std fatal test std fatal test", String("val", "std fatal test "))
	Fatalf("std fatal test %s", "hello")
	Fatalw("std fatal test ", "val", "std fatal test ")
}

func TestStd_V(t *testing.T) {
	V(ErrorLevel).Info("std v info test std v info test")
}

func TestStd_WithValues(t *testing.T) {
	WithValues("withVal", "std with values test ").Info("std with values test std with values test")
}

func TestStd_WithName(t *testing.T) {
	WithName("with name logger").Info("std with name test std with name test")
}

func TestStd_WithContext(t *testing.T) {
	ctx := WithContext(context.Background())
	l := FromContext(ctx)
	l.Info("std with context test std with context test")
}

func TestStd_L(t *testing.T) {
	l := L(context.Background(), "key1", "key2")
	l.Info("std l test std l test")
}

func TestStd_Named(t *testing.T) {
	Named("named-logger")
	Info("std named test std named test")
}
