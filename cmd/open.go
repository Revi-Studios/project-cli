package cmd

import (
	"fmt"
	"os/exec"

	"github.com/spf13/cobra"
)

var OpenCmd = &cobra.Command{
	Use:   "open",
	Short: "Open a project",
	Long:  "Open a project",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			fmt.Println("Please provide a project name")
			return
		}
		projectName := args[0]
		projectPath := GetConfig().ProjectFolderPath + "/" + projectName + "/"
		if projectPath == "" {
			fmt.Printf("Project %s not found\n", projectName)
			return
		}
		err := exec.Command("osascript", "-e", fmt.Sprintf(`tell application "Terminal" to do script "cd '%s'; clear"`, projectPath)).Run()
		if err != nil {
			fmt.Printf("Error opening project %s: %v\n", projectName, err)
			return
		}
		fmt.Printf("Opened project %s\n", projectName)
	},
}

func init() {
	rootCmd.AddCommand(OpenCmd)
}
