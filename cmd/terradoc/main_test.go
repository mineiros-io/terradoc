package main_test

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"testing"
)

var binName = "terradoc"

var terradocBinPath string

func TestMain(m *testing.M) {
	fmt.Println("Building tool...")

	if runtime.GOOS == "windows" {
		binName += ".exe"
	}

	binTmpdir, err := os.MkdirTemp("", "cmd-terradoc-test-")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Cannot create temp dir: %v", err)
		os.Exit(1)
	}

	terradocBinPath = filepath.Join(binTmpdir, binName)

	build := exec.Command("go", "build", "-o", terradocBinPath)
	err = build.Run()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Cannot build %q: %v", binName, err)
		os.Exit(1)
	}

	fmt.Println("Running tests...")

	result := m.Run()

	fmt.Println("Cleaning up...")
	err = os.RemoveAll(binTmpdir)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Cannot clean up temp dir %q: %v", binTmpdir, err)
	}

	os.Exit(result)
}
