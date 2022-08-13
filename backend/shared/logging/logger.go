package logging

type Logger interface {
	Debug(msg string)
	Warn(msg string)
	Error(msg string)
}
