//go:build !windows

package lib

import (
	"errors"
	"fmt"
	"os/exec"
	"runtime"
)

// Opens the folder in the filemanager
func OpenFolder(path string) error {
	var cmd *exec.Cmd

	switch runtime.GOOS {
	case "darwin":
		cmd = exec.Command("open", path)
	case "linux":
		cmd = exec.Command("xdg-open", path)
	default:
		// Fallback for other Unix-like systems
		cmd = exec.Command("xdg-open", path)
	}

	return cmd.Run()
}

func OpenProject(path, target string) error {
	var err error

	switch target {
	case "filemanager", "f":
		err = fmt.Errorf("finder: %w", exec.Command("open", path).Run())
	case "terminal", "t":
		err = fmt.Errorf("terminal: %w", exec.Command("open", "-a", "Terminal", path).Run())
	case "zed", "z":
		err = fmt.Errorf("zed: %w", exec.Command("zed", path).Run())
	default:
		return fmt.Errorf("unknown target: %s", target)
	}

	if errors.Unwrap(err) != nil {
		return fmt.Errorf("opening path: %s with %s: %w", path, target, err)
	}

	return nil
}
