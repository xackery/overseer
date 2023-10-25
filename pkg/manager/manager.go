package manager

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"github.com/xackery/overseer/pkg/flog"
	"github.com/xackery/overseer/pkg/message"
	"github.com/xackery/overseer/pkg/reporter"
	"github.com/xackery/overseer/pkg/runner"
	"github.com/xackery/overseer/pkg/signal"
)

type manager struct {
	ctx           context.Context
	displayName   string
	path          string
	exeName       string
	args          []string
	startDelay    time.Duration
	lastStartTime time.Time
	state         reporter.AppState
	restartCount  int
	lastError     string
	lastErrorAt   time.Time // When lastErrorAt hits 30 minutes, reset errorCount
	errorCooldown time.Time
	errorCount    int // When errorCount hits 3, set errorCooldown to 30 minutes
	doneChan      chan error
	outChan       chan string
}

type SetupType int

const (
	SetupDefault SetupType = iota
	SetupDocker
)

func (e *manager) setState(state reporter.AppState) {
	e.state = state
	reporter.SetAppState(e.displayName, state)
}

func (e *manager) setPID(pid int) {
	reporter.SetAppPID(e.displayName, pid)
}

// Manage is the main loop for the zone.
func Manage(setup SetupType, displayName string, path string, exeName string, args ...string) {
	go poll(displayName, path, exeName, args...)
}

func InitializeDockerNetwork(networkName string) error {
	if networkName == "" {
		return fmt.Errorf("network is empty")
	}
	if networkName == "bridge" || networkName == "host" || networkName == "none" {
		return fmt.Errorf("network is set to %s (invalid)", networkName)
	}

	ctx := context.Background()

	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		return fmt.Errorf("new env client: %w", err)
	}

	networks, err := cli.NetworkList(ctx, types.NetworkListOptions{})
	if err != nil {
		return fmt.Errorf("network list: %w", err)
	}

	for _, network := range networks {
		if network.Name == networkName {
			return nil
		}
	}

	_, err = cli.NetworkCreate(ctx, networkName, types.NetworkCreate{})
	if err != nil {
		return fmt.Errorf("network create: %w", err)
	}

	message.OKf("Created docker network %s\n", networkName)

	return nil
}

func poll(displayName string, path string, exeName string, args ...string) {
	signal.AddWorker()
	defer signal.FinishWorker()

	mgr := &manager{
		ctx:         signal.Ctx(),
		displayName: displayName,
		path:        path,
		exeName:     exeName,
		args:        args,
		outChan:     make(chan string),
		lastError:   "none",
		doneChan:    make(chan error),
	}

	run := runner.NewProcess(mgr.outChan, mgr.doneChan, mgr.path, mgr.exeName, mgr.args...)
	for {
		select {
		case <-mgr.ctx.Done():
			flog.Printf("[%s] Manager exiting ctx done\n", mgr.displayName)
			mgr.setState(reporter.AppStateStopped)
			run.Stop()
			return
		default:
		}
		mgr.lastStartTime = time.Now()
		go run.Start(mgr.ctx)
		mgr.setState(reporter.AppStateStarting)
		mgr.setPID(run.PID())

		parse(mgr, run)
	}
}

func parse(mgr *manager, run *runner.ProcessRunner) {
	start := time.Now()
	for {
		mgr.setPID(run.PID())
		select {
		case line := <-mgr.outChan:
			mgr.lineParse(line)

			flog.Printf("[%s] Manager line: %s\n", mgr.displayName, line)
		case <-mgr.ctx.Done():
			flog.Printf("[%s] Manager exiting ctx done\n", mgr.displayName)
			return
		case <-mgr.doneChan:
			mgr.restartCount++

			flog.Printf("[%s] Manager exited after %s seconds, %d restarts. Last error: %s\n", mgr.displayName, time.Since(start).Round(time.Second), mgr.restartCount, mgr.lastError)
			if time.Since(start) > 3*time.Minute {
				mgr.startDelay = 0
			}
			mgr.startDelay += 3000 * time.Millisecond
			if mgr.startDelay > 30000*time.Millisecond {
				mgr.startDelay = 5000 * time.Millisecond
			}
			flog.Printf("[%s] Restarting in %s\n", mgr.displayName, mgr.startDelay)
			mgr.lastError = ""
			mgr.setState(reporter.AppStateRestarting)
			mgr.errorCooldown = time.Now().Add(30 * time.Minute)
			mgr.errorCount = 0
			time.Sleep(mgr.startDelay)
			return
		case <-time.After(10 * time.Second):
			if time.Since(mgr.lastStartTime) > 10*time.Second && mgr.state == reporter.AppStateStarting {
				mgr.setState(reporter.AppStateRunning)
			}
		}
	}
}

func (mgr *manager) lineParse(line string) {
	if strings.Contains(line, "[Error]") {

		mgr.lastError = line
		mgr.lastErrorAt = time.Now()
		flog.Printf("[%s] Error: %s\n", mgr.displayName, line)
		mgr.errorCount++
		if mgr.errorCount >= 10 || mgr.state == reporter.AppStateStarting {
			mgr.errorCooldown = time.Now().Add(30 * time.Minute)
			mgr.setState(reporter.AppStateErroring)
		}
		return
	}
	if strings.Contains(mgr.exeName, "zone") && strings.Contains(line, "Entering sleep mode") {
		flog.Printf("[%s] Entering sleep mode\n", mgr.displayName)
		mgr.setState(reporter.AppStateSleeping)
		return
	}

	if strings.Contains(mgr.exeName, "zone") &&
		mgr.state == reporter.AppStateSleeping &&
		strings.Contains(line, "Zone booted successfully") {
		flog.Printf("[%s] Zone booted successfully\n", mgr.displayName)
		mgr.setState(reporter.AppStateRunning)
		return
	}

	if mgr.exeName == "world" && strings.Contains(line, "Starting EQ Network server on") {
		flog.Printf("[world] Started, got 'Starting EQ Network server on'\n")
		mgr.setState(reporter.AppStateRunning)
		return
	}
	if mgr.exeName == "ucs" && strings.Contains(line, "Connected to World") {
		flog.Printf("[ucs] Connected to World\n")
		mgr.setState(reporter.AppStateRunning)
		return
	}

	if time.Since(mgr.lastStartTime) > 10*time.Second && mgr.state == reporter.AppStateStarting {
		mgr.setState(reporter.AppStateRunning)
		return
	}

}
