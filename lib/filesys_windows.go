package lib

import (
	"errors"
	"fmt"
	"os/exec"
)

// Opens the folder in the filemanager
func OpenFolder(path string) error {
	return exec.Command("explorer", path).Run()
}

func OpenProject(path, target string) error {
	var err error

	switch target {
	case "filemanager", "f":
		err = fmt.Errorf("explorer: %w", exec.Command("open", path).Run())
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
