package utils

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"

	"log/slog"

	"github.com/stringintech/go-bitcoinkernel/kernel"
)

// LogParser handles parsing of raw log messages from Kernel
type LogParser struct {
	masterRegex *regexp.Regexp
	options     kernel.LoggingOptions
}

func NewLogParser(options kernel.LoggingOptions) (*LogParser, error) {
	parser := &LogParser{
		options: options,
	}

	// Compile regex patterns based on enabled options
	if err := parser.compileRegexes(); err != nil {
		return nil, fmt.Errorf("failed to compile regex patterns: %w", err)
	}

	return parser, nil
}

// Parse extracts structured information from a raw log message using single regex
func (p *LogParser) Parse(message string) (*LogEntry, error) {
	entry := LogEntry{
		Message: message,
		Level:   slog.LevelInfo, // Default level
	}

	if p.masterRegex == nil {
		return nil, fmt.Errorf("regex not compiled")
	}

	matches := p.masterRegex.FindStringSubmatch(strings.TrimSpace(message))
	if matches == nil {
		return nil, fmt.Errorf("failed to parse log format")
	}

	// Extract groups based on regex pattern:
	// 1: timestamp, 2: thread, 3: source location, 4: function, 5: category, 6: level, 7: message
	if len(matches) < 8 {
		return nil, fmt.Errorf("insufficient regex groups: expected 8, got %d", len(matches))
	}

	// Parse timestamp
	if matches[1] != "" {
		// supports parsing time with/without microseconds
		ts, err := time.Parse(time.RFC3339, matches[1])
		if err != nil {
			return nil, fmt.Errorf("failed to parse timestamp: %w", err)
		}
		entry.Time = ts
	}

	// Parse thread
	entry.Thread = matches[2]

	// Parse source location and extract filename/line number
	if matches[3] != "" {
		entry.Source = matches[3]
		if parts := strings.Split(matches[3], ":"); len(parts) >= 2 {
			entry.Path = strings.TrimSpace(parts[0])
			lineNum, err := strconv.Atoi(parts[len(parts)-1])
			if err != nil {
				return nil, fmt.Errorf("failed to parse line number: %w", err)
			}
			entry.LineNumber = lineNum
			// Extract the filename from the full path
			if idx := strings.LastIndex(entry.Path, "/"); idx != -1 {
				entry.File = entry.Path[idx+1:]
			} else {
				entry.File = entry.Path
			}
		}
	}

	// Parse function name
	entry.Function = matches[4]

	// Parse category and level
	entry.Category = matches[5]
	if matches[6] != "" {
		entry.Level = parseLogLevel(matches[6])
	}

	// Extract message
	entry.Message = strings.TrimSpace(matches[7])

	return &entry, nil
}

// compileRegexes builds and compiles a single masterRegex for parsing
func (p *LogParser) compileRegexes() error {
	pattern := "^"

	// Time group (optional)
	if p.options.LogTimestamps {
		//TODO consider p.options.LogTimeMicros
		pattern += `(\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}(?:\.\d+)?Z?)\s+`
	} else {
		pattern += "()"
	}

	// Thread name group (optional)
	if p.options.LogThreadNames {
		pattern += `\[([^\]]+)\]\s+`
	} else {
		pattern += "()"
	}

	// Source location group (optional)
	if p.options.LogSourceLocations {
		pattern += `\[([^\]]+)\]\s+`
	} else {
		pattern += "()"
	}

	// Function name group (optional)
	if p.options.LogSourceLocations {
		pattern += `\[([^\]]+)\]\s+`
	} else {
		pattern += "()"
	}

	// Category and level group (optional)
	if p.options.AlwaysPrintCategoryLevel {
		pattern += `\[([^:]+):([^\]]+)\]\s+`
	} else {
		pattern += "()()"
	}

	// Message (everything remaining)
	pattern += "(.+)$"

	var err error
	p.masterRegex, err = regexp.Compile(pattern)
	if err != nil {
		return fmt.Errorf("failed to compile master regex: %w", err)
	}

	return nil
}

// parseLogLevel converts string log level to slog.Level
func parseLogLevel(levelStr string) slog.Level {
	switch strings.ToLower(levelStr) {
	case "trace":
		return slog.LevelDebug - 4 // Custom trace level
	case "debug":
		return slog.LevelDebug
	case "info":
		return slog.LevelInfo
	default:
		return slog.LevelInfo
	}
}
