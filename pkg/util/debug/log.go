package debug

import (
	"fmt"
	"log"
	"runtime"
	"strings"
)

const symbol = "✦✦✦DEBUG✦✦✦"

type Logger struct {
	includeCallStack         bool
	includeFormatPackageName bool
	callerDepth              int
}

func NewLogger() *Logger {
	return &Logger{
		callerDepth: 1,
	}
}

// func (l *Logger) WithCallStack() *Logger {
// 	l.includeCallStack = true
// 	l.callerDepth += 1
// 	return l
// }

func (l *Logger) WithFormatPackageName() *Logger {
	l.includeFormatPackageName = true
	return l
}

func (l *Logger) Printf(format string, args ...interface{}) {
	_, file, line, _ := runtime.Caller(l.callerDepth)
	parts := strings.Split(file, "/")

	pkg := "main"
	if len(parts) >= 2 {
		pkg = parts[len(parts)-2]
	}
	entry := symbol
	if l.includeFormatPackageName {
		entry += fmt.Sprintf("[%s:%d][%s] %s", file, line, pkg, fmt.Sprintf(format, args...))
	} else {
		entry += fmt.Sprintf("[%s:%d] %s", file, line, fmt.Sprintf(format, args...))
	}

	if l.includeCallStack {
		buf := make([]byte, 1024)
		runtime.Stack(buf, false)
		entry += "\nCall stack:\n" + strings.TrimSpace(string(buf))
	}

	log.Println(entry)
}

func (l *Logger) Println(args ...interface{}) {
	_, file, line, _ := runtime.Caller(l.callerDepth)
	parts := strings.Split(file, "/")

	pkg := "main"
	if len(parts) >= 2 {
		pkg = parts[len(parts)-2]
	}
	entry := symbol
	if l.includeFormatPackageName {
		entry += fmt.Sprintf("[%s:%d][%s] %s", file, line, pkg, fmt.Sprint(args...))
	} else {
		entry += fmt.Sprintf("[%s:%d] %s", file, line, fmt.Sprint(args...))
	}

	if l.includeCallStack {
		buf := make([]byte, 1024)
		runtime.Stack(buf, false)
		entry += "\nCall stack:\n" + strings.TrimSpace(string(buf))
	}

	log.Println(entry)
}
