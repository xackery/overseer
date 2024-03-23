package main

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/erikgeiser/promptkit/confirmation"
	"github.com/erikgeiser/promptkit/selection"
	"github.com/shirou/gopsutil/v3/process"
	"github.com/xackery/overseer/pkg/config"
	"github.com/xackery/overseer/pkg/message"
	"github.com/xackery/overseer/pkg/operation"
	"github.com/xackery/overseer/pkg/sanity"
)

var (
	Version = "0.0.0"
)

func main() {
	err := run()
	if err != nil {
		message.Badf("Start failed: %s\n", err)
		operation.Exit(1)
	}
	operation.Exit(0)
}

func run() error {
	winExt := ""
	if runtime.GOOS == "windows" {
		winExt = ".exe"
	}

	cfg, err := config.LoadOverseerConfig("overseer.ini")
	if err != nil {
		return fmt.Errorf("load overseer config: %w", err)
	}
	message.Banner("Start v" + Version)
	cwd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("getwd: %w", err)
	}

	fmt.Printf("Working directory: %s\n", cwd)

	err = sanity.EssentialFiles(cfg)
	if err != nil {
		return fmt.Errorf("essential files: %w", err)
	}

	optionList := []string{
		"overseer (all)",
		"shared_memory",
		"world",
		"zone",
		"queryserv",
		"ucs",
		"loginserver",
	}

	command := ""
	if len(os.Args) > 1 {
		command = os.Args[1]
	} else {
		command, err = selection.New("Start which program?", optionList).RunPrompt()
		if err != nil {
			return fmt.Errorf("start: %w", err)
		}
	}

	fi, err := os.Stat(cwd + "/" + cfg.ServerPath)
	if err == nil {
		cfg.ServerPath = cwd + "/" + cfg.ServerPath
	} else {
		fi, err = os.Stat(cfg.ServerPath)
		if err != nil {
			return fmt.Errorf("stat: %w", err)
		}
	}
	if !fi.IsDir() {
		return fmt.Errorf("server path is not a directory")
	}

	fi, err = os.Stat(cwd + "/" + cfg.BinPath)
	if err == nil {
		cfg.BinPath = cwd + "/" + cfg.BinPath
	} else {
		fi, err = os.Stat(cfg.BinPath)
		if err != nil {
			return fmt.Errorf("stat: %w", err)
		}
	}
	if !fi.IsDir() {
		return fmt.Errorf("bin path is not a directory")
	}

	choice := command
	dir, err := filepath.Abs(cfg.ServerPath)
	if err != nil {
		return fmt.Errorf("abs: %w", err)
	}
	args := []string{}
	isScreen := false
	switch choice {
	case "overseer (all)":
		dir = cwd
		if cfg.IsScreenStart {
			screenPath, err := exec.LookPath("screen")
			if err != nil {
				message.Bad("screen is not installed. Skipping screen start.")
			}
			args = append(args, "-S", "overseer", "-t", "overseer", "-dm", "./overseer"+winExt)
			command = screenPath
			isScreen = true
		} else {
			command = "./overseer" + winExt
		}
	case "overseer":
		command = "./overseer" + winExt
		dir = cwd
	case "shared_memory":
		command, err = filepath.Rel(dir, cwd+"/"+cfg.BinPath+"/shared_memory"+winExt)
		if err != nil {
			return fmt.Errorf("rel: %w", err)
		}
	case "world":
		command, err = filepath.Rel(dir, cwd+"/"+cfg.BinPath+"/world"+winExt)
		if err != nil {
			return fmt.Errorf("rel: %w", err)
		}
	case "zone":
		command, err = filepath.Rel(dir, cwd+"/"+cfg.BinPath+"/zone"+winExt)
		if err != nil {
			return fmt.Errorf("rel: %w", err)
		}
	case "queryserv":
		command, err = filepath.Rel(dir, cwd+"/"+cfg.BinPath+"/queryserv"+winExt)
		if err != nil {
			return fmt.Errorf("rel: %w", err)
		}
	case "ucs":
		command, err = filepath.Rel(dir, cwd+"/"+cfg.BinPath+"/ucs"+winExt)
		if err != nil {
			return fmt.Errorf("rel: %w", err)
		}
	case "loginserver":
		command, err = filepath.Rel(dir, cwd+"/"+cfg.BinPath+"/loginserver"+winExt)
		if err != nil {
			return fmt.Errorf("rel: %w", err)
		}
	default:
		return fmt.Errorf("unknown argument %s", command)
	}

	if choice == "overseer (all)" {
		choice = "overseer"
	}
	/*if !service.IsDatabaseUp() {
		if cfg.PortableDatabase != 1 {
			message.Bad("Database is not running and we aren't portable. Please run it first.")
			return fmt.Errorf("database not running")
		}
		err = service.DatabaseStart()
		if err != nil {
			return fmt.Errorf("database start: %w", err)
		}

		if !service.IsDatabaseUp() {
			message.Bad("Database is not running even after trying to start it.")
			return fmt.Errorf("database failed to start")
		}
	}*/

	if !strings.Contains(choice, "zone") {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		processes, err := process.ProcessesWithContext(ctx)
		if err != nil {
			return fmt.Errorf("processes: %w", err)
		}
		var n string
		for _, p := range processes {
			n, err = p.NameWithContext(ctx)
			if err != nil {
				continue
			}
			if !strings.Contains(n, choice) {
				continue
			}
			isOK, err := confirmation.New(fmt.Sprintf("%s is already running. Start another copy?", n), confirmation.No).RunPrompt()
			if err != nil {
				return fmt.Errorf("confirmation: %w", err)
			}
			if !isOK {
				message.OK("OK, exiting")
				return nil
			}
			break
		}
	}

	fmt.Println("Running", command, strings.Join(args, " "), "from", dir)
	start := time.Now()
	cmd := exec.Command(command, args...)
	cmd.Dir = dir
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Env = os.Environ()

	if isScreen {
		err = cmd.Start()
		if err != nil {
			return fmt.Errorf("start %s %+v: %w", command, cmd, err)
		}
		message.OK("Screen of overseer started. You can use `screen -r overseer` to view it.")
		return nil
	} else {
		err = cmd.Run()
	}
	if err != nil {
		message.Badf("Start %s exited after %0.2f seconds\n", command, time.Since(start).Seconds())
		if exitError, ok := err.(*exec.ExitError); ok {
			return fmt.Errorf("start %s: %s", command, exitError.Error())
		}
		return fmt.Errorf("start %s: %w", command, err)
	}
	message.OKf("Start %s exited after %0.2f seconds\n", choice, time.Since(start).Seconds())

	return nil
}
