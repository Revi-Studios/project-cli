package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"path"
	"strings"

	"github.com/Revi-Studios/project/lib"
	"github.com/spf13/cobra"
	"golang.org/x/term"
)

func init() {
	rootCmd.AddCommand(quick)
}

var quick = &cobra.Command{
	Use:   "quick",
	Short: "Open the Quick opener",
	Long:  "Open the Quick opener that makes it quicker to open projects",
	RunE: func(cmd *cobra.Command, args []string) error {
		folders, err := func() ([]string, error) {
			files, err := os.ReadDir(lib.Config.ProjectPath() + "/")
			if err != nil {
				return nil, fmt.Errorf("reading the project directory: %w", err)
			}
			folders := make([]string, 0, len(files))
			for _, file := range files {
				if file.IsDir() {
					folders = append(folders, file.Name())
				}
			}
			return folders, nil
		}()

		if err != nil {
			return fmt.Errorf("getting folders: %w", err)
		}

		oldState, err := term.MakeRaw(int(os.Stdin.Fd()))

		if err != nil {
			return fmt.Errorf("setting raw mode: %e", err)
		}

		defer term.Restore(int(os.Stdin.Fd()), oldState)

		var selection int

		fmt.Print("\033[?1049h\033[?25l")
		defer fmt.Print("\033[?1049l\033[?25h")

		fmt.Println(renderListWithSelection(folders, selection))

		buf := make([]byte, 3)

		for {
			numRead, err := os.Stdin.Read(buf)
			if err != nil {
				break
			}

			if numRead == 3 && buf[0] == 27 && buf[1] == 91 {
				switch buf[2] {
				case 65:
					if selection <= 0 {
						selection = len(folders)
					}
					selection--
					fmt.Println(fmt.Sprintf("\033[%dA\r", len(folders)+1), renderListWithSelection(folders, selection))
				case 66:
					selection++
					fmt.Println(fmt.Sprintf("\033[%dA\r", len(folders)+1), renderListWithSelection(folders, selection))
				}
			} else if numRead == 1 {
				// Handle Single Keypresses
				switch buf[0] {
				case 13:
					term.Restore(int(os.Stdin.Fd()), oldState)
					exec.Command("open", path.Join(lib.Config.ProjectPath(), folders[selection%len(folders)])).Run()
					return nil
				case 'q', 'Q', 3:
					return nil
				}
			}
		}

		return nil
	},
}

func renderListWithSelection(list []string, selection int) string {
	var builder strings.Builder
	builder.WriteString("\r")
	for i, str := range list {
		if i == (selection)%len(list) {
			builder.WriteString(" > \033[7m")
			builder.WriteString(str)
			builder.WriteString("\033[0m")
			builder.WriteString("\n\r")
		} else {
			builder.WriteString("   ")
			builder.WriteString(str)
			builder.WriteString(strings.Repeat(" ", 10))
			builder.WriteString("\n\r")
		}

	}

	return builder.String()
}
