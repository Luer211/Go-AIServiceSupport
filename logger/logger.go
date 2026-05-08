package logger

import "log"

// Todo: 未引入Zap

type Logger interface {
	Info(msg string, fields ...any)
	Error(msg string, fields ...any)
}

type StdLogger struct{}

func New() Logger {
	return &StdLogger{}
}

func (l *StdLogger) Info(msg string, fields ...any) {
	log.Println(append([]any{"INFO", msg}, fields...)...)
}

func (l *StdLogger) Error(msg string, fields ...any) {
	log.Println(append([]any{"ERROR", msg}, fields...)...)
}
