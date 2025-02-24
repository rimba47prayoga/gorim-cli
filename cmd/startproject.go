/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"gorim.org/gorim-cli/generator"
)

// startprojectCmd represents the startproject command
var startprojectCmd = &cobra.Command{
	Use:   "startproject",
	Short: "Start a new Gorim project",
	Long: `Creates a Gorim project directory structure for the given project name in the current directory or optionally in the given
directory.`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			fmt.Println("Project name is required")
			return
		}

		projectName := args[0]
		generator.StartProject(projectName)
	},
}

func init() {
	rootCmd.AddCommand(startprojectCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// startprojectCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// startprojectCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
