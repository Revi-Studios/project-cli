package cmd

import "github.com/spf13/cobra"

func init() {
	rootCmd.AddCommand(quick)
}

var quick = &cobra.Command{
	Use:   "quick",
	Short: "Open the Quick opener",
	Long:  "Open the Quick opener that makes it quicker to open projects",
	RunE: func(cmd *cobra.Command, args []string) error {
		return nil
	},
}
