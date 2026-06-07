package cmd

import (
	"fmt"
	"log"
	"math"
	"slices"
	"time"

	"os"
	"strconv"
	"strings"
	"unicode/utf8"

	"github.com/Revi-Studios/project/lib"
	"github.com/spf13/cobra"
	"golang.org/x/term"
)

var listCmd = &cobra.Command{
	Use:   "list [flags]",
	Short: "List projects",
	Long:  "Filter different types of projects and list them",
	Args:  cobra.ExactArgs(0),
	RunE: func(cmd *cobra.Command, args []string) error {
		os.Chdir(lib.Config.ProjectPath() + "/")

		folders, err := func() ([]os.DirEntry, error) {
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
			return folders, nil
		}()

		if err != nil {
			return fmt.Errorf("getting folders: %w", err)
		}

		filters, _ := cmd.Flags().GetStringSlice("filter")
		excludes, _ := cmd.Flags().GetStringSlice("exclude")
		browses, _ := cmd.Flags().GetStringSlice("browse")

		filterFolders(&folders, filters)
		excludeFolders(&folders, excludes)
		browseFolders(&folders, browses)

		switch true {
		case len(folders) == 0 && (len(filters) != 0 || len(excludes) != 0):
			fmt.Println("No matching projects found.")
			return nil
		case len(folders) == 0:
			fmt.Println("No projects found.")
			return nil
		}

		var longestFileName int
		var longestTagName int

		str, _ := cmd.Flags().GetString("view")
		if str == "" {
			str = lib.Config.Defaults.View
		}

		switch str {

		case "grid", "g":

			// Format data
			projectsMap := make(map[string][]string, len(folders))
			projectData := make(map[string][2]int, len(folders))
			for _, folder := range folders {
				longestWordLength := 0
				name := " - " + folder.Name() + " "

				tmp, _ := lib.GetTags(folder.Name())
				tags := strings.Join(tmp, ", ")
				tags = "[" + tags + "] "

				// Adjusts the longest name if needed
				if len := utf8.RuneCountInString(name); len > longestWordLength {
					longestWordLength = len
				}
				if len := utf8.RuneCountInString(tags); len > longestWordLength {
					longestWordLength = len
				}

				if tags == "[] " {
					tags = "[Without Tags]"
				}
				if projectsMap[tags] == nil {
					projectsMap[tags] = append(projectsMap[tags], tags)
				}
				projectsMap[tags] = append(projectsMap[tags], name)

				// Project Data
				data := projectData[tags]

				data[1]++
				if data[0] < longestWordLength {
					data[0] = longestWordLength
				}

				projectData[tags] = data
			}

			// Calculate default box size
			var height, width int
			for _, n := range projectData {
				width += n[0]
				height += n[1] + 1
			}
			width /= len(projectData)
			height /= len(projectData)

			// Calculating category sizes relative to the defaul box sizes
			for c, n := range projectData {
				projectData[c] = [2]int{int(math.Ceil(float64(n[0]) / float64(width))), int(math.Ceil(float64(n[1]+2) / float64(height)))}
			}

			// Gets the terminal width
			termWidth, _, err := term.GetSize(int(os.Stdout.Fd()))
			if err != nil {
				return fmt.Errorf("reading terminal width: %w", err)
			}

			// Create needed variables
			var index int = 0
			var viewW int = termWidth / width
			var bitMap uint64
			var renderBuf = make(map[int]string)

			//	Maps the blocks into the renderbuffer
		MapBlocks:
			for {
				if len(projectData) == 0 {
					break MapBlocks
				}

				// Find emty space
				if getBit(&bitMap, index) == 0 {
					space := nextSetBitOrNewRow(&bitMap, viewW, index)

					round := 0
				FindBlock:
					for {
						for c, n := range projectData {
							// Finds a block that fits
							if n[0] == space-round {
								delete(projectData, c)
								insertBlockIntoBitMap(&bitMap, n, viewW, index)
								// projectsMap[c] = append([]string{projectsMap[c][0] /* + " [" + strconv.Itoa(n[0]) + ":" + strconv.Itoa(n[1]) + "] " */ + strings.Repeat("-", n[0]*width-len(projectsMap[c][0])-1) + "+"}, projectsMap[c][1:len(projectsMap[c])]...)

								for i, s := range projectsMap[c] {
									for column := range (len(s) + width) / width {
										str := s

										str = string([]rune(str)[min(column*width, len(str)):min(column*width+width, len(str))])

										if str == "" {
											continue
										}

										str = str + strings.Repeat(" ", width-len(str))

										if in := ((index+1)/viewW)*viewW*(height-1) + index + i*viewW + column; renderBuf[in] != "" {
											log.Fatal("\"", str, "\"", " cell: ", in, " already taken in buffer by: \"", renderBuf[in], "\"")
										}

										renderBuf[((index+1)/viewW)*viewW*(height-1)+index+i*viewW+column] = str

										if lib.Config.Debug && lib.Config.Debugger.SkipGrid == false {
											renderGrid(bitMap, viewW, width, height, renderBuf)
										}
									}
								}

								index += n[0] - 1
								break FindBlock
							}
						}
						round++
					}

				}

				index++
			}
			if lib.Config.Debug && lib.Config.Debugger.SkipGrid == false {
				return nil
			}

			var renderBuilder = strings.Builder{}

			// Writing every string piece to the renderBuilder for string building
			renderBuilder.Grow(len(strconv.FormatUint(bitMap, 2)) * height)
			for i := range len(strconv.FormatUint(bitMap, 2)) * height {
				var str string

				switch true {
				case renderBuf[i] != "":
					str = renderBuf[i]
					if len(browses) != 0 {
						for _, browse := range browses {
							str = strings.ReplaceAll(str, browse, "\033[7m"+browse+"\033[0m")
						}
					}
				default:
					str = strings.Repeat(" ", width)
				}

				renderBuilder.WriteString(str)

				if (i+1)%viewW == 0 {
					renderBuilder.WriteString("\n")
				}
			}

			fmt.Println(strings.Repeat("-", viewW*width/2-5), "Projects:", strings.Repeat("-", viewW*width/2-6))
			fmt.Println(renderBuilder.String())

		case "category", "c":
			projectsMap := make(map[string][]string)
			for _, folder := range folders {
				name := folder.Name()
				tmp, _ := lib.GetTags(name)
				tags := strings.Join(tmp, ", ")

				// Adjusts the longest names if needed
				if len := utf8.RuneCountInString(name); len > longestFileName {
					longestFileName = len
				}
				if len := utf8.RuneCountInString(tags); len > longestTagName {
					longestTagName = len
				}

				if tags == "" {
					tags = "[Without Tags]"
				}
				projectsMap[tags] = append(projectsMap[tags], name)
			}
			fmt.Println(strings.Repeat("-", 3), "Projects:", strings.Repeat("-", longestTagName+longestFileName-8))
			for category, prjs := range projectsMap {
				fmt.Println(category)
				for _, prjname := range prjs {
					fmt.Println(" -", func() string {
						if len(browses) != 0 {
							str := prjname
							for _, browse := range browses {
								str = strings.ReplaceAll(str, browse, "\033[7m"+browse+"\033[0m")
							}
							return str
						}
						return prjname
					}())
				}
				fmt.Println("")
			}
			fmt.Println(strings.Repeat("-", longestFileName+longestTagName+7))

		default:
			projects := make([][2]string, len(folders))
			for i, folder := range folders {
				name := folder.Name()
				tmp, _ := lib.GetTags(name)
				tags := strings.Join(tmp, ", ")

				if len := utf8.RuneCountInString(name); len > longestFileName {
					longestFileName = len
				}
				if len := utf8.RuneCountInString(tags); len > longestTagName {
					longestTagName = len
				}
				projects[i] = [2]string{name, tags}
			}
			fmt.Println(strings.Repeat("-", safeRepeat(longestFileName-2)), "Projects:", strings.Repeat("-", safeRepeat(longestTagName-2)))
			for _, prj := range projects {
				if prj == [2]string{"", ""} {
					continue
				}
				strLeft := strings.Repeat(" ", longestFileName-utf8.RuneCountInString(prj[0]))
				fmt.Println(func() string {
					if len(browses) != 0 {
						str := prj[0]
						for _, browse := range browses {
							str = strings.ReplaceAll(str, browse, "\033[7m"+browse+"\033[0m")
						}
						return str
					}
					return prj[0]
				}(), strLeft, "|", prj[1])
			}
			fmt.Println(strings.Repeat("-", longestFileName+longestTagName+7))

		}

		return nil
	},
}

func init() {
	listCmd.Flags().StringP("view", "v", "", "--view <view-type>")
	listCmd.Flags().StringSliceP("filter", "f", []string{}, "--filter <tags>")
	listCmd.Flags().StringSliceP("exclude", "e", []string{}, "--exclude <tags>")
	listCmd.Flags().StringSliceP("browse", "b", []string{}, "--browse <string>")
	listCmd.RegisterFlagCompletionFunc("view", func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return []string{"list", "category", "grid"}, cobra.ShellCompDirectiveNoFileComp
	})
	listCmd.RegisterFlagCompletionFunc("filter", func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		tags := make([]string, 0, len(lib.Config.Tags))
		for _, tag := range lib.Config.Tags {
			tags = append(tags, "\""+tag+"\"")
		}
		return tags, cobra.ShellCompDirectiveNoFileComp
	})
	listCmd.RegisterFlagCompletionFunc("exlude", func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		tags := make([]string, 0, len(lib.Config.Tags))
		for _, tag := range lib.Config.Tags {
			tags = append(tags, "\""+tag+"\"")
		}
		return tags, cobra.ShellCompDirectiveNoFileComp
	})
	rootCmd.AddCommand(listCmd)
}

// Functions
func getBit(bitMap *uint64, index int) int {
	return int((*bitMap >> index) & 1)
}

func nextSetBitOrNewRow(bitMap *uint64, rowLength, index int) int {
	r := 0
	for r = range rowLength - (index % rowLength) {
		if index++; getBit(bitMap, index) == 1 {
			return r + 1
		}
	}
	return r + 1
}

func insertBlockIntoBitMap(bitMap *uint64, block [2]int, rowLength, index int) {
	for line := range block[1] {
		for column := range block[0] {
			*bitMap |= (1 << uint(index+column+(line*rowLength)))
		}
	}
}

func renderGrid(bitMap uint64, viewW, width, height int, renderBuf map[int]string) {

	var renderBuilder = strings.Builder{}
	renderBuilder.Grow(viewW * height * viewW)

	// Writing every string piece to the renderBuilder for string building
	for i := range len(strconv.FormatUint(bitMap, 2)) * height {
		var str string

		switch true {
		case renderBuf[i] != "":
			str = renderBuf[i]
		default:
			str = strings.Repeat(" ", width)
		}

		renderBuilder.WriteString(str)

		if (i+1)%viewW == 0 {
			renderBuilder.WriteString("\n")
		}
	}
	fmt.Printf("\033[H")
	// Print the renderd string
	fmt.Println(renderBuilder.String())

	fmt.Println("Info:", width, height, viewW, strings.Repeat(" ", viewW*2))
	fmt.Println("[bitMap]", strings.Repeat(" ", viewW*2))
	for i, char := range reverse(strconv.FormatUint(bitMap, 2)) {
		fmt.Printf("|%c", char)
		// Check if we've reached the end of a row
		if (i+1)%viewW == 0 {
			fmt.Print("\n")
		}
	}

	time.Sleep(500 * time.Millisecond)
}

func reverse(s string) string {
	r := []rune(s)
	for i, j := 0, len(r)-1; i < j; i, j = i+1, j-1 {
		r[i], r[j] = r[j], r[i]
	}
	return string(r)
}

func safeRepeat(n int) int {
	if n < 0 {
		return 0
	}
	return n
}

func filterFolders(folders *[]os.DirEntry, filters []string) {
	if len(filters) == 0 {
		return
	}

	for i, f := range filters {
		filters[i] = lib.Config.Shorthands.Tag(f)
	}

	var filtered []os.DirEntry
	for _, folder := range *folders {
		tags, _ := lib.GetTags(folder.Name())
		if len(tags) == 0 {
			tags = []string{"Without Tags"}
		}

		for _, tag := range tags {
			if slices.Contains(filters, tag) {
				filtered = append(filtered, folder)
			}
		}
	}
	*folders = filtered
}

func excludeFolders(folders *[]os.DirEntry, excludes []string) {
	if len(excludes) == 0 {
		return
	}

	for i, e := range excludes {
		excludes[i] = lib.Config.Shorthands.Tag(e)
	}

	var filtered []os.DirEntry
FolderFilter:
	for _, folder := range *folders {
		tags, _ := lib.GetTags(folder.Name())
		if len(tags) == 0 {
			tags = []string{"Without Tags"}
		}

		for _, tag := range tags {
			if slices.Contains(excludes, tag) {
				continue FolderFilter
			}
		}
		filtered = append(filtered, folder)
	}
	*folders = filtered
}

func browseFolders(folders *[]os.DirEntry, browses []string) {
	if len(browses) == 0 {
		return
	}

	var filtered []os.DirEntry

	for _, folder := range *folders {
		for _, str := range browses {
			if strings.Contains(folder.Name(), str) {
				filtered = append(filtered, folder)
				break
			}
		}
	}

	if len(browses) == 1 {
		fmt.Println("Search string:", strings.Join(browses, ", "))
	} else {
		fmt.Println("Search strings:", strings.Join(browses, ", "))
	}

	*folders = filtered
}
