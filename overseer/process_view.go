package main

import (
	"sort"

	"github.com/xackery/overseer/share/slog"
	"github.com/xackery/wlk/walk"
)

type ProcessViewEntry struct {
	Icon    *walk.Bitmap
	Name    string
	PID     string
	Status  string
	Uptime  string
	checked bool
}

type ProcessView struct {
	gui *Gui
	walk.TableModelBase
	walk.SorterBase
	sortColumn int
	sortOrder  walk.SortOrder
	items      []*ProcessViewEntry
}

func NewProcessView(gui *Gui) *ProcessView {
	m := new(ProcessView)
	m.gui = gui
	m.ResetRows()
	return m
}

func newProcessViewStyler(gui *Gui) *fileViewStyler {
	return &fileViewStyler{
		gui: gui,
	}
}

// Called by the TableView from SetModel and every time the model publishes a
// RowsReset event.
func (m *ProcessView) RowCount() int {
	return len(m.items)
}

// Called by the TableView when it needs the text to display for a given cell.
func (m *ProcessView) Value(row, col int) interface{} {
	if row < 0 || row >= len(m.items) {
		return nil
	}
	item := m.items[row]

	switch col {
	case -1:
		return nil
	case 0:
		return ""
	case 1:
		return item.Name
	case 2:
		return item.PID
	case 3:
		return item.Status
	case 4:
		return item.Uptime
	}

	slog.Printf("invalid col: %d\n", col)
	return nil
}

// Called by the TableView to retrieve if a given row is checked.
func (m *ProcessView) Checked(row int) bool {
	return m.items[row].checked
}

// Called by the TableView when the user toggled the check box of a given row.
func (m *ProcessView) SetChecked(row int, checked bool) error {
	m.items[row].checked = checked

	return nil
}

// Called by the TableView to sort the model.
func (m *ProcessView) Sort(col int, order walk.SortOrder) error {
	m.sortColumn, m.sortOrder = col, order

	sort.SliceStable(m.items, func(i, j int) bool {
		a, b := m.items[i], m.items[j]

		c := func(ls bool) bool {
			if m.sortOrder == walk.SortAscending {
				return ls
			}

			return !ls
		}

		switch m.sortColumn {
		case -1:
			return false
		case 0:
			return c(a.Name < b.Name)
		case 1:
			return c(a.Name < b.Name)
		case 2:
			return c(a.PID < b.PID)
		case 3:
			return c(a.Status < b.Status)
		case 4:
			return c(a.Uptime < b.Uptime)
		}

		slog.Printf("invalid sort col: %d", m.sortColumn)
		return false
	})

	return m.SorterBase.Sort(col, order)
}

func (m *ProcessView) ResetRows() {
	m.items = nil

	m.PublishRowsReset()

	m.Sort(m.sortColumn, m.sortOrder)
}

func (m *ProcessView) SetItems(items []*ProcessViewEntry) {
	m.items = items

	m.PublishRowsReset()

	m.Sort(m.sortColumn, m.sortOrder)
}

type fileViewStyler struct {
	gui *Gui
}

func (fv *fileViewStyler) StyleCell(style *walk.CellStyle) {
	if style.Col() != 0 {
		return
	}

	gui := fv.gui

	if style.Row() >= len(gui.procView.items) {
		return
	}

	item := gui.procView.items[style.Row()]
	if item == nil {
		slog.Printf("item %d is nil\n", style.Row())
		return
	}

	if item.Icon == nil {
		//slog.Printf("item %d icon is nil\n", style.Row())
		return
	}

	style.Image = item.Icon

	/* canvas := style.Canvas()
	if canvas == nil {
		return
	}
	bounds := style.Bounds()
	bounds.X += 2
	bounds.Y += 2
	bounds.Width = 16
	bounds.Height = 16
	err := canvas.DrawBitmapPartWithOpacityPixels(item.Icon, bounds, walk.Rectangle{X: 0, Y: 0, Width: 16, Height: 16}, 127)
	if err != nil {
		slog.Printf("failed to draw bitmap: %s\n", err.Error())
	} */

	/*

		switch style.Col() {
		case 1:
			if canvas := style.Canvas(); canvas != nil {
				bounds := style.Bounds()
				bounds.X += 2
				bounds.Y += 2
				bounds.Width = int((float64(bounds.Width) - 4) / 5 * float64(len(item.Bar)))
				bounds.Height -= 4
				canvas.DrawBitmapPartWithOpacity(barBitmap, bounds, walk.Rectangle{0, 0, 100 / 5 * len(item.Bar), 1}, 127)

				bounds.X += 4
				bounds.Y += 2
				canvas.DrawText(item.Bar, tv.Font(), 0, bounds, walk.TextLeft)
			}

		case 2:
			if item.Baz >= 900.0 {
				style.TextColor = walk.RGB(0, 191, 0)
				style.Image = goodIcon
			} else if item.Baz < 100.0 {
				style.TextColor = walk.RGB(255, 0, 0)
				style.Image = badIcon
			}

		case 3:
			if item.Quux.After(time.Now().Add(-365 * 24 * time.Hour)) {
				style.Font = boldFont
			}
		}
	*/
}
