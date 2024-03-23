package main

import (
	"context"
	"fmt"
	"strings"
	"sync"

	"github.com/xackery/wlk/cpl"
	"github.com/xackery/wlk/walk"
)

type Gui struct {
	ctx         context.Context
	cancel      context.CancelFunc
	mw          *walk.MainWindow
	splash      *walk.ImageView
	isAutoPatch *walk.CheckBox
	isAutoPlay  *walk.CheckBox
	patchButton *walk.PushButton
	playButton  *walk.PushButton
	progress    *walk.ProgressBar
	statusBar   *walk.StatusBarItem
}

var (
	isAutoMode bool
	mu         sync.RWMutex
)

// NewMainWindow creates a new main window
func NewMainWindow(ctx context.Context, cancel context.CancelFunc, version string) (*Gui, error) {
	gui := &Gui{
		ctx:    ctx,
		cancel: cancel,
	}

	var err error
	cmw := cpl.MainWindow{
		Title:   "diagnose v" + version,
		MinSize: cpl.Size{Width: 405, Height: 371},
		Layout:  cpl.VBox{},
		Visible: false,
		Name:    "diagnose",
		Children: []cpl.Widget{
			cpl.HSplitter{Children: []cpl.Widget{
				cpl.VSplitter{Children: []cpl.Widget{
					cpl.Label{Text: "Files"},
					cpl.PushButton{
						Text: "Restart",
					},
				}},
			}},
		},
		AssignTo: &gui.mw,
		StatusBarItems: []cpl.StatusBarItem{
			{
				AssignTo: &gui.statusBar,
				Text:     "Ready",
				OnClicked: func() {
					fmt.Println("status bar clicked")
				},
			},
		},
	}
	err = cmw.Create()
	if err != nil {
		return nil, fmt.Errorf("create main window: %w", err)
	}

	return gui, nil
}

func (g *Gui) Run() int {
	if g.mw == nil {
		return 1
	}
	g.mw.SetVisible(true)
	return g.mw.Run()
}

func (g *Gui) SubscribeClose(fn func(cancelled *bool, reason byte)) {
	if g.mw == nil {
		return
	}
	g.mw.Closing().Attach(fn)
}

// Logf logs a message to the gui
func (g *Gui) Logf(format string, a ...interface{}) {
	if g.mw == nil {
		return
	}

	line := fmt.Sprintf(format, a...)
	if strings.Contains(line, "\n") {
		line = line[0:strings.Index(line, "\n")]
	}
	g.statusBar.SetText(line)
}

func (g *Gui) Close() error {
	if g.mw == nil {
		return nil
	}
	return g.mw.Close()
}

func (g *Gui) SetTitle(title string) {
	if g.mw == nil {
		return
	}
	g.mw.SetTitle(title)
}
