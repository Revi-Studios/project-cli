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

		fmt.Println("Projects:")
		os.Chdir(config.ProjectFolderPath + "/")

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
		return nil
	},
}
