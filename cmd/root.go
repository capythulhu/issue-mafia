/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "issue-mafia",
	Short: "Git hooks enforcer in local repositories",
	Long: `issue-mafia is a Go application that locally enforces
Git hooks fetched from a remote repository. That way,
hooks can be easily managed and updated.
`,
	Run: func(cmd *cobra.Command, args []string) {

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
	var recursive bool
	rootCmd.PersistentFlags().BoolVar(&recursive, "recursive", false, "search repos recursively in subfolders (default is false)")
}
