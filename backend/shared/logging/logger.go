package logging

import "net/http"

type Logger interface {
	Debug(msg string)
	Warn(msg string)
	Error(msg string)

	Middleware(next http.Handler) http.Handler
}
