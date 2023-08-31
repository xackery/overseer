package runner

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
)

// Runner handles running and polling output of a process
type ProcessRunner struct {
	outChan  chan (string)
	doneChan chan (error)
	name     string
	args     []string
	cmd      *exec.Cmd
}

func NewProcess(outChan chan (string), doneChan chan (error), name string, args ...string) *ProcessRunner {
	return &ProcessRunner{
		outChan:  outChan,
		doneChan: doneChan,
		name:     name,
		args:     args,
	}
}

func (r *ProcessRunner) Start() {
	if r.cmd != nil {
		return
	}
	r.cmd = exec.Command(r.name, r.args...)
	err := r.run()
	r.cmd = nil
	r.doneChan <- err
}

func (r *ProcessRunner) run() error {
	var err error

	cwd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("getwd: %w", err)
	}

	r.cmd.Dir = cwd

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

	err = r.cmd.Start()
	if err != nil {
		return fmt.Errorf("start: %w", err)
	}

	err = r.cmd.Wait()
	if err != nil {
		return fmt.Errorf("wait: %w", err)
	}

	return nil
}

func (r *ProcessRunner) Stop() error {
	if r.cmd == nil {
		return nil
	}
	err := r.cmd.Process.Signal(os.Interrupt)
	if err != nil {
		return fmt.Errorf("signal: %w", err)
	}

	return nil
}
