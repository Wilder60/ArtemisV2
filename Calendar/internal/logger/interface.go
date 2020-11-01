package logger

// The Logger interface will define
type Logger interface {
	Debug(string)
	Info(string)
	Warn(string)
	Error(string)
	Fatal(string)
	Panic(string)
}
