package lib

import "os/exec"

func OpenFolder(path string) error {
	err := exec.Command("open", path).Run()

	return err
}
