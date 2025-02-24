/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"gorim.org/gorim-cli/generator"
)

// startappCmd represents the startapp command
var startappCmd = &cobra.Command{
	Use:   "startapp",
	Short: "Start a new Gorim app",
	Long: `Creates a Gorim app directory structure for the given app name in the current directory or optionally in the given
directory.`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			fmt.Println("App name is required")
			return
		}

		appName := args[0]
		generator.StartApp(appName)
	},
}

func init() {
	rootCmd.AddCommand(startappCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// startappCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// startappCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
