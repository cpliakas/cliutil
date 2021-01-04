package cliutil_test

import (
	"bytes"
	"context"
	"fmt"
	"strings"
	"testing"

	"github.com/cpliakas/cliutil"
)

func TestLogLevel(t *testing.T) {
	tests := []struct {
		level string
		ex    int
	}{
		{"none", cliutil.LogLevelNone},
		{"fatal", cliutil.LogLevelFatal},
		{"error", cliutil.LogLevelError},
		{"notice", cliutil.LogLevelNotice},
		{"info", cliutil.LogLevelInfo},
		{"debug", cliutil.LogLevelDebug},
		{"xyz", cliutil.LogLevelInfo},
	}

	for _, tt := range tests {
		if actual := cliutil.LogLevel(tt.level); actual != tt.ex {
			t.Errorf("got %v, expected %v", actual, tt.ex)
		}
	}
}
func TestLogLevelValid(t *testing.T) {
	tests := []struct {
		level string
		ex    bool
	}{
		{"none", true},
		{"fatal", true},
		{"error", true},
		{"notice", true},
		{"info", true},
		{"debug", true},
		{"xyz", false},
	}

	for _, tt := range tests {
		if actual := cliutil.LogLevelValid(tt.level); actual != tt.ex {
			t.Errorf("got %t, expected %t", actual, tt.ex)
		}
	}
}

func TestNewLoggerWithContext(t *testing.T) {
	ctx := context.WithValue(context.Background(), "test", "test")
	ctx, logger, xid := cliutil.NewLoggerWithContext(ctx, cliutil.LogInfo)

	var buf bytes.Buffer
	logger.SetOutput(&buf)
	logger.Info(ctx, "test message")

	ex := fmt.Sprintf("transid=%s", xid.String())
	if !strings.Contains(buf.String(), ex) {
		t.Error("expected true, got false")
	}
}

func TestContextWithLogTag(t *testing.T) {
	ctx := context.Background()
	ctx, logger, _ := cliutil.NewLoggerWithContext(ctx, cliutil.LogInfo)

	ctx = cliutil.ContextWithLogTag(ctx, "testkey", "test value")

	var buf bytes.Buffer
	logger.SetOutput(&buf)
	logger.Info(ctx, "test message")

	ex := `testkey="test value"`
	if !strings.Contains(buf.String(), ex) {
		t.Error("expected true, got false")
	}
}

func TestContextWithLogTagPanic(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("expeccted panic")
		}
	}()

	cliutil.ContextWithLogTag(context.Background(), "test key", "test value")
}
