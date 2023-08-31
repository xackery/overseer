//go:build mage

package main

import (
	"fmt"
	"os"
	"runtime"
	"strings"

	"github.com/magefile/mage/sh"
)

var (
	winExt = ""
)

func init() {
	if runtime.GOOS == "windows" {
		winExt = ".exe"
	}
}

// Runs build diagnose
func Build(target string) error {
	err := os.Chdir(target)
	if err != nil {
		return fmt.Errorf("Error changing directory: %v", err)
	}

	cmd := "build -o ../bin/" + target + winExt + " ."
	fmt.Println("Running: go " + cmd)
	err = sh.Run("go", strings.Split(cmd, " ")...)
	if err != nil {
		return fmt.Errorf("Failed build: %v", err)
	}
	return nil
}
