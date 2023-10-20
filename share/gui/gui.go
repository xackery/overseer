package gui

var (
	gui GuiWrapper
)

// GuiWrapper is a wrapper for gui functions
type GuiWrapper interface {
	Run() int
	Logf(format string, a ...interface{})
	LogClear()
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

// Logf logs a message to the gui
func Logf(format string, a ...interface{}) {
	if gui == nil {
		return
	}
	gui.Logf(format, a...)
}

func LogClear() {
	if gui == nil {
		return
	}
	gui.LogClear()
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
