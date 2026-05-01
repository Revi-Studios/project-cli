//go:build !windows

package lib

import (
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
