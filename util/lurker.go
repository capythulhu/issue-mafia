package util

import (
	"os"
	"path/filepath"
)

// Returns if current directory is a repo
func isRepo(path string) bool {
	_, err := os.Stat(path + "/.git")
	return err == nil
}

// Returns if current directory has config file
func hasConfig(path string) bool {
	_, err := os.Stat(path + "/.issue-mafia")
	return err == nil
}

// Check if current directory is a git repository
func CurrentDirStats() (dirIsRepo, dirHasConfig bool) {
	ex, _ := os.Executable()
	exPath := filepath.Dir(ex)

	return isRepo(exPath), hasConfig(exPath)
}

// Scan repositories
func ScanRepos(recursive bool) {
	if recursive {

	} else {
		ex, _ := os.Executable()
		exPath := filepath.Dir(ex)

		if !isRepo(exPath) {
			ErrorLogger.Fatalln("current directory is not a git repository. if you want issue-mafia to look for repos in sub-directories, run \u001b[100m issue-mafia --recursive \u001b[0m.")
		} else if hasConfig(exPath) {
			ErrorLogger.Fatalln("current repository has no \u001b[100m .issue-mafia \u001b[0m config file. please, run \u001b[100m issue-mafia init \u001b[0m to setup a config file on this directory.")
		}
	}
}
