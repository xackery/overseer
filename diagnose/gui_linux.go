//go:build !windows
// +build !windows

package main

import (
	"context"
)

type Gui struct {
	ctx    context.Context
	cancel context.CancelFunc
}

var ()

// NewMainWindow creates a new main window
func NewMainWindow(ctx context.Context, cancel context.CancelFunc, version string) (*Gui, error) {
	gui := &Gui{
		ctx:    ctx,
		cancel: cancel,
	}

	return gui, nil
}

func (g *Gui) Run() int {
	return 0
}

func (g *Gui) SubscribeClose(fn func(cancelled *bool, reason byte)) {
}

// Logf logs a message to the gui
func (g *Gui) Logf(format string, a ...interface{}) {
}

func (g *Gui) Close() error {
	return nil
}

func (g *Gui) SetTitle(title string) {
}
