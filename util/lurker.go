package util

import (
	"fmt"
	"io"
	"io/fs"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

// Returns if current directory is a repo
func IsRepo(path string) bool {
	_, err := os.Stat(path + "/.git")
	return err == nil
}

// Returns if current directory has config file
func HasConfig(path string) bool {
	_, err := os.Stat(path + "/.issue-mafia")
	return err == nil
}

// Get current directory path
func GetCurrentDir() string {
	ex, _ := os.Executable()
	return filepath.Dir(ex)
}

// Download single file and attribute specific permissions to it
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
func UpdateRepo(path string) (dirIsRepo, dirHasConfig bool) {
	// Check if directory is repository and if it has config file
	dirIsRepo, dirHasConfig = IsRepo(path), HasConfig(path)
	if !dirIsRepo || !dirHasConfig {
		if dirHasConfig && !dirIsRepo {
			ErrorLogger.Println(path, "has an \u001b[100m .issue-mafia \u001b[0m config file, but is not a git repository.")
		}
		return
	}

	// Read configuration file
	repo, branch, ok := readConfigFile(path)
	if !ok {
		return
	}

	// Fetch files from remote repo
	files, status := FetchIntersectingFiles(repo, branch)
	if files == nil {
		ErrorLogger.Println("could not fetch files from", repo+". received status "+fmt.Sprintf("%d", status)+".")
	}

	// Download files from remote repo
	for _, file := range files {
		completePath := path + "/.git/hooks/" + file
		downloadFile(completePath, "https://raw.githubusercontent.com/"+repo+"/"+branch+"/"+file)
		err := os.Chmod(completePath, 0700)
		if err != nil {
			log.Fatal(err)
		}
	}

	// Log success
	InfoLogger.Println(path, "hooks synchronized from github.com/"+repo)

	return
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
func UpdateRepos() {
	currentPath := GetCurrentDir()
	err := filepath.WalkDir(currentPath,
		func(path string, d fs.DirEntry, err error) error {
			if err != nil {
				return err
			}
			if path != currentPath {
				UpdateRepo(path)
			}
			return nil
		})
	if err != nil {
		log.Println(err)
	}
}
