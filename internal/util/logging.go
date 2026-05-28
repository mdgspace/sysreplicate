package util

import (
	"io"
	"log"
	"os"
)

type Level int

const (
	DebugLevel Level = iota
	InfoLevel
	WarnLevel
	ErrorLevel
)

type Logger struct {
	inner *log.Logger
	level Level
}

func NewLogger(level Level) *Logger {
	return &Logger{
		inner: log.New(os.Stderr, "", log.LstdFlags),
		level: level,
	}
}

func (l *Logger) SetOutput(w io.Writer) {
	l.inner.SetOutput(w)
}

func (l *Logger) Debug(format string, v ...interface{}) {
	if l.level <= DebugLevel {
		l.inner.Printf("[DEBUG] "+format, v...)
	}
}

func (l *Logger) Info(format string, v ...interface{}) {
	if l.level <= InfoLevel {
		l.inner.Printf("[INFO] "+format, v...)
	}
}

func (l *Logger) Warn(format string, v ...interface{}) {
	if l.level <= WarnLevel {
		l.inner.Printf("[WARN] "+format, v...)
	}
}

func (l *Logger) Error(format string, v ...interface{}) {
	if l.level <= ErrorLevel {
		l.inner.Printf("[ERROR] "+format, v...)
	}
}
