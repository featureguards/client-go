package logger

type LogLevel int

const (
	Error LogLevel = iota
	Warn  LogLevel = iota + 1
	Info  LogLevel = iota + 2
	Debug LogLevel = iota + 3
)
