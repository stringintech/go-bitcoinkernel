package utils

import (
	"bytes"
	"encoding/json"
	"github.com/stringintech/go-bitcoinkernel/kernel"
	"log/slog"
	"strings"
	"testing"
)

func TestLogger(t *testing.T) {
	// Enable LogKernel category to capture the block decode failure log
	kernel.AddLogLevelCategory(kernel.LogKernel, kernel.LogLevelDebug)
	kernel.EnableLogCategory(kernel.LogKernel)
	defer kernel.DisableLogCategory(kernel.LogKernel)

	var buf bytes.Buffer
	handler := slog.NewJSONHandler(&buf, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	})

	options := kernel.LoggingOptions{
		LogTimestamps:            true,
		LogTimeMicros:            true,
		LogThreadNames:           true,
		LogSourceLocations:       true,
		AlwaysPrintCategoryLevel: true,
	}

	viewer, err := NewLogger(handler, options)
	if err != nil {
		t.Fatalf("Failed to create log viewer: %v", err)
	}
	defer viewer.Destroy()

	_, err = kernel.NewBlockFromRaw([]byte{0xab})
	if err == nil {
		t.Error("Expected error for invalid block data")
	}

	// Parse the last JSON entry and assert message
	logOutput := buf.String()
	lines := strings.Split(strings.TrimSpace(logOutput), "\n")
	if len(lines) == 0 {
		t.Fatal("No log output captured")
	}

	lastLine := lines[len(lines)-1]
	var logEntry map[string]interface{}
	if err := json.Unmarshal([]byte(lastLine), &logEntry); err != nil {
		t.Fatalf("Failed to parse last JSON log entry: %v", err)
	}

	expectedMessage := "Block decode failed."
	if msg, ok := logEntry["msg"]; !ok {
		t.Error("Log entry missing 'msg' field")
	} else if msg != expectedMessage {
		t.Errorf("Expected message '%s', got '%s'", expectedMessage, msg)
	}
}
