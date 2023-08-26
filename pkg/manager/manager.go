package manager

import (
	"context"
	"strings"
	"time"

	"github.com/xackery/overseer/pkg/reporter"
	"github.com/xackery/overseer/pkg/runner"
	"github.com/xackery/overseer/pkg/signal"
)

type manager struct {
	ctx           context.Context
	displayName   string
	exeName       string
	args          []string
	startDelay    time.Duration
	state         reporter.AppState
	restartCount  int
	lastError     string
	lastErrorAt   time.Time // When lastErrorAt hits 30 minutes, reset errorCount
	errorCooldown time.Time
	errorCount    int // When errorCount hits 3, set errorCooldown to 30 minutes
	doneChan      chan error
	outChan       chan string
}

func (e *manager) setState(state reporter.AppState) {
	e.state = state
	reporter.SetAppState(e.displayName, state)
}

// Manage is the main loop for the zone.
func Manage(displayName string, exeName string, args ...string) {
	go poll(displayName, exeName, args...)
}

func poll(displayName string, exeName string, args ...string) {
	signal.AddWorker()
	defer signal.FinishWorker()

	mgr := &manager{
		ctx:         signal.Ctx(),
		displayName: displayName,
		exeName:     exeName,
		args:        args,
		outChan:     make(chan string),
		lastError:   "none",
		doneChan:    make(chan error),
	}

	run := runner.New(mgr.outChan, mgr.doneChan, mgr.exeName, mgr.args...)
	for {
		select {
		case <-mgr.ctx.Done():
			mgr.setState(reporter.AppStateStopped)
			run.Stop()
			return
		default:
		}
		go run.Start()
		mgr.setState(reporter.AppStateStarting)

		parse(mgr)
	}
}

func parse(mgr *manager) {
	start := time.Now()
	for {
		select {
		case line := <-mgr.outChan:
			mgr.lineParse(line)

			//fmt.Printf("[zone %d] %s\n", mgr.port, line)
		case <-mgr.ctx.Done():
			return
		case <-mgr.doneChan:
			mgr.restartCount++
			//fmt.Printf("[zone %d] Exited after %s seconds, %d restarts. Last error: %s\n", mgr.port, time.Since(start).Round(time.Second), mgr.restartCount, mgr.lastError)
			if time.Since(start) > 3*time.Minute {
				mgr.startDelay = 0
			}
			mgr.startDelay += 3000 * time.Millisecond
			if mgr.startDelay > 30000*time.Millisecond {
				mgr.startDelay = 5000 * time.Millisecond
			}
			//fmt.Printf("[zone %d] Restarting in %s\n", mgr.port, mgr.startDelay)
			mgr.lastError = ""
			mgr.setState(reporter.AppStateRestarting)
			mgr.errorCooldown = time.Now().Add(30 * time.Minute)
			mgr.errorCount = 0
			time.Sleep(mgr.startDelay)
			return
		}
	}
}

func (mgr *manager) lineParse(line string) {
	if strings.Contains(line, "[Error]") {
		mgr.lastError = line
		mgr.lastErrorAt = time.Now()
		mgr.errorCount++
		if mgr.errorCount >= 10 || mgr.state == reporter.AppStateStarting {
			mgr.errorCooldown = time.Now().Add(30 * time.Minute)
			mgr.setState(reporter.AppStateErroring)
		}
		return
	}
	if mgr.exeName == "zone" && strings.Contains(line, "Entering sleep mode") {
		mgr.setState(reporter.AppStateSleeping)
		return
	}
	if mgr.exeName == "world" && strings.Contains(line, "UDP Listening on") {
		mgr.setState(reporter.AppStateRunning)
		return
	}
	if mgr.exeName == "ucs" && strings.Contains(line, "Connected to World") {
		mgr.setState(reporter.AppStateRunning)
		return
	}
}
