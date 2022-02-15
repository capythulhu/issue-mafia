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
		// Formatting flag
		hard, _ := cmd.Flags().GetBool("hard")

		// Check if remote repository has hook files
		fmt.Println("Warning! This action will remove all hooks synchronized by issue-mafia,\nincluding its configuration file (if present).\nDo you really want to proceed? \u001b[90m(\u001b[1mY\u001b[0m\u001b[90m/\u001b[1mn\u001b[0m\u001b[90m)\u001b[0m: \u001b[1m")
		var answer string
		fmt.Scanf("%s", &answer)
		fmt.Print("\u001b[0m")
		fmt.Println()
		// Validate answer
		if answer != "Y" && answer != "y" {
			util.WarningLogger.Fatalln("no changes made.")
		}

		util.CleanRepo(".", hard)
	},
}

func init() {
	// Formatting flag
	removeCmd.Flags().Bool("hard", false, "remove all hooks on repository")

	rootCmd.AddCommand(removeCmd)
}
