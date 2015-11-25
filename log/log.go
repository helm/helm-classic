package log

import (
	"fmt"
	"io"
	"log"
	"os"

	pretty "github.com/deis/pkg/prettyprint"
)

// Stdout is the logging destination for normal messages.
var Stdout io.Writer = os.Stdout

// Stderr is the logging destination for error messages.
var Stderr io.Writer = os.Stderr

// Stdin is the input alternative for logging.
//
// Applications that take console-like input should use this.
var Stdin io.Reader = os.Stdin

// IsDebugging toggles whether or not to enable debug output and behavior.
var IsDebugging = false

// ErrorState denotes if application is in an error state.
var ErrorState = false

// New creates a *log.Logger that writes to this source.
func New() *log.Logger {
	ll := log.New(Stdout, pretty.Colorize("{{.Yellow}}--->{{.Default}} "), 0)
	return ll
}

// Msg passes through the formatter, but otherwise prints exactly as-is.
//
// No prettification.
func Msg(format string, v ...interface{}) {
	fmt.Fprintf(Stdout, appendNewLine(format), v...)
}

// Die prints an error and then call os.Exit(1).
func Die(format string, v ...interface{}) {
	Err(format, v...)
	if IsDebugging {
		panic(fmt.Sprintf(format, v...))
	}
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
	ErrorState = true
}

// Info prints a green-tinted message.
func Info(format string, v ...interface{}) {
	fmt.Fprint(Stderr, pretty.Colorize("{{.Green}}--->{{.Default}} "))
	fmt.Fprintf(Stderr, appendNewLine(format), v...)
}

// Debug prints a cyan-tinted message if IsDebugging is true.
func Debug(format string, v ...interface{}) {
	if IsDebugging {
		fmt.Fprint(Stderr, pretty.Colorize("{{.Cyan}}[DEBUG]{{.Default}} "))
		fmt.Fprintf(Stderr, appendNewLine(format), v...)
	}
}

// Warn prints a yellow-tinted warning message.
func Warn(format string, v ...interface{}) {
	fmt.Fprint(Stderr, pretty.Colorize("{{.Yellow}}[WARN]{{.Default}} "))
	fmt.Fprintf(Stderr, appendNewLine(format), v...)
}

func appendNewLine(format string) string {
	return format + "\n"
}
