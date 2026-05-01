package lib

import (
	"os/exec"
)

// Opens the folder in the filemanager
func OpenFolder(path string) error {
	return exec.Command("explorer", path).Run()
}