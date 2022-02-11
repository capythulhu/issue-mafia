/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"os"

	"github.com/spf13/cobra"
	"github.com/thzoid/issue-mafia/util"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "issue-mafia",
	Short: "Fetch and update git hooks on repository",
	Long:  "Synchronize local git hooks with a remote repository\nspecified in the local .issue-mafia configuration file.",
	Run: func(cmd *cobra.Command, args []string) {
		recursive, _ := cmd.Root().Flags().GetBool("recursive")
		util.ScanRepos(recursive)
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	// Recursiveness flag
	rootCmd.Flags().BoolP("recursive", "r", false, "search repos recursively in subfolders")
}
