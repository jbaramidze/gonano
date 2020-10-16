package logger

import (
	"fmt"
	"os"
)

// L Logger
var L Logger

// Logger ss
type Logger struct {
	f *os.File
}

// Log to log file
func (l Logger) Log(s string) {
	l.f.WriteString(s + "\n")
}

func (l Logger) close() {
	l.f.Close()
}

func init() {
	f, err := os.OpenFile("log", os.O_TRUNC|os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		fmt.Println("Error", err)
		os.Exit(-1)
	}
	L = Logger{f}
}
