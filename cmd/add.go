package cmd

import (
	"fmt"
	"os"

	"github.com/Revi-Studios/project-cli/config"
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
		if len(args) == 0 {
			fmt.Println("Please provide a project name")
			return
		}
		projectName := args[0]
		projectPath := config.GetConfig().ProjectFolderPath + "/" + projectName
		err := os.Mkdir(projectPath, 0755)
		if err != nil {
			fmt.Println("Error creating project:", err)
			return
		}
		fmt.Println("Project created:", projectPath)
	},
}
