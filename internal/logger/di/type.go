package di

type Logger interface {
	Debug(format string, a ...Field)
	Info(format string, a ...Field)
	Warn(format string, a ...Field)
	Error(format string, a ...Field)
}
type Field struct {
	Key string
	Val any
}
