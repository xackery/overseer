package main

import (
	"context"
	"fmt"
	"os"
	"runtime"
	"syscall"
	"time"

	"github.com/erikgeiser/promptkit/confirmation"
	"github.com/shirou/gopsutil/v3/process"
	"github.com/xackery/overseer/share/config"
	"github.com/xackery/overseer/share/message"
	"github.com/xackery/overseer/share/operation"
)

var (
	Version = "0.0.0"
	winExt  = ""
)

// icon link: https://prefinem.com/simple-icon-generator/#eyJiYWNrZ3JvdW5kQ29sb3IiOiIjMDA4MGZmIiwiYm9yZGVyQ29sb3IiOiIjMDAwMDAwIiwiYm9yZGVyV2lkdGgiOiI0IiwiZXhwb3J0U2l6ZSI6IjI1NiIsImV4cG9ydGluZyI6dHJ1ZSwiZm9udEZhbWlseSI6IkFiaGF5YSBMaWJyZSIsImZvbnRQb3NpdGlvbiI6IjU1IiwiZm9udFNpemUiOiIyMyIsImZvbnRXZWlnaHQiOjYwMCwiaW1hZ2UiOiIiLCJpbWFnZU1hc2siOiIiLCJpbWFnZVNpemUiOjUwLCJzaGFwZSI6InNxdWFyZSIsInRleHQiOiJJTlNUQUxMIn0
func main() {
	err := run()
	if err != nil {
		message.Badf("Stop failed: %s\n", err)
		operation.Exit(1)
	}
	message.OK("Stop complete\n")
	operation.Exit(0)
}

func run() error {
	config, err := config.LoadOverseerConfig("overseer.ini")
	if err != nil {
		return fmt.Errorf("load overseer config: %w", err)
	}
	message.Banner("Stop v" + Version)
	cwd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("getwd: %w", err)
	}

	fmt.Printf("Working directory: %s\n", cwd)

	choice, err := confirmation.New("Stop overseer and all processes?", confirmation.Yes).RunPrompt()
	if err != nil {
		return fmt.Errorf("select stop: %w", err)
	}
	if !choice {
		message.OK("OK, exiting\n")
		return nil
	}

	if runtime.GOOS == "windows" {
		winExt = ".exe"
	}

	processes := []string{
		"overseer",
		"world",
		"zone",
	}

	processes = append(processes, config.Apps...)

	totalCount := 0

	delay := 5 * time.Second
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	count, err := KillAllProcessesCtx(ctx, "overseer", 1)
	if err != nil {
		return fmt.Errorf("kill overseer: %w", err)
	}
	if count > 0 {
		message.OK("Sent overseer stop signal")
		totalCount += count
		message.Emote("⏱️", " Waiting 5 seconds for overseer to stop")
		time.Sleep(delay)
	} else {
		message.OK("Overseer isn't running")
	}

	message.Emote("💻", "Stopping everything gracefully")
	for _, process := range processes {
		ctxTimeout := 500 * time.Millisecond
		ctx, cancel := context.WithTimeout(context.Background(), ctxTimeout)
		defer cancel()

		count, err := KillAllProcessesCtx(ctx, process, 1)
		if err != nil {
			return fmt.Errorf("kill %s: %w", process, err)
		}
		totalCount += count
		if count > 0 {
			message.OKf("Stopped %d %s processes\n", count, process)
		}
	}
	if totalCount > 0 {
		message.OKf("Stopped %d total processes\n", totalCount)
	}
	if totalCount > 0 {
		message.Emote("⏱️", " Waiting 5 seconds for programs to stop")
		time.Sleep(5 * time.Second)
	}

	if totalCount == 0 {
		message.OK("Nothing to stop")
		return nil
	}

	message.Emote("💻", "Killing everything not exiting")
	totalCount = 0
	for _, process := range processes {
		ctxTimeout := 500 * time.Millisecond
		ctx, cancel := context.WithTimeout(context.Background(), ctxTimeout)
		defer cancel()

		count, err := KillAllProcessesCtx(ctx, process, 9)
		if err != nil {
			return fmt.Errorf("kill %s: %w", process, err)
		}
		totalCount += count
		if count > 0 {
			message.OKf("Killed %d %s processes\n", count, process)
		}
	}

	if totalCount > 0 {
		message.OKf("Killed %d total processes\n", totalCount)
	}
	return nil
}

func KillAllProcessesCtx(ctx context.Context, name string, signal int) (count int, err error) {
	processes, err := process.ProcessesWithContext(ctx)
	if err != nil {
		err = fmt.Errorf("processes: %w", err)
		return
	}
	var n string
	for _, p := range processes {
		n, err = p.NameWithContext(ctx)
		if err != nil {
			err = fmt.Errorf("name: %w", err)
			return
		}
		if n != name && n != name+winExt {
			continue
		}
		count++
		if signal == 1 {
			err = p.SendSignal(syscall.SIGINT)
			if err != nil {
				err = fmt.Errorf("send signal: %w", err)
				return
			}
			continue
		}

		err = p.KillWithContext(ctx)
		if err != nil {
			err = fmt.Errorf("kill: %w", err)
			return
		}
	}
	return
}
