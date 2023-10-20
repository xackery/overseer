//go:build windows
// +build windows

package main

import (
	"context"
	"fmt"
	"strings"

	"github.com/xackery/overseer/share/slog"
	"github.com/xackery/wlk/cpl"
	"github.com/xackery/wlk/walk"
)

type Gui struct {
	ctx         context.Context
	cancel      context.CancelFunc
	mw          *walk.MainWindow
	statusBar   *walk.StatusBarItem
	table       *walk.TableView
	procView    *ProcessView
	procEntries []*ProcessViewEntry
}

// NewMainWindow creates a new main window
func NewMainWindow(ctx context.Context, cancel context.CancelFunc, version string) (*Gui, error) {
	gui := &Gui{
		ctx:    ctx,
		cancel: cancel,
	}

	var err error
	fvs := newProcessViewStyler(gui)
	gui.procView = NewProcessView(gui)
	cmw := cpl.MainWindow{
		Title:   "overseer v" + version,
		MinSize: cpl.Size{Width: 305, Height: 371},
		Size:    cpl.Size{Width: 365, Height: 371},
		Layout:  cpl.VBox{},
		Visible: false,
		Name:    "overseer",
		Children: []cpl.Widget{
			cpl.HSplitter{Children: []cpl.Widget{
				cpl.VSplitter{Children: []cpl.Widget{
					cpl.TableView{
						AssignTo:              &gui.table,
						Name:                  "tableView",
						AlternatingRowBG:      true,
						ColumnsOrderable:      true,
						MultiSelection:        false,
						Model:                 gui.procView,
						OnCurrentIndexChanged: gui.onTableSelect,
						StyleCell:             fvs.StyleCell,
						MinSize:               cpl.Size{Width: 360, Height: 0},
						Columns: []cpl.TableViewColumn{
							{Name: " ", Width: 30},
							{Name: "Name", Width: 160},
							{Name: "PID", Width: 40},
							{Name: "Status", Width: 60},
							{Name: "Uptime", Width: 60},
						},
					},
				}},
			}},
		},
		AssignTo: &gui.mw,
		StatusBarItems: []cpl.StatusBarItem{
			{
				AssignTo: &gui.statusBar,
				Text:     "  Ready",
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

	gui.SubscribeClose(func(canceled *bool, reason walk.CloseReason) {
		if ctx.Err() != nil {
			fmt.Println("Accepting exit")
			return
		}
		*canceled = true
		fmt.Println("Got close message")
		cancel()
	})

	return gui, nil
}

func (gui *Gui) Run() int {
	if gui.mw == nil {
		return 1
	}
	gui.mw.SetVisible(true)
	return gui.mw.Run()
}

func (gui *Gui) SubscribeClose(fn func(cancelled *bool, reason walk.CloseReason)) {
	if gui.mw == nil {
		return
	}
	gui.mw.Closing().Attach(fn)
}

// Logf logs a message to the gui
func (gui *Gui) Logf(format string, a ...interface{}) {
	if gui.mw == nil {
		return
	}

	line := fmt.Sprintf(format, a...)
	if strings.Contains(line, "\n") {
		line = line[0:strings.Index(line, "\n")]
	}
	gui.statusBar.SetText(line)
}

func (gui *Gui) LogClear() {
	if gui.mw == nil {
		return
	}
	gui.statusBar.SetText("")
}

func (gui *Gui) Close() error {
	if gui.mw == nil {
		return nil
	}
	return gui.mw.Close()
}

func (gui *Gui) SetTitle(title string) {
	if gui.mw == nil {
		return
	}
	gui.mw.SetTitle(title)
}

func (gui *Gui) onTableSelect() {
	if len(gui.procEntries) == 0 {
		slog.Printf("No files to open")
		return
	}

	if gui.table.CurrentIndex() < 0 || gui.table.CurrentIndex() >= len(gui.procEntries) {
		//slog.Printf("Invalid file index %d", gui.table.CurrentIndex())
		return
	}
	name := gui.procEntries[gui.table.CurrentIndex()].Name
	slog.Printf("Selected %s\n", name)
}

func (gui *Gui) SetProcessViewItems(items []*ProcessViewEntry) {
	if gui == nil {
		return
	}
	gui.procEntries = items
	gui.procView.SetItems(items)
}
