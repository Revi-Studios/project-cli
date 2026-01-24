package cmd

import (
	"fmt"
	"os"
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
	Run: func(cmd *cobra.Command, args []string) {
		files, err := os.ReadDir(lib.GetConfig().ProjectFolderPath + "/")
		if err != nil {
			fmt.Println("Error listing projects:", err)
			return
		}
		fmt.Println("Projects:")
		os.Chdir(lib.GetConfig().ProjectFolderPath + "/")

		var projects [][2]string
		var longestFileName int

		for _, file := range files {
			if file.IsDir() {
				name := file.Name()
				tags, err := lib.GetTags(name)
				if err != nil {
				}
				if len := utf8.RuneCountInString(name); len > longestFileName {
					longestFileName = len
				}
				projects = append(projects, [2]string{name, tags})

			}
		}
		for _, prj := range projects {
			strLeft := strings.Repeat(" ", longestFileName-utf8.RuneCountInString(prj[0]))
			fmt.Println(prj[0], strLeft, "|", prj[1])
		}

	},
}
