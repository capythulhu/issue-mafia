package util

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strings"
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

// DownloadFile will download a url to a local file. It's efficient because it will
// write as it downloads and not load the whole file into memory.
func downloadFile(path string, url string) error {
	// Get the data
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Create the file
	out, err := os.Create(path)
	if err != nil {
		return err
	}
	defer out.Close()

	// Write the body to file
	_, err = io.Copy(out, resp.Body)
	return err
}

// Update individual repository
func UpdateRepo(path, repo, branch string) {
	files, status := FetchIntersectingFiles(repo, branch)
	if files == nil {
		ErrorLogger.Println("could not fetch files from", repo+". received status "+fmt.Sprintf("%d", status)+".")
	}

	for _, file := range files {
		completePath := path + "/.git/hooks/" + file
		downloadFile(completePath, "https://raw.githubusercontent.com/"+repo+"/"+branch+"/"+file)
		err := os.Chmod(completePath, 0700)
		if err != nil {
			log.Fatal(err)
		}
	}
}

// Read configuration file on directory
func readConfigFile(path string) (repo, branch string, ok bool) {
	// Read file
	content, err := os.ReadFile(path + "/.issue-mafia")
	if err != nil {
		ErrorLogger.Println("could not read config file at", path)
		return "", "", false
	}

	// Validate repository
	re := regexp.MustCompile(`^[-a-zA-Z0-9_]+\/[-a-zA-Z0-9_]+ [-a-zA-Z0-9_]+$`)
	if !re.Match(content) {
		ErrorLogger.Println("bad config file format at", path)
		return "", "", false
	}

	// Parse information
	info := strings.Split(string(content), " ")

	return info[0], info[1], true
}

// Update repositories
func UpdateRepos(recursive bool) {
	// Update current directory repository
	ex, _ := os.Executable()
	exPath := filepath.Dir(ex)
	if isRepo(exPath) {
		if hasConfig(exPath) {
			if repo, branch, ok := readConfigFile(exPath); ok {
				UpdateRepo(exPath, repo, branch)
			}
		} else if !recursive {
			ErrorLogger.Fatalln("current repository has no \u001b[100m .issue-mafia \u001b[0m config file. please, run \u001b[100m issue-mafia init \u001b[0m to setup a config file on this directory.")
		}
	} else if !recursive {
		ErrorLogger.Fatalln("current directory is not a git repository. if you want issue-mafia to look for repos in sub-directories, run \u001b[100m issue-mafia --recursive \u001b[0m.")
	}

	if recursive {

	}
}
