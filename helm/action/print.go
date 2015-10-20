package action

import (
	"fmt"
	"os"
)

// Die prints an error and then call os.Exit(1).
func Die(msg string, v ...interface{}) {
	Err(msg, v...)
	os.Exit(1)
}

// Err prints an error message. It does not cause an exit.
func Err(msg string, v ...interface{}) {
	fmt.Fprintf(os.Stderr, msg, v...)
	fmt.Fprintln(os.Stderr)
}

// Info prints a message.
func Info(msg string, v ...interface{}) {
	fmt.Fprintf(os.Stderr, msg, v...)
	fmt.Fprintln(os.Stderr)
}

func Warn(msg string, v ...interface{}) {
	Info(msg, v...)
}
