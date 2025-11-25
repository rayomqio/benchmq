package logger

import (
	"context"
	"io"
	"log/slog"
	"os"
	"strings"
	"sync"
)

// LogLevel represents logging levels
type LogLevel int

const (
	LevelDebug LogLevel = iota
	LevelInfo
	LevelWarn
	LevelError
	LevelFatal
)

// Logger wraps slog.Logger with MQTT-specific functionality
type Logger struct {
	*slog.Logger
	level     LogLevel
	component string
}

// Config holds logger configuration
type Config struct {
	Level       LogLevel
	Format      string // "json" or "text"
	Output      io.Writer
	Component   string
	ShowCaller  bool
	AddSource   bool
	TimeFormat  string
	Environment string
	Service     string
}

var (
	globalLogger *Logger
	mu           sync.RWMutex
)

// New creates a new logger with the given configuration
func New(config Config) *Logger {
	var handler slog.Handler

	opts := &slog.HandlerOptions{
		Level:     convertLevel(config.Level),
		AddSource: config.AddSource,
	}

	if config.Output == nil {
		config.Output = os.Stdout
	}

	switch strings.ToLower(config.Format) {
	case "json":
		handler = slog.NewJSONHandler(config.Output, opts)
	default:
		handler = slog.NewTextHandler(config.Output, opts)
	}

	// Add default fields if configured
	if config.Environment != "" || config.Service != "" {
		attrs := make([]slog.Attr, 0, 3)
		if config.Service != "" {
			attrs = append(attrs, slog.String("service", config.Service))
		}
		if config.Environment != "" {
			attrs = append(attrs, slog.String("environment", config.Environment))
		}
		handler = handler.WithAttrs(attrs)
	}

	if config.Component != "" {
		handler = handler.WithGroup(config.Component)
	}

	return &Logger{
		Logger:    slog.New(handler),
		level:     config.Level,
		component: config.Component,
	}
}

// GetGlobalLogger returns the global logger
func GetGlobalLogger() *Logger {
	mu.RLock()
	gl := globalLogger
	mu.RUnlock()
	if gl != nil {
		return gl
	}
	// Slow path
	mu.Lock()
	defer mu.Unlock()
	if globalLogger == nil {
		globalLogger = New(DevelopmentConfig())
	}
	return globalLogger
}

// InitGlobalLogger initializes the global logger
func InitGlobalLogger(config Config) {
	mu.Lock()
	defer mu.Unlock()
	globalLogger = New(config)
}

// DevelopmentConfig returns a development-friendly configuration
func DevelopmentConfig() Config {
	return Config{
		Level:       LevelDebug,
		Format:      "text",
		Output:      os.Stdout,
		ShowCaller:  true,
		Service:     "benchmq",
		Environment: "development",
	}
}

// ProductionConfig returns a production-ready configuration
func ProductionConfig() Config {
	return Config{
		Level:       LevelInfo,
		Format:      "json",
		Output:      os.Stdout,
		ShowCaller:  false,
		Service:     "benchmq",
		Environment: "production",
	}
}

// Debug logs a debug message
func (l *Logger) Debug(msg string, attrs ...slog.Attr) {
	l.LogAttrs(context.Background(), slog.LevelDebug, msg, attrs...)
}

// Info logs an info message
func (l *Logger) Info(msg string, attrs ...slog.Attr) {
	l.LogAttrs(context.Background(), slog.LevelInfo, msg, attrs...)
}

// Warn logs a warning message
func (l *Logger) Warn(msg string, attrs ...slog.Attr) {
	l.LogAttrs(context.Background(), slog.LevelWarn, msg, attrs...)
}

// Error logs an error message
func (l *Logger) Error(msg string, attrs ...slog.Attr) {
	l.LogAttrs(context.Background(), slog.LevelError, msg, attrs...)
}

// Fatal logs a fatal message and exits
func (l *Logger) Fatal(msg string, attrs ...slog.Attr) {
	l.LogAttrs(context.Background(), slog.LevelError, msg, attrs...)
	os.Exit(1)
}

// Debug logs a debug message using the global logger
func Debug(msg string, attrs ...slog.Attr) {
	GetGlobalLogger().Debug(msg, attrs...)
}

// Info logs an info message using the global logger
func Info(msg string, attrs ...slog.Attr) {
	GetGlobalLogger().Info(msg, attrs...)
}

// Warn logs a warning message using the global logger
func Warn(msg string, attrs ...slog.Attr) {
	GetGlobalLogger().Warn(msg, attrs...)
}

// Error logs an error message using the global logger
func Error(msg string, attrs ...slog.Attr) {
	GetGlobalLogger().Error(msg, attrs...)
}

// Fatal logs a fatal message using the global logger and exits
func Fatal(msg string, attrs ...slog.Attr) {
	GetGlobalLogger().Fatal(msg, attrs...)
}

func convertLevel(level LogLevel) slog.Level {
	switch level {
	case LevelDebug:
		return slog.LevelDebug
	case LevelInfo:
		return slog.LevelInfo
	case LevelWarn:
		return slog.LevelWarn
	case LevelError:
		return slog.LevelError
	case LevelFatal:
		return slog.LevelError
	default:
		return slog.LevelInfo
	}
}

// NewBenchmarkLogger creates a component-specific logger for benchmark operations
func NewBenchmarkLogger(component string) *Logger {
	global := GetGlobalLogger()

	// Create a new logger with component group
	handler := global.Handler().WithGroup(component)

	return &Logger{
		Logger:    slog.New(handler),
		level:     global.level,
		component: component,
	}
}

// LogClientConnection logs client connection events
func (l *Logger) LogClientConnection(clientID string, attrs ...slog.Attr) {
	baseAttrs := []slog.Attr{
		ClientID(clientID),
	}
	baseAttrs = append(baseAttrs, attrs...)

	l.LogAttrs(context.Background(), slog.LevelInfo, "connected", baseAttrs...)
}

// LogPublish logs PUBLISH packet details
func (l *Logger) LogPublish(clientID, topic string, qos int, attrs ...slog.Attr) {
	baseAttrs := []slog.Attr{
		slog.String("client_id", clientID),
		slog.String("topic", topic),
		slog.Int("qos", qos),
	}
	baseAttrs = append(baseAttrs, attrs...)

	l.LogAttrs(context.Background(), slog.LevelInfo, "published", baseAttrs...)
}

// LogSubscribe logs SUBSCRIBE packet details
func (l *Logger) LogSubscribe(clientID, topic string, qos int, attrs ...slog.Attr) {
	baseAttrs := []slog.Attr{
		ClientID(clientID),
		slog.String("topic", topic),
		slog.Int("qos", qos),
	}
	baseAttrs = append(baseAttrs, attrs...)

	l.LogAttrs(context.Background(), slog.LevelInfo, "subscribed", baseAttrs...)
}

// ClientID creates a client_id attribute
func ClientID(clientID string) slog.Attr {
	return slog.String("client_id", clientID)
}

// State creates a state attribute
func State(state string) slog.Attr {
	return slog.String("state", state)
}

// String creates a string attribute
func String(key, value string) slog.Attr {
	return slog.String(key, value)
}

// Int creates an int attribute
func Int(key string, value int) slog.Attr {
	return slog.Int(key, value)
}

// Bool creates a bool attribute
func Bool(key string, value bool) slog.Attr {
	return slog.Bool(key, value)
}

// Float creates a float attribute
func Float(key string, value float64) slog.Attr {
	return slog.Float64(key, value)
}

// Any creates an attribute with any value
func Any(key string, value any) slog.Attr {
	return slog.Any(key, value)
}

// ErrorAttr creates an error attribute
func ErrorAttr(err error) slog.Attr {
	return slog.String("error", err.Error())
}
