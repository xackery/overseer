package runner

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/xackery/overseer/pkg/flog"
)

// Runner handles running and polling output of a process
type ProcessRunner struct {
	outChan  chan (string)
	doneChan chan (error)
	path     string
	name     string
	args     []string
	cmd      *exec.Cmd
}

func NewProcess(outChan chan (string), doneChan chan (error), path string, name string, args ...string) *ProcessRunner {
	return &ProcessRunner{
		outChan:  outChan,
		doneChan: doneChan,
		path:     path,
		name:     name,
		args:     args,
	}
}

func (r *ProcessRunner) Start(ctx context.Context) {
	if r.cmd != nil {
		flog.Printf("process %s already running\n", r.name)
		return
	}

	absPath, err := filepath.Abs(r.path)
	if err != nil {
		absPath = r.path
	}

	flog.Printf("Runner exec Starting process '%s %s' from path %s\n", r.name, strings.Join(r.args, " "), absPath)
	r.cmd = exec.CommandContext(ctx, r.name, r.args...)
	r.cmd.Dir = r.path
	err = r.run()
	r.cmd = nil
	flog.Printf("Runner %s finished with error: %s\n", r.name, err)
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

	flog.Printf("Runner start process %s\n", r.name)
	err = r.cmd.Start()
	if err != nil {
		return fmt.Errorf("start: %w", err)
	}

	flog.Printf("Runner wait process %s\n", r.name)
	err = r.cmd.Wait()
	if err != nil {
		return fmt.Errorf("wait: %w", err)
	}

	flog.Printf("Runner process %s self exit\n", r.name)
	return nil
}

func (r *ProcessRunner) Stop() error {
	if r.cmd == nil {
		return nil
	}
	flog.Printf("Stopping process %s\n", r.name)
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
