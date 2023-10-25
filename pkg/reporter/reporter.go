package reporter

import (
	"strings"
	"sync"
	"time"
)

var (
	mu             sync.RWMutex
	apps           = make(map[string]*App)
	SendUpdateChan = make(chan bool, 1000)
)

type App struct {
	Status AppState
	PID    int
	start  time.Time
}

func (a *App) Uptime() string {

	// return format string to be nearest whole number
	uptime := time.Since(a.start)
	if uptime < time.Second {
		return "0s"
	}
	if uptime < time.Minute {
		return uptime.Round(time.Second).String()
	}
	if uptime < time.Hour {
		return uptime.Round(time.Minute).String()
	}
	if uptime < 24*time.Hour {
		return uptime.Round(time.Hour).String()
	}

	return uptime.Round(24 * time.Hour).String()
}

type AppState int

const (
	AppStateUnknown AppState = iota
	AppStateStarting
	AppStateRunning
	AppStateStopped
	AppStateRestarting
	AppStateSleeping
	AppStateErroring
)

type AppStateReport struct {
	States         map[string]AppState
	ZoneTotal      int
	ZoneUnknown    int
	ZoneStarting   int
	ZoneRunning    int
	ZoneStopped    int
	ZoneRestarting int
	ZoneSleeping   int
	ZoneErroring   int
}

// ZoneUpdate updates the status of a zone.
func SetAppState(name string, status AppState) {
	mu.Lock()
	defer mu.Unlock()

	app, ok := apps[name]
	if !ok {
		app = &App{
			start: time.Now(),
		}
		apps[name] = app
	}

	isUpdate := false
	if app.Status != status {
		isUpdate = true
	}
	app.Status = status
	if isUpdate {
		SendUpdateChan <- true
	}
}

func SetAppPID(name string, pid int) {
	mu.Lock()
	defer mu.Unlock()
	app, ok := apps[name]
	if !ok {
		app = &App{
			start: time.Now(),
		}
		apps[name] = app
	}
	if pid == 0 {
		app.start = time.Now()
	}
	isUpdate := false
	if app.PID != pid {
		isUpdate = true
	}
	app.PID = pid
	if isUpdate {
		SendUpdateChan <- true
	}
}

// AppPtr is used by windows for showing a GUI of apps
func AppPtr() map[string]*App {
	mu.RLock()
	defer mu.RUnlock()
	return apps
}

// AppStates returns the status of all zones
func AppStates() *AppStateReport {
	mu.RLock()
	defer mu.RUnlock()

	result := &AppStateReport{
		States: make(map[string]AppState),
	}
	for k, v := range apps {
		if strings.Contains(k, "zone") {
			result.ZoneTotal++
			switch v.Status {
			case AppStateUnknown:
				result.ZoneUnknown++
			case AppStateStarting:
				result.ZoneStarting++
			case AppStateRunning:
				result.ZoneRunning++
			case AppStateStopped:
				result.ZoneStopped++
			case AppStateRestarting:
				result.ZoneRestarting++
			case AppStateSleeping:
				result.ZoneSleeping++
			case AppStateErroring:
				result.ZoneErroring++
			}
			continue
		}
		result.States[k] = v.Status
	}
	return result
}

func AppStateString(in AppState) string {
	switch in {
	case AppStateUnknown:
		return "Unknown"
	case AppStateStarting:
		return "Starting"
	case AppStateRunning:
		return "Running"
	case AppStateStopped:
		return "Stopped"
	case AppStateRestarting:
		return "Restarting"
	case AppStateSleeping:
		return "Sleeping"
	case AppStateErroring:
		return "Erroring"
	}
	return "Unknown"
}
