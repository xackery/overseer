package runner

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"os/exec"

	"github.com/xackery/overseer/pkg/flog"
)

// Runner handles running and polling output of a process
type ProcessRunner struct {
	outChan     chan (string)
	doneChan    chan (error)
	displayName string
	wdPath      string
	exePath     string
	name        string
	args        []string
	cmd         *exec.Cmd
}

func NewProcess(outChan chan (string), doneChan chan (error), displayName string, wdPath string, exePath string, name string, args ...string) *ProcessRunner {
	return &ProcessRunner{
		outChan:     outChan,
		doneChan:    doneChan,
		displayName: displayName,
		wdPath:      wdPath,
		exePath:     exePath,
		name:        name,
		args:        args,
	}
}

// Start starts the process
func (r *ProcessRunner) Start(ctx context.Context) {
	if r.cmd != nil {
		flog.Printf("[runner][%s] already running\n", r.displayName)
		return
	}

	fullCmd := r.name
	for _, arg := range r.args {
		fullCmd += " " + arg
	}

	flog.Printf("[runner][%s] priming wdPath: '%s', exePath: '%s', exeCommand: '%s'\n", r.displayName, r.wdPath, r.exePath, fullCmd)
	r.cmd = exec.CommandContext(ctx, r.exePath+"/"+r.name, r.args...)

	r.cmd.Dir = r.wdPath
	err := r.run()
	if err != nil {
		if err.Error() != "wait: signal: killed" {
			flog.Printf("[runner][%s] finished with error: %s\n", r.displayName, err)
		}
	}
	r.cmd = nil
	flog.Printf("[runner][%s] done\n", r.displayName)
	r.doneChan <- err
}

func (r *ProcessRunner) run() error {
	var err error

	stdout, err := r.cmd.StdoutPipe()
	if err != nil {
		return fmt.Errorf("stdout pipe: %w", err)
	}

	go func() {
		scanner := bufio.NewScanner(stdout)
		scanner.Split(bufio.ScanLines)
		for scanner.Scan() {
			r.outChan <- scanner.Text()
		}
	}()

	// don't pop up window for new process
	r.cmd.SysProcAttr = newProcAttr()

	flog.Printf("[runner][%s] starting process\n", r.displayName)
	err = r.cmd.Start()
	if err != nil {
		return fmt.Errorf("start %+v: %w", r.cmd, err)
	}

	flog.Printf("[runner][%s] wait process\n", r.displayName)
	err = r.cmd.Wait()
	if err != nil {
		return fmt.Errorf("wait: %w", err)
	}

	flog.Printf("[runner][%s] self exit\n", r.displayName)
	return nil
}

func (r *ProcessRunner) Stop() error {
	if r.cmd == nil {
		return nil
	}
	flog.Printf("[runner][%s] stopping\n", r.displayName)
	err := r.cmd.Process.Signal(os.Interrupt)
	if err != nil {
		return fmt.Errorf("signal: %w", err)
	}

	return nil
}

func (r *ProcessRunner) PID() int {
	if r.cmd == nil {
		return 0
	}
	if r.cmd.Process == nil {
		return 0
	}
	return r.cmd.Process.Pid
}
