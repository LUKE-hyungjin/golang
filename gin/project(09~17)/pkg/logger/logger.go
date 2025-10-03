package logger

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"time"

	"github.com/example/banking-system/internal/config"
)

// Logger interface
type Logger struct {
	level  string
	format string
	output io.Writer
}

// Log levels
const (
	DEBUG = "debug"
	INFO  = "info"
	WARN  = "warn"
	ERROR = "error"
	FATAL = "fatal"
)

// NewLogger creates a new logger instance
func NewLogger(cfg config.LoggingConfig) *Logger {
	var output io.Writer

	switch cfg.Output {
	case "file":
		file, err := os.OpenFile(cfg.FilePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		if err != nil {
			log.Printf("Failed to open log file, using stdout: %v", err)
			output = os.Stdout
		} else {
			output = file
		}
	default:
		output = os.Stdout
	}

	return &Logger{
		level:  cfg.Level,
		format: cfg.Format,
		output: output,
	}
}

// LogEntry represents a log entry
type LogEntry struct {
	Timestamp string                 `json:"timestamp"`
	Level     string                 `json:"level"`
	Message   string                 `json:"message"`
	Fields    map[string]interface{} `json:"fields,omitempty"`
}

// log writes a log entry
func (l *Logger) log(level, message string, fields map[string]interface{}) {
	if !l.shouldLog(level) {
		return
	}

	entry := LogEntry{
		Timestamp: time.Now().UTC().Format(time.RFC3339),
		Level:     level,
		Message:   message,
		Fields:    fields,
	}

	if l.format == "json" {
		data, _ := json.Marshal(entry)
		fmt.Fprintln(l.output, string(data))
	} else {
		// Text format
		fieldsStr := ""
		if len(fields) > 0 {
			fieldsData, _ := json.Marshal(fields)
			fieldsStr = " " + string(fieldsData)
		}
		fmt.Fprintf(l.output, "[%s] %s: %s%s\n", entry.Timestamp, entry.Level, entry.Message, fieldsStr)
	}
}

// shouldLog checks if the message should be logged based on level
func (l *Logger) shouldLog(level string) bool {
	levels := map[string]int{
		DEBUG: 0,
		INFO:  1,
		WARN:  2,
		ERROR: 3,
		FATAL: 4,
	}

	currentLevel, ok1 := levels[l.level]
	messageLevel, ok2 := levels[level]

	if !ok1 || !ok2 {
		return true
	}

	return messageLevel >= currentLevel
}

// Debug logs a debug message
func (l *Logger) Debug(message string, fields ...map[string]interface{}) {
	f := make(map[string]interface{})
	if len(fields) > 0 {
		f = fields[0]
	}
	l.log(DEBUG, message, f)
}

// Info logs an info message
func (l *Logger) Info(message string, fields ...map[string]interface{}) {
	f := make(map[string]interface{})
	if len(fields) > 0 {
		f = fields[0]
	}
	l.log(INFO, message, f)
}

// Warn logs a warning message
func (l *Logger) Warn(message string, fields ...map[string]interface{}) {
	f := make(map[string]interface{})
	if len(fields) > 0 {
		f = fields[0]
	}
	l.log(WARN, message, f)
}

// Error logs an error message
func (l *Logger) Error(message string, fields ...map[string]interface{}) {
	f := make(map[string]interface{})
	if len(fields) > 0 {
		f = fields[0]
	}
	l.log(ERROR, message, f)
}

// Fatal logs a fatal message and exits
func (l *Logger) Fatal(message string, fields ...map[string]interface{}) {
	f := make(map[string]interface{})
	if len(fields) > 0 {
		f = fields[0]
	}
	l.log(FATAL, message, f)
	os.Exit(1)
}

// WithFields returns a logger with fields
func (l *Logger) WithFields(fields map[string]interface{}) *LoggerWithFields {
	return &LoggerWithFields{
		logger: l,
		fields: fields,
	}
}

// LoggerWithFields is a logger with preset fields
type LoggerWithFields struct {
	logger *Logger
	fields map[string]interface{}
}

// Debug logs a debug message with fields
func (lf *LoggerWithFields) Debug(message string) {
	lf.logger.log(DEBUG, message, lf.fields)
}

// Info logs an info message with fields
func (lf *LoggerWithFields) Info(message string) {
	lf.logger.log(INFO, message, lf.fields)
}

// Warn logs a warning message with fields
func (lf *LoggerWithFields) Warn(message string) {
	lf.logger.log(WARN, message, lf.fields)
}

// Error logs an error message with fields
func (lf *LoggerWithFields) Error(message string) {
	lf.logger.log(ERROR, message, lf.fields)
}