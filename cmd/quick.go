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
	"golang.design/x/clipboard"
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

		fmt.Print("\033[?1049h\033[?25l\033[H")
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
				case 'A': // 65
					if quick.selection <= 0 {
						quick.selection = len(*quick.f_folders)
					}
					quick.selection--

					quick.RenderUI()
				case 'B': // 66
					quick.selection++
					if quick.selection >= len(*quick.f_folders) {
						quick.selection = 0
					}
					quick.RenderUI()
				case 'C': // 67
					quick.selection += min(quick.last_list_length, len(*quick.f_folders)-1-quick.selection)
					quick.RenderUI()
				case 'D': // 68
					quick.selection -= min(quick.last_list_length, quick.selection)
					quick.RenderUI()
				}
			} else if numRead == 1 {
				// Handle Single Keypresses
				switch buf[0] {
				case 13:
					term.Restore(int(os.Stdin.Fd()), oldState)
					if len(*quick.f_folders)-1 >= quick.selection {
						err := clipboard.Init()
						if err != nil {
							panic(err)
						}
						clipboard.Write(clipboard.FmtText, []byte((*quick.f_folders)[quick.selection].Name()))
					}
					return nil
				case 'q', 'Q', 3:
					return nil
				case 't', 'T':
					quick.selection = (quick.selection / quick.last_list_length) * quick.last_list_length
					quick.RenderUI()
				case 'b', 'B':
					quick.selection = min((quick.selection/quick.last_list_length)*quick.last_list_length+quick.last_list_length-1, len(*quick.f_folders)-1)
					quick.RenderUI()
				case 'm', 'M':
					quick.selection = min((quick.selection/quick.last_list_length)*quick.last_list_length+(quick.last_list_length-1)/2, len(*quick.f_folders)-1)
					quick.RenderUI()
				case 's', 'S', ':':
					quick.InitSearch()
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

	selection      int // quick.selection
	hide_selection bool

	filters  []string
	excludes []string
	browses  []string

	render_filters bool
	render_help    bool

	last_height      int
	last_list_length int
	left_space       int // How many rows are left before the terminal size is used up

	search bool
}

// Filters the folders with filters, excludes and browses and puts the result in f_folder
func (quick *Quick) filter() {
	*quick.f_folders = *quick.folders

	filterFolders(quick.f_folders, quick.filters)
	excludeFolders(quick.f_folders, quick.excludes)
	browseFolders(quick.f_folders, quick.browses)
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

func (quick *Quick) renderHelp(builder *strings.Builder) {
	builder.WriteString("\r")

	builder.WriteString("* Help Message:\n\r")
	builder.WriteString("* :h -> Show this help message.\n\r")
	builder.WriteString("* :a -> Show all the active filtering.\n\r")
	builder.WriteString("*\n\r")
	builder.WriteString("* :f -> Adds filters.\n\r")
	builder.WriteString("* :e -> Adds excludes.\n\r")
	builder.WriteString("* :b -> Adds browses.\n\r")
	builder.WriteString("*\n\r")
	builder.WriteString("* :f:c -> Clears filters.\n\r")
	builder.WriteString("* :e:c -> Clears excludes.\n\r")
	builder.WriteString("* :b:c -> Clears browses.\n\r")
	builder.WriteString("*\n\r")
	builder.WriteString("* :c -> Clears all filterings.\n\r")

	builder.WriteString("\n\r")
}

func (quick *Quick) renderFilters(builder *strings.Builder) {
	builder.WriteString("\r")

	builder.WriteString("* Filters: [\x1b[94m")
	builder.WriteString(strings.Join(quick.filters, "\x1b[0m, \x1b[94m"))
	builder.WriteString("\x1b[0m]\n\r")

	builder.WriteString("* Excludes: [\x1b[94m")
	builder.WriteString(strings.Join(quick.excludes, "\x1b[0m, \x1b[94m"))
	builder.WriteString("\x1b[0m]\n\r")

	builder.WriteString("* Browses: [")
	if len(quick.browses) > 0 {
		builder.WriteString("\x1b[92m\"")
	}
	builder.WriteString(strings.Join(quick.browses, "\x1b[0m\", \x1b[92m\""))
	if len(quick.browses) > 0 {
		builder.WriteString("\"\x1b[0m")
	}
	builder.WriteString("]\n\r")

	builder.WriteString("\n\r")
}

func (quick *Quick) renderSearchBar(builder *strings.Builder) {
	builder.WriteString("\r:\n\r")
}

func (quick *Quick) renderList(builder *strings.Builder) {
	os.Chdir(lib.Config.ProjectPath())

	space := quick.getFreeSpace(builder)

	if space < 1 {
		return
	}

	avgW := 0
	folders := 0

	quick.last_list_length = min(space, len(*quick.f_folders))

	builder.WriteString("\r")

	for i, folder := range (*quick.f_folders)[(quick.selection/space)*space : min((quick.selection/space)*space+space, len(*quick.f_folders))] {
		tags, _ := lib.GetTags(folder.Name())
		tags_str := strings.Join(tags, ", ")
		avgW += len(folder.Name()) + 3
		folders++
		if i+(quick.selection/space)*space == (quick.selection)%len(*quick.f_folders) && !quick.hide_selection {
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
	if min((quick.selection/space)*space+space, len(*quick.f_folders))-(quick.selection/space)*space == 0 {
		builder.WriteString("   no projects found.")
	}

	if !(len(*quick.f_folders) <= space) {
		avgW /= folders
		left_margin := avgW/2 + 5

		str := strconv.Itoa((quick.selection/space)+1) + "/" + strconv.Itoa((len(*quick.f_folders)-1)/space+1)

		switch true {
		case quick.selection/space == 0:
			builder.WriteString(strings.Repeat(" ", left_margin))
			builder.WriteString(str)
			builder.WriteString(" >")
		case (quick.selection / space) == (len(*quick.f_folders)-1)/space:
			builder.WriteString(strings.Repeat(" ", left_margin-2))
			builder.WriteString("< ")
			builder.WriteString(str)
		default:
			builder.WriteString(strings.Repeat(" ", left_margin))
			builder.WriteString(str)
		}
	}
}

func (quick *Quick) RenderUI() {
	var builder strings.Builder

	builder.WriteString("\x1b[2K")
	builder.WriteString(strings.Repeat("\x1b[1A\x1b[2K", quick.last_height))

	if quick.search {
		quick.renderSearchBar(&builder)
	}
	if quick.render_help {
		quick.renderHelp(&builder)
	}
	if quick.render_filters {
		quick.renderFilters(&builder)
	}
	quick.renderList(&builder)

	result := builder.String()

	if len(result) > 0 {
		quick.last_height = bytes.Count([]byte(result), []byte{'\n'}) + 1
	}

	fmt.Print(builder.String())
}

func (quick *Quick) InitSearch() {
	quick.search = true
	quick.hide_selection = true
	quick.RenderUI()

	fmt.Print("\033[?25h")
	defer fmt.Print("\033[?25l")

	fmt.Print("\x1b[H\x1b[K2\r:")

	buf := make([]byte, 1)
	var input []byte

	for {
		_, err := os.Stdin.Read(buf)
		if err != nil {
			break
		}

		char := buf[0]

		// Handle Enter and Esc press
		if char == 13 || char == 10 || char == 27 {
			if len(input) > 0 && !(char == 27) {
				prefix := string(input[:2])
				doublePrefix := string(input[:4])

				switch true {
				case doublePrefix == ":f:c" || doublePrefix == ":F:c":
					quick.filters = []string{}

				case prefix == ":f" || prefix == ":F":
					quick.filters = append(quick.filters, strings.Split(string(input[2:]), ",")...)

				case doublePrefix == ":e:c" || doublePrefix == ":E:c":
					quick.excludes = []string{}

				case prefix == ":e" || prefix == ":E":
					quick.excludes = append(quick.excludes, strings.Split(string(input[2:]), ",")...)

				case doublePrefix == ":b:c" || doublePrefix == ":B:c":
					quick.browses = []string{}

				case prefix == ":b" || prefix == ":B":
					strs := strings.Split(string(input[2:]), ",")

					for i := range strs {
						strs[i] = strings.TrimSpace(strs[i])
					}

					quick.browses = append(quick.browses, strs...)

				case prefix == ":c" || prefix == ":C":
					quick.filters = []string{}
					quick.excludes = []string{}
					quick.browses = []string{}

				case prefix == ":a" || prefix == ":A":
					quick.render_filters = !quick.render_filters

				case prefix == ":h" || prefix == ":H":
					quick.render_help = !quick.render_help

				default:
					quick.browses = []string{string(input)}
				}

				quick.filter()
			}

			quick.search = false

			quick.hide_selection = false
			quick.selection = 0
			fmt.Print("\x1b[" + strconv.Itoa(quick.last_height) + "B")
			quick.RenderUI()
			break
		}

		// 3. Handle Backspace (ASCII 127 or 8)
		if (char == 127 || char == 8) && len(input) > 0 {
			input = input[:len(input)-1]
			fmt.Print("\b \b")

			continue
		}

		// 4. Handle standard printable characters
		if char >= 32 && char <= 126 {
			input = append(input, char)
			fmt.Print(string(char)) // Manually echo the character to the terminal screen
		}
	}
}
