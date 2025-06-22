package utils

import (
	"fmt"
	"github.com/stringintech/go-bitcoinkernel/kernel"
	"testing"
	"time"
)

func TestLogParserParse(t *testing.T) {
	timeStr := "2025-03-19T12:14:55Z"
	thread := "unknown"
	filename := "context.cpp"
	path := fmt.Sprintf("depend/bitcoin/src/kernel/%s", filename)
	lineno := 20
	function := "operator()"
	category := "all"
	level := "info"
	msg := "Using the 'arm_shani(1way,2way)' SHA256 implementation"
	logString := fmt.Sprintf("%s [%s] [%s:%d] [%s] [%s:%s] %s",
		timeStr, thread, path, lineno, function, category, level, msg)

	options := kernel.LoggingOptions{
		LogTimestamps:            true,
		LogTimeMicros:            false,
		LogThreadNames:           true,
		LogSourceLocations:       true,
		AlwaysPrintCategoryLevel: true,
	}

	parser, err := NewLogParser(options)
	if err != nil {
		t.Fatalf("Failed to create log parser: %v", err)
	}

	// Parse the log string
	entry, err := parser.Parse(logString)
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}

	expectedTime, _ := time.Parse(time.RFC3339, timeStr)
	if !entry.Time.Equal(expectedTime) {
		t.Errorf("Expected timestamp %v, got %v", expectedTime, entry.Time)
	}

	if entry.Thread != thread {
		t.Errorf("Expected thread %s, got %s", thread, entry.Thread)
	}

	if entry.Path != path {
		t.Errorf("Expected path %s, got %s", filename, entry.File)
	}

	if entry.File != filename {
		t.Errorf("Expected filename %s, got %s", filename, entry.File)
	}

	if entry.LineNumber != lineno {
		t.Errorf("Expected line number %d, got %d", lineno, entry.LineNumber)
	}

	if entry.Function != function {
		t.Errorf("Expected function %s, got %s", function, entry.Function)
	}

	if entry.Category != category {
		t.Errorf("Expected category %s, got %s", category, entry.Category)
	}

	expectedLevel := parseLogLevel(level)
	if entry.Level != expectedLevel {
		t.Errorf("Expected level %v, got %v", expectedLevel, entry.Level)
	}

	if entry.Message != msg {
		t.Errorf("Expected message %s, got %s", msg, entry.Message)
	}
}

// TestRegexCompilation tests regex compilation with different options
func TestRegexCompilation(t *testing.T) {
	testCases := []struct {
		name    string
		options kernel.LoggingOptions
	}{
		{
			name: "no options",
			options: kernel.LoggingOptions{
				LogTimestamps:            false,
				LogTimeMicros:            false,
				LogThreadNames:           false,
				LogSourceLocations:       false,
				AlwaysPrintCategoryLevel: false,
			},
		},
		{
			name: "timestamps only",
			options: kernel.LoggingOptions{
				LogTimestamps:            true,
				LogTimeMicros:            false,
				LogThreadNames:           false,
				LogSourceLocations:       false,
				AlwaysPrintCategoryLevel: false,
			},
		},
		{
			name: "microsecond timestamps",
			options: kernel.LoggingOptions{
				LogTimestamps:            true,
				LogTimeMicros:            true,
				LogThreadNames:           false,
				LogSourceLocations:       false,
				AlwaysPrintCategoryLevel: false,
			},
		},
		{
			name: "all options",
			options: kernel.LoggingOptions{
				LogTimestamps:            true,
				LogTimeMicros:            true,
				LogThreadNames:           true,
				LogSourceLocations:       true,
				AlwaysPrintCategoryLevel: true,
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			parser := &LogParser{options: tc.options}
			err := parser.compileRegexes()
			if err != nil {
				t.Errorf("Failed to compile regexes for %s: %v", tc.name, err)
			}

			if parser.masterRegex == nil {
				t.Errorf("Expected regex to be compiled for %s", tc.name)
			}
		})
	}
}
