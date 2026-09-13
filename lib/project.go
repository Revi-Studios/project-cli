package lib

import (
	"errors"
	"fmt"
	"os/exec"
)

func OpenProjectInApp(name, target string) error {
	path := Config.ProjectPath() + "/" + name + "/"

	if target == "" {
		target = Config.Defaults.Open
	}

	var err error

	switch target {
	case "finder", "f":
		err = exec.Command("open", path).Run()
		target = "Finder"
	case "terminal", "t":
		err = exec.Command("open", "-a", "Terminal", path).Run()
		target = "Terminal"
	case "editor", "e":
		if Config.Editor == "" {
			err = errors.New("editor: not set in config")
			break
		}
		err = exec.Command(Config.Editor, path).Run()
		target = Config.Editor
	default:
		return fmt.Errorf("unknown target: %s", target)
	}

	return err
}
