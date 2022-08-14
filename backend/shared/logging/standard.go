package logging

import (
	"log"
	"net/http"
)

type DefaultLogger struct {
	log *log.Logger
}

func NewDefaultLogger(l *log.Logger) *DefaultLogger {
	return &DefaultLogger{l}
}

func (l *DefaultLogger) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		l.log.Printf("%s - %s", r.Method, r.URL.Path)

		next.ServeHTTP(rw, r)
	})
}

func (l *DefaultLogger) Debug(msg string) {
	l.log.Println("[INFO]", msg)
}

func (l *DefaultLogger) Warn(msg string) {
	l.log.Println("[WARN]", msg)
}

func (l *DefaultLogger) Error(msg string) {
	l.log.Println("[ERROR]", msg)
}
