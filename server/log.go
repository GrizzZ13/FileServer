package server

import "fmt"

type Logger struct {
}

func NewLogger() *Logger {
	return &Logger{}
}

func (l *Logger) Log(a ...any) {
	fmt.Println(a...)
}

func (l *Logger) Error(a ...any) {
	fmt.Println(a...)
}
