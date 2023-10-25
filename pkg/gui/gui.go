package gui

var (
	gui GuiWrapper
)

// GuiWrapper is a wrapper for gui functions
type GuiWrapper interface {
	Run() int
	Close() error
	SetTitle(title string)
}

// New sets a gui instance to the gui wrapper
func New(g GuiWrapper) {
	gui = g
}

func Run() int {
	if gui == nil {
		return 1
	}
	return gui.Run()
}

func Close() error {
	if gui == nil {
		return nil
	}
	return gui.Close()
}

func SetTitle(title string) {
	if gui == nil {
		return
	}
	gui.SetTitle(title)
}
