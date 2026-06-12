package cmd

import (
	"bytes"
	"fmt"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"

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
		folders, err := func() (*[]os.DirEntry, error) {
			files, err := os.ReadDir(lib.Config.ProjectPath() + "/")
			if err != nil {
				return nil, fmt.Errorf("reading the project directory: %w", err)
			}
			folders := make([]os.DirEntry, 0, len(files))
			for _, file := range files {
				if file.IsDir() {
					folders = append(folders, file)
				}
			}
			return &folders, nil
		}()

		if err != nil {
			return fmt.Errorf("getting folders: %w", err)
		}

		quick := NewQuickMenu(folders)

		oldState, err := term.MakeRaw(int(os.Stdin.Fd()))

		if err != nil {
			return fmt.Errorf("setting raw mode: %e", err)
		}

		defer term.Restore(int(os.Stdin.Fd()), oldState)

		os.Chdir(lib.Config.ProjectPath()) // Change directory to the project folder so the filter functions read the correct folder data

		// quick.filters = []string{"go", "rust", "dart"}
		// quick.browses = []string{"l"}
		quick.filter()

		fmt.Print("\033[?1049h\033[?25l")
		defer fmt.Print("\033[?1049l\033[?25h")

		quick.RenderUI()

		buf := make([]byte, 3)
		sigChan := make(chan os.Signal, 1)

		signal.Notify(sigChan, syscall.SIGWINCH)

		go func() {
			for range sigChan {
				quick.RenderUI()
			}
		}()

		for {

			numRead, err := os.Stdin.Read(buf)
			if err != nil {
				break
			}

			if numRead == 3 && buf[0] == 27 && buf[1] == 91 {
				switch buf[2] {
				case 65:
					if quick.selection <= 0 {
						quick.selection = len(*quick.f_folders)
					}
					quick.selection--

					quick.RenderUI()
					// renderQuick(f_folders, quick.selection, filters, -1, &last_render_height)
					// fmt.Println(fmt.Sprintf("\033[%dA\r", len(*f_folders)+1), renderListWithSelection(f_folders, quick.selection))
				case 66:
					quick.selection++
					if quick.selection >= len(*quick.f_folders) {
						quick.selection = 0
					}
					quick.RenderUI() // renderQuick(f_folders, quick.selection, filters, -1, &last_render_height) // fmt.Println(fmt.Sprintf("\033[%dA\r", len(*f_folders)+1), renderListWithSelection(f_folders, quick.selection))
				}
			} else if numRead == 1 {
				// Handle Single Keypresses
				switch buf[0] {
				case 13:
					term.Restore(int(os.Stdin.Fd()), oldState)
					// exec.Command("open", path.Join(lib.Config.ProjectPath(), (*f_folders)[quick.selection].Name())).Run()
					return nil
				case 'q', 'Q', 3:
					return nil
				}
			}
		}

		return nil
	},
}

func NewQuickMenu(folders *[]os.DirEntry) *Quick {
	return &Quick{
		folders:   folders,
		f_folders: &[]os.DirEntry{},
		selection: 0,
	}
}

type Quick struct {
	folders   *[]os.DirEntry // All folders
	f_folders *[]os.DirEntry // Filtered folders

	selection int // quick.selection

	filters  []string
	excludes []string
	browses  []string

	last_height int
	left_space  int // How many rows are left before the terminal size is used up
}

// Filters the folders with filters, excludes and browses and puts the result in f_folder
func (quick *Quick) filter() {
	*quick.f_folders = *quick.folders

	browseFolders(quick.f_folders, quick.browses)
	filterFolders(quick.f_folders, quick.filters)
	excludeFolders(quick.f_folders, quick.excludes)
}

func (quick *Quick) getSize() (int, int, error) {
	width, heihgt, err := term.GetSize(int(os.Stderr.Fd()))

	if err != nil {
		return 0, 0, fmt.Errorf("reading terminal width: %w", err)
	}

	return width, heihgt, nil
}

func (quick *Quick) getFreeSpace(builder *strings.Builder) int {
	str := builder.String()
	_, height, err := quick.getSize()

	if err != nil {
		return 0
	}

	if len(str) > 0 {
		return height - (bytes.Count([]byte(str), []byte{'\n'}) + 1)
	}
	return height
}

func (quick *Quick) renderTextInput(builder *strings.Builder) {
	builder.WriteString("\rFilters: ")
	builder.WriteString("\x1b[3;94m" + strings.Join(quick.filters, "\x1b[0m, \x1b[3;94m") + "\x1b[0m\n\r")
	builder.WriteString("Excludes: ")
	builder.WriteString("\x1b[3;90m" + strings.Join(quick.excludes, "\x1b[0m, \x1b[3;2m") + "\x1b[0m\n\n\r")
	builder.WriteString("Search: \x1b[4;m l f:\x1b[94mgo\x1b[0m,\x1b[94mrust\x1b[0m e:\x1b[94mweb\x1b[0m\n\n\r")
	builder.WriteString(strings.Repeat("-", 20))
	builder.WriteString("\n")
}

func (quick *Quick) renderList(builder *strings.Builder) {
	os.Chdir(lib.Config.ProjectPath())

	space := quick.getFreeSpace(builder)

	avgW := 0
	folders := 0

	builder.WriteString("\r")

	for i, folder := range (*quick.f_folders)[(quick.selection/space)*space : min((quick.selection/space)*space+space, len(*quick.f_folders))] {
		tags, _ := lib.GetTags(folder.Name())
		tags_str := strings.Join(tags, ", ")
		avgW += len(folder.Name()) + 3
		folders++
		if i+(quick.selection/space)*space == (quick.selection)%len(*quick.f_folders) {
			builder.WriteString(" > \033[7m")
			builder.WriteString(folder.Name())
			builder.WriteString("\033[0m")
			if len(tags) > 0 {
				builder.WriteString(strings.Repeat(" ", 8))
				builder.WriteString("\033[90m[")
				builder.WriteString(tags_str)
				builder.WriteString("]\033[0m")
			}
			builder.WriteString("\n\r")
		} else {
			builder.WriteString("   ")
			builder.WriteString(folder.Name())
			builder.WriteString(strings.Repeat(" ", 8+2+len([]rune(tags_str))))
			builder.WriteString("\n\r")
		}
	}

	if !(len(*quick.f_folders) <= space) {
		avgW /= folders

		builder.WriteString(strings.Repeat(" ", avgW/2+5))
		// builder.WriteString("[" + strconv.Itoa((quick.selection / space)) + ":" + strconv.Itoa((quick.selection % len(*quick.f_folders))) + "]")
		builder.WriteString(strconv.Itoa((quick.selection/space)+1) + "/" + strconv.Itoa(((len(*quick.f_folders) / space) + 1)))
		if ((quick.selection/space)+1)*space <= len(*quick.f_folders) {
			builder.WriteString(" >")
		}
	}
}

func (quick *Quick) RenderUI() {
	var builder strings.Builder

	builder.WriteString("\x1b[2K")
	builder.WriteString(strings.Repeat("\x1b[1A\x1b[2K", quick.last_height))

	// quick.renderTextInput(&builder)
	quick.renderList(&builder)

	result := builder.String()

	if len(result) > 0 {
		quick.last_height = bytes.Count([]byte(result), []byte{'\n'}) + 1
	}

	fmt.Print(builder.String())
}

/*
	func readTextInput(oldState *term.State) string {
		// Temporarily show the blinking text cursor so the user knows where they are typing
		fmt.Print("\033[?25h")
		defer fmt.Print("\033[?25l") // Hide it again when done

		var inputBytes []byte
		buf := make([]byte, 1)

		for {
			_, err := os.Stdin.Read(buf)
			if err != nil {
				break
			}

			char := buf[0]

			// 1. Handle Enter (Finished Typing)
			if char == 13 || char == 10 {
				break
			}

			// 2. Handle Ctrl+C (Abort)
			if char == 3 {
				return ""
			}

			// 3. Handle Backspace (ASCII 127 or 8)
			if char == 127 || char == 8 {
				if len(inputBytes) > 0 {
					inputBytes = inputBytes[:len(inputBytes)-1]
					// Move cursor back, overwrite char with a space, move cursor back again
					fmt.Print("\b \b")
				}
				continue
			}

			// 4. Handle standard printable characters
			if char >= 32 && char <= 126 {
				inputBytes = append(inputBytes, char)
				fmt.Print(string(char)) // Manually echo the character to the terminal screen
			}
		}

		return string(inputBytes)
	}
*/
