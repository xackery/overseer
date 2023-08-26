package reporter

import (
	"strings"
	"sync"
)

var (
	mu             sync.RWMutex
	appStates      = make(map[string]AppState)
	SendUpdateChan = make(chan bool, 1000)
)

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

	appStates[name] = status
	SendUpdateChan <- true
}

// AppStates returns the status of all zones
func AppStates() *AppStateReport {
	mu.RLock()
	defer mu.RUnlock()

	result := &AppStateReport{
		States: make(map[string]AppState),
	}
	for k, v := range appStates {
		if strings.Contains(k, "zone") {
			result.ZoneTotal++
			switch v {
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
		result.States[k] = v
	}
	return result
}
