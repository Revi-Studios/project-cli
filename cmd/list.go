package cmd

import (
	"fmt"
	"math"

	"os"
	"strconv"
	"strings"
	"unicode/utf8"

	"github.com/Revi-Studios/project/lib"
	"github.com/spf13/cobra"
	"golang.org/x/term"
)

func init() {
	rootCmd.AddCommand(listCmd)
}

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List all projects",
	Long:  "List all projects",
	Args:  cobra.ExactArgs(0),
	RunE: func(cmd *cobra.Command, args []string) error {
		config, err := lib.GetConfig()

		if err != nil {
			return fmt.Errorf("getting the config: %w", err)
		}

		files, err := os.ReadDir(config.ProjectFolderPath + "/")
		if err != nil {
			return fmt.Errorf("reading the project directory: %w", err)
		}

		os.Chdir(config.ProjectFolderPath + "/")

		var longestFileName int
		var longestTagName int

		str, _ := cmd.Flags().GetString("view")

		switch str {

		case "grid", "g":

			// Format data
			projectsMap := make(map[string][]string)
			projectData := make(map[string][2]int)
			for _, file := range files {
				if file.IsDir() {
					longestWordLength := 0
					name := " - " + file.Name()
					tags, _ := lib.GetTags(file.Name())

					// Adjusts the longest name if needed
					if len := utf8.RuneCountInString(name); len > longestWordLength {
						longestWordLength = len
					}
					if len := utf8.RuneCountInString(tags); len > longestWordLength {
						longestWordLength = len
					}

					if tags == "" {
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
				projectData[c] = [2]int{int(math.Ceil(float64(n[0]) / float64(width))), int(math.Ceil(float64(n[1]+3) / float64(height)))}
			}

			// Gets the terminal width
			termWidth, _, err := term.GetSize(int(os.Stdout.Fd()))

			if err != nil {
				return fmt.Errorf("reading terminal width: %w", err)
			}

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
					space := nextSetBitOrNewRow(&bitMap, viewW, index) + 1

					round := 0
				FindBlock:
					for {
						for c, n := range projectData {
							// Finds a block that fits
							if n[0] == space-round {
								delete(projectData, c)
								insertBlockIntoBitMap(&bitMap, n, viewW, index)

								for i, s := range projectsMap[c] {
									for column := range (len(s) + width) / width {
										str := s

										str = string([]rune(str)[min(column*width, len(str)):min(column*width+width, len(str))])
										str = str + strings.Repeat(" ", width-len(str))
										renderBuf[index*height+i*viewW+column] = str
									}
								}

								index += n[0] - 1
								break FindBlock
							}
							// Doesn't find a block that fits
							if space-round <= 0 {
								index += space
								break FindBlock
							}
						}
						round++
					}

				}

				index++
			}

			var renderBuilder = strings.Builder{}

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

			fmt.Println("Info:", width, height, viewW)
			fmt.Println("Projects:", len(projectsMap))

			// Print the renderd string
			fmt.Println(renderBuilder.String())

		case "category", "c":
			projectsMap := make(map[string][]string)
			for _, file := range files {
				if file.IsDir() {

					name := file.Name()
					tags, _ := lib.GetTags(name)

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
			}
			fmt.Println(strings.Repeat("-", 3), "Projects:", strings.Repeat("-", longestTagName+longestFileName-8))
			for category, prjs := range projectsMap {
				fmt.Println(category)
				for _, prjname := range prjs {
					fmt.Println(" -", prjname)
				}
				fmt.Println("")
			}
			fmt.Println(strings.Repeat("-", longestFileName+longestTagName+7))

		default:
			projects := make([][2]string, len(files))
			for i, file := range files {
				if file.IsDir() {
					name := file.Name()
					tags, err := lib.GetTags(name)
					if err != nil {
					}
					if len := utf8.RuneCountInString(name); len > longestFileName {
						longestFileName = len
					}
					if len := utf8.RuneCountInString(tags); len > longestTagName {
						longestTagName = len
					}
					projects[i] = [2]string{name, tags}

				}
			}
			fmt.Println(strings.Repeat("-", longestFileName-2), "Projects:", strings.Repeat("-", longestTagName-2))
			for _, prj := range projects {
				if prj == [2]string{"", ""} {
					continue
				}
				strLeft := strings.Repeat(" ", longestFileName-utf8.RuneCountInString(prj[0]))
				fmt.Println(prj[0], strLeft, "|", prj[1])
			}
			fmt.Println(strings.Repeat("-", longestFileName+longestTagName+7))

		}

		return nil
	},
}

func init() {
	listCmd.Flags().StringP("view", "v", "", "--view <view-type>")
	rootCmd.AddCommand(listCmd)
}

func getBit(bitMap *uint64, index int) int {
	return int((*bitMap >> index) & 1)
}

func nextSetBitOrNewRow(bitMap *uint64, rowLength, index int) int {
	r := 0
	for r = range rowLength - (index % rowLength) - 1 {
		if index++; getBit(bitMap, index) == 1 {
			return r
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
