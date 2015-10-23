package log

import (
	"fmt"
	"io"
	"os"

	pretty "github.com/deis/pkg/prettyprint"
)

var Stdout io.Writer = os.Stdout
var Stderr io.Writer = os.Stderr

var IsDebugging = false

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

// CleanExit prints a message and then exits with 0.
func CleanExit(format string, v ...interface{}) {
	Info(format, v...)
	os.Exit(0)
}

// Err prints an error message. It does not cause an exit.
func Err(format string, v ...interface{}) {
	fmt.Fprint(Stderr, pretty.Colorize("{{.Red}}[ERROR]{{.Default}} "))
	fmt.Fprintf(Stderr, appendNewLine(format), v...)
}

// Info prints a message.
func Info(format string, v ...interface{}) {
	fmt.Fprint(Stderr, pretty.Colorize("{{.Green}}--->{{.Default}} "))
	fmt.Fprintf(Stderr, appendNewLine(format), v...)
}

func Debug(msg string, v ...interface{}) {
	if IsDebugging {
		fmt.Fprint(Stderr, pretty.Colorize("{{.Cyan}}[DEBUG]{{.Default}} "))
		Msg(msg, v...)
	}
}

func Warn(format string, v ...interface{}) {
	fmt.Fprint(Stderr, pretty.Colorize("{{.Yellow}}[WARN]{{.Default}} "))
	Msg(format, v...)
}

func appendNewLine(format string) string {
	return format + "\n"
}
