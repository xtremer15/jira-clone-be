// Package logger provides structured JSON logging for the Jira Clone Backend application.
// It supports multiple log levels and structured logging with custom fields.
package logger

import (
	"encoding/json"
	"log"
	"os"
	"time"
)

// LogLevel represents the logging level with increasing severity.
// Higher values indicate more severe log levels.
type LogLevel int

const (
	DEBUG LogLevel = iota // Debug level for detailed information
	INFO                  // Info level for general information
	WARN                  // Warning level for potentially harmful situations
	ERROR                 // Error level for error events
	FATAL                 // Fatal level for severe errors that cause termination
)

// Logger represents a structured logger that outputs JSON-formatted log entries.
// It supports multiple log levels and allows adding custom fields to log entries.
type Logger struct {
	level  LogLevel    // Current logging level threshold
	logger *log.Logger // Underlying Go logger instance
}

// LogEntry represents a structured log entry with timestamp, level, message, and optional fields.
// All log entries are formatted as JSON for easy parsing by log aggregation systems.
type LogEntry struct {
	Level     string                 `json:"level"`            // Log level as string (DEBUG, INFO, WARN, ERROR, FATAL)
	Timestamp string                 `json:"timestamp"`        // ISO 8601 timestamp in UTC
	Message   string                 `json:"message"`          // Log message content
	Fields    map[string]interface{} `json:"fields,omitempty"` // Additional structured fields
}

// New creates a new logger instance with the specified logging level.
// The logger outputs structured JSON to stdout.
//
// Parameters:
//   - level: Log level string ("debug", "info", "warn", "error", "fatal")
//
// Returns a new Logger instance configured with the specified level.
// If an invalid level is provided, it defaults to INFO level.
func New(level string) *Logger {
	logLevel := parseLogLevel(level)

	return &Logger{
		level:  logLevel,
		logger: log.New(os.Stdout, "", 0),
	}
}

// parseLogLevel parses the log level string and returns the corresponding LogLevel constant.
// If an invalid level is provided, it defaults to INFO level.
//
// Parameters:
//   - level: Log level string ("debug", "info", "warn", "error", "fatal")
//
// Returns the corresponding LogLevel constant or INFO if invalid.
func parseLogLevel(level string) LogLevel {
	switch level {
	case "debug":
		return DEBUG
	case "info":
		return INFO
	case "warn":
		return WARN
	case "error":
		return ERROR
	case "fatal":
		return FATAL
	default:
		return INFO
	}
}

// log writes a structured log entry if the level meets the logger's threshold.
// It formats the entry as JSON and outputs it to stdout.
//
// Parameters:
//   - level: The log level of the entry
//   - message: The log message
//   - fields: Additional structured fields to include in the log entry
func (l *Logger) log(level LogLevel, message string, fields map[string]interface{}) {
	if level < l.level {
		return
	}

	entry := LogEntry{
		Level:     l.getLevelString(level),
		Timestamp: time.Now().UTC().Format(time.RFC3339),
		Message:   message,
		Fields:    fields,
	}

	jsonData, err := json.Marshal(entry)
	if err != nil {
		l.logger.Printf("Failed to marshal log entry: %v", err)
		return
	}

	l.logger.Println(string(jsonData))
}

// getLevelString returns the string representation of the log level.
//
// Parameters:
//   - level: The LogLevel constant
//
// Returns the string representation of the log level.
func (l *Logger) getLevelString(level LogLevel) string {
	switch level {
	case DEBUG:
		return "DEBUG"
	case INFO:
		return "INFO"
	case WARN:
		return "WARN"
	case ERROR:
		return "ERROR"
	case FATAL:
		return "FATAL"
	default:
		return "INFO"
	}
}

// Debug logs a debug level message with optional structured fields.
// Debug messages are used for detailed information that is typically only of interest
// when diagnosing problems.
//
// Parameters:
//   - message: The debug message
//   - fields: Optional structured fields to include in the log entry
func (l *Logger) Debug(message string, fields ...map[string]interface{}) {
	fieldsMap := mergeFields(fields...)
	l.log(DEBUG, message, fieldsMap)
}

// Info logs an info level message with optional structured fields.
// Info messages are used for general information about the application's operation.
//
// Parameters:
//   - message: The info message
//   - fields: Optional structured fields to include in the log entry
func (l *Logger) Info(message string, fields ...map[string]interface{}) {
	fieldsMap := mergeFields(fields...)
	l.log(INFO, message, fieldsMap)
}

// Warn logs a warning level message with optional structured fields.
// Warning messages are used for potentially harmful situations that should be noted.
//
// Parameters:
//   - message: The warning message
//   - fields: Optional structured fields to include in the log entry
func (l *Logger) Warn(message string, fields ...map[string]interface{}) {
	fieldsMap := mergeFields(fields...)
	l.log(WARN, message, fieldsMap)
}

// Error logs an error level message with optional structured fields.
// Error messages are used for error events that might still allow the application to continue running.
//
// Parameters:
//   - message: The error message
//   - fields: Optional structured fields to include in the log entry
func (l *Logger) Error(message string, fields ...map[string]interface{}) {
	fieldsMap := mergeFields(fields...)
	l.log(ERROR, message, fieldsMap)
}

// Fatal logs a fatal level message with optional structured fields and exits the application.
// Fatal messages are used for very severe error events that will presumably lead the application to abort.
// After logging the message, the application will exit with code 1.
//
// Parameters:
//   - message: The fatal message
//   - fields: Optional structured fields to include in the log entry
func (l *Logger) Fatal(message string, fields ...map[string]interface{}) {
	fieldsMap := mergeFields(fields...)
	l.log(FATAL, message, fieldsMap)
	os.Exit(1)
}

// mergeFields merges multiple field maps into a single map.
// If duplicate keys exist, later maps will override earlier ones.
//
// Parameters:
//   - fields: Variable number of field maps to merge
//
// Returns a single map containing all fields from the input maps.
func mergeFields(fields ...map[string]interface{}) map[string]interface{} {
	result := make(map[string]interface{})
	for _, fieldMap := range fields {
		for k, v := range fieldMap {
			result[k] = v
		}
	}
	return result
}
