package util

import (
	"fmt"
	"io"
	"io/fs"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"sync"
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
	// ex, _ := os.Executable()
	// return filepath.Dir(ex)
	return "."
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
func UpdateRepo(path string) (dirIsRepo, dirHasConfig, ok bool) {
	// Check if directory is repository and if it has config file
	dirIsRepo, dirHasConfig = IsRepo(path), HasConfig(path)
	if !dirIsRepo || !dirHasConfig {
		if dirHasConfig && !dirIsRepo {
			ErrorLogger.Println(path, "has an \u001b[90m.issue-mafia\u001b[0m config file, but is not a git repository.")
		}
		return dirIsRepo, dirHasConfig, false
	}

	// Read configuration file
	repo, branch, ok := readConfigFile(path)
	if !ok {
		return dirIsRepo, dirHasConfig, false
	}

	// Fetch files from remote repo
	files, status := FetchIntersectingFiles(repo, branch)
	if files == nil {
		ErrorLogger.Println("could not fetch files from", repo+". received status "+fmt.Sprintf("%d", status)+".")
	}

	// Wait group for download synchronization
	var wg sync.WaitGroup

	// Download files from remote repo
	for _, file := range files {
		wg.Add(1)
		go func(file string) {
			defer wg.Done()
			completePath := path + "/.git/hooks/" + file
			downloadFile(completePath, "https://raw.githubusercontent.com/"+repo+"/"+branch+"/"+file)
			os.Chmod(completePath, 0700)
		}(file)
	}

	// Wait for downloads to be finished
	wg.Wait()

	// Log success
	var readablePath string
	if path == "." {
		readablePath = "current repository"
	} else {
		readablePath = path
	}

	InfoLogger.Println(readablePath, "hooks synchronized from github.com/"+repo)

	return dirIsRepo, dirHasConfig, true
}

// Revert individual repository
func RevertRepo(path string) {
	// Check if directory is repository and if it has config file
	dirIsRepo, dirHasConfig := IsRepo(path), HasConfig(path)
	if !dirIsRepo || !dirHasConfig {
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

	// Wait group for download synchronization
	var wg sync.WaitGroup

	// Download files from remote repo
	for _, file := range files {
		wg.Add(1)
		go func(file string) {
			defer wg.Done()
			os.Remove(path + "/.git/hooks/" + file)
		}(file)
	}

	// Wait for downloads to be finished
	wg.Wait()

	InfoLogger.Println("hooks from github.com/"+repo, "removed successfully.")

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

// Scan directories in folder
func ScanDirs() []string {
	paths := []string{}
	currentPath := GetCurrentDir()
	filepath.WalkDir(currentPath,
		func(path string, d fs.DirEntry, err error) error {
			if err != nil {
				return err
			}
			if IsRepo(path) && HasConfig(path) && currentPath != path {
				paths = append(paths, path)
			}
			return nil
		})

	return paths
}

// Update repositories
func UpdateRepos(paths []string) {
	// Get current path
	currentPath := GetCurrentDir()

	// Wait group for download synchronization
	var wg sync.WaitGroup

	// Recursively update repositories
	updatedRepos := 0
	for _, path := range paths {
		if path != currentPath {
			wg.Add(1)
			go func() {
				defer wg.Done()
				_, _, ok := UpdateRepo(path)
				if ok {
					updatedRepos++
				}
			}()
		}
	}

	// Wait for repositories to be updated
	wg.Wait()

	// Show final log
	if updatedRepos == 0 {
		WarningLogger.Println("no issue-mafia repositories found on sub-directories.")
	} else if updatedRepos == 1 {
		InfoLogger.Println("1 repository synchronized successfully.")
	} else {
		InfoLogger.Println(updatedRepos, " repositories synchronized successfully.")
	}
}
