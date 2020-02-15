package server

import (
	"fmt"
	"log"
	"os"
	"runtime"
)

type Logger interface {
	Error(msg ...interface{})
	Trace(msg ...interface{})
}

type StdLogger struct {
	err   *log.Logger
	trace *log.Logger
}

func NewStdLogger(prefix string) *StdLogger {
	return &StdLogger{
		err:   log.New(os.Stdout, prefix+"[ERROR]", log.Ldate|log.Ltime),
		trace: log.New(os.Stdout, prefix+"[TRACE]", log.Ldate|log.Ltime),
	}
}

func (l *StdLogger) Error(msg ...interface{}) {
	l.err.Println(msg...)
	l.err.Println(callstack()...)
}

func (l *StdLogger) Trace(msg ...interface{}) {
	l.trace.Println(msg...)
}

func callstack() []interface{} {
	var cs []interface{}
	cs = append(cs, "Callstack:\n") // start with a new line
	for skip := 0; ; skip++ {
		_, file, line, ok := runtime.Caller(skip)
		if !ok {
			break
		}
		cs = append(cs, fmt.Sprintf("%s:%d\n", file, line))
	}
	return cs
}
