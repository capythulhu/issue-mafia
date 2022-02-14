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
	Long:  "Synchronize local git hooks with a remote repository specified in the local .issue-mafia configuration file.",
	Run: func(cmd *cobra.Command, args []string) {
		// Recursive flag
		recursive, _ := cmd.Root().Flags().GetBool("recursive")

		// Get available repositories
		paths := util.ScanDirs()
		if recursive {
			if len(paths) == 1 {
				util.InfoLogger.Println("updating 1 repository...")
			} else if len(paths) > 1 {
				util.InfoLogger.Println("updating", len(paths), "repositories...")
			}
			util.UpdateRepos(paths)
		} else {
			if len(paths) == 1 {
				util.InfoLogger.Println("updating current repository...")
			}
			// Update current repository (if avaliable)
			currentPath := util.GetCurrentDir()
			dirIsRepo, dirHasConfig, _ := util.UpdateRepo(currentPath)

			// Show error messages for current repository
			if !recursive && !dirHasConfig {
				if dirIsRepo {
					util.ErrorLogger.Fatalln("current directory has no \u001b[90m.issue-mafia\u001b[0m config file.")
				} else {
					util.ErrorLogger.Fatalln("current directory is not a git repository. if you want issue-mafia to look for repos in sub-directories, run \u001b[90missue-mafia --recursive\u001b[0m.")
				}
			}
		}
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
