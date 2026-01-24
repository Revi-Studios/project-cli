package cmd

import (
	"fmt"
	"os"

	"github.com/Revi-Studios/project/lib"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(add)
}

var add = &cobra.Command{
	Use:   "add",
	Short: "Add a new project",
	Long:  `Add a new project to the project list`,
	Run: func(cmd *cobra.Command, args []string) {
		switch {
		case len(args) == 0:
			fmt.Println("Please provide a project name and tag")
			return
		case len(args) == 1:
			fmt.Println("Please provide a tag")
			return
		}
		projectName := args[0]
		projectPath := lib.GetConfig().ProjectFolderPath + "/" + projectName
		err := os.Mkdir(projectPath, 0755)
		if err != nil {
			fmt.Println("Error creating project:", err)
			return
		}
		if err := lib.SetTag(projectPath, args[1]); err != nil {
			fmt.Println("Error adding tag:", err)
			return
		}
		fmt.Println("Project created:", projectPath)
	},
}
