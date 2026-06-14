package cmd

import "github.com/spf13/cobra"

func init() {
	rootCmd.AddCommand(config)
	config.AddCommand(configSet)
	config.AddCommand(configGet)
}

var config = &cobra.Command{}

var configSet = &cobra.Command{}

var configGet = &cobra.Command{}
