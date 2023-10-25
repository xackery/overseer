package flog

import (
	"fmt"
	"os"
)

var (
	w *os.File
)

// New creates a new file logger
func New(path string) error {
	var err error
	w, err = os.Create(path)
	if err != nil {
		return fmt.Errorf("create: %w", err)
	}
	return nil
}

// Printf prints a formatted string to the file
func Printf(format string, a ...interface{}) {
	if w == nil {
		return
	}

	fmt.Fprintf(w, format, a...)
}

// Println prints a string to the file
func Println(a ...interface{}) {
	if w == nil {
		return
	}

	fmt.Fprintln(w, a...)
}

// Close closes the file
func Close() error {
	if w == nil {
		return nil
	}
	return w.Close()
}
