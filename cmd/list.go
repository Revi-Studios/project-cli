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
			projectData := make(map[string][3]int)
			for _, file := range files {
				if file.IsDir() {
					longestWordLength := 0
					name := file.Name()
					tags, _ := lib.GetTags(name)

					// Adjusts the longest name if needed
					if len := utf8.RuneCountInString(name) + 3; len > longestWordLength {
						longestWordLength = len
					}
					if len := utf8.RuneCountInString(tags); len > longestWordLength {
						longestWordLength = len
					}

					if tags == "" {
						tags = "[Without Tags]"
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
				width += n[0] + 1
				height += n[1] + 2
			}
			width /= len(projectData)
			height /= len(projectData)

			// Calculating category sizes relative to the defaul box sizes
			for c, n := range projectData {
				projectData[c] = [3]int{int(math.Ceil(float64(n[0]) / float64(width))), int(math.Ceil(float64(n[1]+3) / float64(height)))}
			}

			/*
				// Default
				println("Default 1:1")
				fmt.Println("+" + strings.Repeat("-", width) + "+")
				for range height {
					fmt.Println("|" + strings.Repeat(" ", width) + "|")
				}
				fmt.Println("+" + strings.Repeat("-", width) + "+")
				println("")

				// Print out boxes for refrence
				for c, n := range projectData {
					println(c, strconv.Itoa(n[0])+":"+strconv.Itoa(n[1]))

					fmt.Println("+" + strings.Repeat("-", n[0]*width) + "+")
					for range n[1] * height {
						fmt.Println("|" + strings.Repeat(" ", n[0]*width) + "|")
					}
					fmt.Println("+" + strings.Repeat("-", n[0]*width) + "+")
					println("")
				}
			*/

			var index int = 0
			var viewW int = 4
			var bitMap uint64

			var done int
			var projectsDone = make(map[string]bool, len(projectData))

			// Inserts the blocks into the bitmap
		MapBlocks:
			for {
				if done == len(projectData) {
					break MapBlocks
				}

				if getBit(&bitMap, index) == 0 {
					space := nextBitChangeOrWall(&bitMap, viewW, index)

					round := 0
				FindBlock:
					for {
						for c, n := range projectData {
							// Finds a block that fits
							if n[0] == space-round && !projectsDone[c] {
								projectsDone[c] = true

								temp := projectData[c]
								temp[2] = (index - 1) * height
								projectData[c] = temp

								done++
								insertBlockIntoBitMap(&bitMap, [2]int{n[0], n[1]}, viewW, index)

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
			var renderBuf = make(map[int]string, (len(strconv.FormatUint(bitMap, 2)) * height))

			// Insert values into the render buffer
			for c, data := range projectData {
				strs := append([]string{c}, projectsMap[c]...)

				for i := range data[1] * height * data[0] {
					var str string

					switch true {
					case i/data[0] == 0:
						str = strs[i/data[0]]
					case i >= len(strs)*data[0]:
						continue
					default:
						str = " - " + strs[i/data[0]]
					}

					str = string([]rune(str)[min(i%data[0]*width, len(str)):min(i%data[0]*width+width, len(str))])
					if str == "" {
						continue
					}
					str = str + strings.Repeat(" ", width-len(str))

					str = "\033[38;5;" + strconv.Itoa(data[2]*100%255) + "m" + str + "\033[0m"

					renderBuf[(i/data[0]*viewW + data[2] + i%data[0])] = str
				}

			}

			// Writing every string piece to the renderBuilder for string building
			for i := range len(strconv.FormatUint(bitMap, 2)) * height {
				var str string

				switch true {
				case renderBuf[i] != "":
					str = renderBuf[i]
				default:
					str = strings.Repeat(" ", width)
				}
				str = "|" + str

				if (i+1)%viewW == 0 {
					str += "\n"
				}

				if i%(viewW*height) == 0 {
					str += "\n" + strings.Repeat("+"+strings.Repeat("-", width), viewW) + "+\n"
				}

				renderBuilder.WriteString(str)
			}

			fmt.Println("Info:", width, height, viewW)
			fmt.Println("Projects:", len(projectData))
			// Print the renderd string
			fmt.Println(renderBuilder.String())

			/*
				// Print grid with values
				println(strings.Repeat("+"+strings.Repeat("-", width), viewW) + "+")
				for y := range 12 {
					for h := range height {
						for i := range viewW {
							if h == int(math.Ceil(float64(height)/2)) {
								print("|" + strings.Repeat(" ", int(math.Floor(float64(width)/2))-1) + strconv.Itoa(int((bitMap>>uint((y*viewW)+i+1))&1)) + strings.Repeat(" ", int(math.Ceil(float64(width)/2))))
								continue
							}
							print("|" + strings.Repeat(" ", width))
						}
						print("|\n")
					}
					for range viewW {
						print("+" + strings.Repeat("-", width))
					}
					print("+\n")
				}
			*/
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

func nextBitChangeOrWall(bitMap *uint64, rowLength, index int) int {
	firstBit, firstline := getBit(bitMap, index), math.Floor(float64(index)/float64(rowLength))
	times := 0
	for {
		bit := getBit(bitMap, index)

		if math.Floor(float64(index)/float64(rowLength)) != firstline || bit != firstBit {
			break
		}
		index++
		times++
	}
	return times
}

func insertBlockIntoBitMap(bitMap *uint64, block [2]int, rowLength, index int) {
	for line := range block[0] {
		for column := range block[1] {
			*bitMap |= (1 << uint(index+(line*rowLength)+column))
		}
	}
}

func min(x, y int) int {
	if x > y {
		return y
	}
	return x
}
