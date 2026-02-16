package logger

import (
	"fmt"
	"io"
)

type Level int

const (
	LevelDebug Level = iota // 0 - Very verbose
	LevelInfo               // 1 - Normal
	LevelWarn               // 2 - Warnings
	LevelError              // 3 - Errors only
)

type Logger struct {
	Writer io.Writer
	Level  Level
}

func New(w io.Writer, level Level) *Logger {
	return &Logger{Writer: w, Level: level}
}

func (l *Logger) Debug(format string, args ...any) {
	if l.Level <= LevelDebug {
		fmt.Fprintf(l.Writer, "[DEBUG] "+format, args...)
	}
}

func (l *Logger) Info(format string, args ...any) {
	if l.Level <= LevelInfo {
		fmt.Fprintf(l.Writer, format, args...)
	}
}

func (l *Logger) Warn(format string, args ...any) {
	if l.Level <= LevelWarn {
		fmt.Fprintf(l.Writer, "[!] "+format, args...)
	}
}

func (l *Logger) Error(format string, args ...any) {
	if l.Level <= LevelError {
		fmt.Fprintf(l.Writer, "[ERROR] "+format, args...)
	}
}
