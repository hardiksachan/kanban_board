package logging

import "log"

type DefaultLogger struct {
	l *log.Logger
}

func NewDefaultLogger(l *log.Logger) *DefaultLogger {
	return &DefaultLogger{l}
}

func (l DefaultLogger) Debug(msg string) {
	l.l.Println("[INFO]", msg)
}

func (l DefaultLogger) Warn(msg string) {
	l.l.Println("[WARN]", msg)
}

func (l DefaultLogger) Error(msg string) {
	l.l.Println("[ERROR]", msg)
}
