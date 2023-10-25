//go:build !windows
// +build !windows

package main

import "context"

type Gui struct {
}

// NewMainWindow creates a new main window
func NewMainWindow(ctx context.Context, cancel context.CancelFunc, version string) (*Gui, error) {
	g := &Gui{}
	return g, nil
}

func (g *Gui) Close() error {
	return nil
}

func (g *Gui) Run() int {
	return 0
}

func runWindows(ctx context.Context, g *Gui) error {
	return nil
}

func (g *Gui) SetTitle(title string) {

}
