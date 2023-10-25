//go:build windows
// +build windows

package main

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/xackery/overseer/pkg/handler"
	"github.com/xackery/overseer/pkg/reporter"
	"github.com/xackery/overseer/pkg/signal"
	"github.com/xackery/overseer/pkg/slog"
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
		MinSize: cpl.Size{Width: 165, Height: 200},
		Size:    cpl.Size{Width: 165, Height: 300},
		Layout:  cpl.Grid{Columns: 2},
		Visible: false,
		Name:    "overseer",
		Children: []cpl.Widget{
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
				ContextMenuItems: []cpl.MenuItem{
					cpl.Action{
						Text: "End task",
					},
					cpl.Separator{},
					cpl.Action{
						Text: "Open log",
					},
					cpl.Action{
						Text: "Properties",
					},
				},
				Columns: []cpl.TableViewColumn{
					{Name: "Name", Width: 100},
					{Name: "PID", Width: 50},
					{Name: "Status", Width: 70},
					{Name: "Uptime", Width: 70},
				},
			},
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

	gui.SubscribeClose(func(cancelled *bool, reason byte) {
		if ctx.Err() == nil {
			*cancelled = true
			fmt.Println("Got close message")
			handler.WindowCloseInvoke(cancelled, reason)
			cancel()
			return
		}
		fmt.Println("Officially exiting")
	})

	slog.AddHandler(gui.Logf)

	return gui, nil
}

func (gui *Gui) Run() int {
	if gui.mw == nil {
		return 1
	}
	gui.mw.SetVisible(true)
	return gui.mw.Run()
}

func (gui *Gui) SubscribeClose(fn func(cancelled *bool, reason byte)) {
	if gui.mw == nil {
		return
	}
	gui.mw.Closing().Attach(fn)
}

func (gui *Gui) Close() error {
	if gui.ctx.Err() == nil {
		return nil
	}

	if gui.mw == nil {
		return nil
	}

	walk.App().Exit(0)
	return nil
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

func runWindows(ctx context.Context, gui *Gui) error {
	items := []*ProcessViewEntry{}

	go func() {
		for {
			fmt.Println("listening")
			select {
			case <-ctx.Done():
				return
			case <-reporter.SendUpdateChan:
			}

			isFirstRun := len(items) == 0

			apps := reporter.AppPtr()
			fmt.Println("Got update", len(apps))
			for name, app := range apps {
				if app == nil {
					continue
				}
				isFound := false
				for _, item := range items {
					if item.Name != name {
						continue
					}
					item.PID = fmt.Sprintf("%d", app.PID)
					item.Status = reporter.AppStateString(app.Status)
					item.Uptime = app.Uptime()
					isFound = true
					break
				}
				if !isFound {
					items = append(items, &ProcessViewEntry{
						Name:   name,
						PID:    fmt.Sprintf("%d", app.PID),
						Status: reporter.AppStateString(app.Status),
						Uptime: app.Uptime(),
					})
				}

				if isFirstRun {
					gui.SetProcessViewItems(items)
				}
				gui.procView.PublishRowsReset()

			}
		}
	}()
	go func() {
		<-ctx.Done()
		fmt.Println("Doing clean up process...")
		gui.SetTitle("Shutting down... Please wait, ensuring all processes are exiting!")
		signal.Cancel()
		signal.WaitWorker()
		gui.Close()
		fmt.Println("Done, exiting")
		os.Exit(0)
	}()

	errCode := gui.Run()
	if errCode != 0 {
		fmt.Println("Failed to run:", errCode)
		os.Exit(1)
	}

	return nil
}

// Logf logs a message to the gui
func (gui *Gui) Logf(format string, a ...interface{}) {
	if gui == nil {
		return
	}

	line := fmt.Sprintf(format, a...)
	if strings.Contains(line, "\n") {
		line = "  " + line[0:strings.Index(line, "\n")]
	}
	gui.statusBar.SetText(line)

	//convert \n to \r\n
	//format = strings.ReplaceAll(format, "\n", "\r\n")
	//gui.log.AppendText(fmt.Sprintf(format, a...))
}
