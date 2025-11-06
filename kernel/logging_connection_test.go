package kernel

import (
	"testing"
)

func TestLoggingOptions(t *testing.T) {
	// Test different logging options
	testCases := []struct {
		name    string
		options LoggingOptions
	}{
		{
			name: "minimal options",
			options: LoggingOptions{
				LogTimestamps:            false,
				LogTimeMicros:            false,
				LogThreadNames:           false,
				LogSourceLocations:       false,
				AlwaysPrintCategoryLevel: false,
			},
		},
		{
			name: "full options",
			options: LoggingOptions{
				LogTimestamps:            true,
				LogTimeMicros:            true,
				LogThreadNames:           true,
				LogSourceLocations:       true,
				AlwaysPrintCategoryLevel: true,
			},
		},
		{
			name: "timestamps only",
			options: LoggingOptions{
				LogTimestamps:            true,
				LogTimeMicros:            false,
				LogThreadNames:           false,
				LogSourceLocations:       false,
				AlwaysPrintCategoryLevel: false,
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			SetLoggingOptions(tc.options)
			conn, err := NewLoggingConnection(func(_ string) {})
			if err != nil {
				t.Fatalf("NewLoggingConnection() error = %v", err)
			}
			defer conn.Destroy()

			if conn.ptr == nil {
				t.Error("Expected connection pointer to be set")
			}
		})
	}
}

func TestAddLogLevelCategory(t *testing.T) {
	// Test adding log level categories
	AddLogLevelCategory(LogAll, LogLevelDebug)
	AddLogLevelCategory(LogAll, LogLevelInfo)
	AddLogLevelCategory(LogAll, LogLevelInfo) // Same operation twice should succeed

	AddLogLevelCategory(LogBlockStorage, LogLevelDebug)
	AddLogLevelCategory(LogValidation, LogLevelTrace)
	AddLogLevelCategory(LogKernel, LogLevelInfo)
}

func TestEnableDisableLogCategory(t *testing.T) {
	EnableLogCategory(LogBlockStorage)
	EnableLogCategory(LogBlockStorage) // Same operation twice should succeed
	DisableLogCategory(LogBlockStorage)
	DisableLogCategory(LogBlockStorage) // Same operation twice should succeed

	EnableLogCategory(LogValidation)
	DisableLogCategory(LogValidation)

	EnableLogCategory(LogKernel)
	DisableLogCategory(LogKernel)

	EnableLogCategory(LogAll)
	DisableLogCategory(LogAll)
}
