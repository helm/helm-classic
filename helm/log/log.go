package log

import (
	"fmt"
	"os"
)

var Stdout = os.Stdout
var Stderr = os.Stderr

// TODO: Wire this to the global --debug flag and then set this to false.
var IsDebugging = true

// Msg passes through the formatter, but otherwise prints exactly as-is.
//
// No prettification.
func Msg(format string, v ...interface{}) {
	fmt.Fprintf(Stdout, appendNewLine(format), v...)
}

// Die prints an error and then call os.Exit(1).
func Die(format string, v ...interface{}) {
	Err(format, v...)
	os.Exit(1)
}

// Err prints an error message. It does not cause an exit.
func Err(format string, v ...interface{}) {
	fmt.Fprintf(Stderr, appendNewLine(format), v...)
}

// Info prints a message.
func Info(format string, v ...interface{}) {
	fmt.Fprintf(Stderr, appendNewLine(format), v...)
}

func Debug(msg string, v ...interface{}) {
	if IsDebugging {
		Info(msg, v...)
	}
}

func Warn(format string, v ...interface{}) {
	Info(format, v...)
}

func appendNewLine(format string) string {
	return format + "\n"
}
