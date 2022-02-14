package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/thzoid/issue-mafia/util"
)

// removeCmd represents the remove command
var removeCmd = &cobra.Command{
	Use:   "remove",
	Short: "Remove syncronized hooks",
	Long:  "Delete all Git hooks synchronized by issue-mafia and its config file.",
	Run: func(cmd *cobra.Command, args []string) {
		currentPath := util.GetCurrentDir()
		dirIsRepo, dirHasConfig := util.IsRepo(currentPath), util.HasConfig(currentPath)
		// Check if local directory has config file
		if !dirHasConfig {
			util.ErrorLogger.Fatalln("current directory has no \u001b[90m.issue-mafia\u001b[0m config file.")
		}
		// Check if local directory is repository
		if !dirIsRepo {
			util.ErrorLogger.Fatalln("current directory is not a git repository.")
		}
		// Check if remote repository has hook files
		fmt.Println("Warning! This action will remove all hooks synchronized by issue-mafia, including its configuration file.\nDo you really want to proceed? \u001b[90m(\u001b[1mY\u001b[0m\u001b[90m/\u001b[1mn\u001b[0m\u001b[90m)\u001b[0m: \u001b[1m")
		var answer string
		fmt.Scanf("%s", &answer)
		fmt.Print("\u001b[0m")
		fmt.Println()
		// Validate answer
		if answer != "Y" && answer != "y" {
			util.WarningLogger.Fatalln("no changes made.")
		}

		util.RevertRepo(currentPath)
	},
}

func init() {
	rootCmd.AddCommand(removeCmd)
}
