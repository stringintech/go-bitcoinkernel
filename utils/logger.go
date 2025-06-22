package utils

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/stringintech/go-bitcoinkernel/kernel"
)

// Logger pipes Kernel logs to structured Go slog entries
type Logger struct {
	logger *slog.Logger
	parser *LogParser
	conn   *kernel.LoggingConnection
}

func NewLogger(handler slog.Handler, options kernel.LoggingOptions) (*Logger, error) {
	parser, err := NewLogParser(options)
	if err != nil {
		return nil, fmt.Errorf("failed to create log parser: %w", err)
	}

	viewer := &Logger{
		logger: slog.New(handler),
		parser: parser,
	}

	// Create logging connection with our message handler
	conn, err := kernel.NewLoggingConnection(viewer.handleMessage, options)
	if err != nil {
		return nil, fmt.Errorf("failed to create logging connection: %w", err)
	}

	viewer.conn = conn
	return viewer, nil
}

// handleMessage processes incoming log messages from the Kernel
func (l *Logger) handleMessage(message string) {
	entry, err := l.parser.Parse(message)
	if err != nil {
		l.logger.Warn("Failed to parse log message",
			slog.String("error", err.Error()),
			slog.String("raw_message", message))

		// Create a fallback entry with the raw message
		entry = &LogEntry{
			Message: message,
			Level:   slog.LevelInfo,
			Time:    time.Now(),
		}
	}
	l.log(entry)
}

func (l *Logger) log(entry *LogEntry) {
	l.logger.LogAttrs(context.TODO(), entry.Level, entry.Message, entry.GetSlogAttributes()...)
}

func (l *Logger) Destroy() {
	if l.conn != nil {
		l.conn.Destroy()
		l.conn = nil
	}
}
