package cmd

import (
	"fmt"
	"os"

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
		files, err := os.ReadDir(GetConfig().ProjectFolderPath)
		if err != nil {
			fmt.Println("Error listing projects:", err)
			return
		}
		fmt.Println("Projects:")
		for _, file := range files {
			if file.IsDir() {
				fmt.Println("-", file.Name())
			}
		}
	},
}
