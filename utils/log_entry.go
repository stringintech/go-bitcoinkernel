package utils

import (
	"log/slog"
	"time"
)

// LogEntry represents a parsed log message with extracted components
type LogEntry struct {
	Time       time.Time
	Thread     string
	Category   string
	Level      slog.Level
	Source     string
	Function   string
	File       string
	Path       string
	LineNumber int
	Message    string
}

func (entry *LogEntry) GetSlogAttributes() []slog.Attr {
	var attrs []slog.Attr

	if !entry.Time.IsZero() {
		attrs = append(attrs, slog.Time("time", entry.Time))
	}

	if entry.Thread != "" {
		attrs = append(attrs, slog.String("thread", entry.Thread))
	}

	if entry.Category != "" {
		attrs = append(attrs, slog.String("category", entry.Category))
	}

	if entry.Source != "" {
		attrs = append(attrs, slog.String("source", entry.Source))
	}

	if entry.Path != "" {
		attrs = append(attrs, slog.String("path", entry.Path))
	}

	if entry.File != "" {
		attrs = append(attrs, slog.String("file", entry.File))
	}

	if entry.LineNumber > 0 {
		attrs = append(attrs, slog.Int("line_number", entry.LineNumber))
	}

	if entry.Function != "" {
		attrs = append(attrs, slog.String("function", entry.Function))
	}

	return attrs
}
