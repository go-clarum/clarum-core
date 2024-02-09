package logging

import (
	"fmt"
	"log"
	"log/slog"
	"os"
)

var internalLogger = log.New(os.Stdout, "", log.LstdFlags|log.Lmicroseconds)
var defaultLogger = NewLogger(slog.LevelInfo, "")

type Logger struct {
	level  slog.Level
	prefix string
}

func NewLogger(level slog.Level, prefix string) *Logger {
	return &Logger{level, prefix}
}

func (logger *Logger) Info(message string) {
	if logger.level <= slog.LevelInfo {
		internalLogger.Println("INFO " + logger.prefix + message)
	}
}

func (logger *Logger) Infof(format string, a ...any) {
	if logger.level <= slog.LevelInfo {
		internalLogger.Println("INFO " + logger.prefix + fmt.Sprintf(format, a...))
	}
}

func (logger *Logger) Debug(message string) {
	if logger.level <= slog.LevelDebug {
		internalLogger.Println("DEBUG " + logger.prefix + message)
	}
}

func (logger *Logger) Debugf(format string, a ...any) {
	if logger.level <= slog.LevelDebug {
		internalLogger.Println("DEBUG " + logger.prefix + fmt.Sprintf(format, a...))
	}
}

func (logger *Logger) Warn(message string) {
	if logger.level <= slog.LevelWarn {
		internalLogger.Println("WARN " + logger.prefix + message)
	}
}

func (logger *Logger) Warnf(format string, a ...any) {
	if logger.level <= slog.LevelWarn {
		internalLogger.Println("WARN " + logger.prefix + fmt.Sprintf(format, a...))
	}
}

func (logger *Logger) Error(message string) {
	if logger.level <= slog.LevelError {
		internalLogger.Println("ERROR " + logger.prefix + message)
	}
}

func (logger *Logger) Errorf(format string, a ...any) {
	if logger.level <= slog.LevelError {
		internalLogger.Println("ERROR " + logger.prefix + fmt.Sprintf(format, a...))
	}
}

func (logger *Logger) Prefix() string {
	return logger.prefix
}

// calls on the default logger

func Info(message string) {
	defaultLogger.Info(message)
}

func Infof(format string, a ...any) {
	defaultLogger.Infof(format, a...)
}

func Debug(message string) {
	defaultLogger.Debug(message)
}

func Debugf(format string, a ...any) {
	defaultLogger.Debugf(format, a...)
}

func Warn(message string) {
	defaultLogger.Warn(message)
}

func Warnf(format string, a ...any) {
	defaultLogger.Warnf(format, a...)
}

func Error(message string) {
	defaultLogger.Error(message)
}

func Errorf(format string, a ...any) {
	defaultLogger.Errorf(format, a...)
}
