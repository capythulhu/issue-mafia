package util

import (
	"encoding/json"
	"io"
	"net/http"
)

var hooks = map[string]struct{}{
	"pre-commit":         {},
	"prepare-commit-msg": {},
	"commit-msg":         {},
	"post-commit":        {},
	"applypatch-msg":     {},
	"pre-applypatch":     {},
	"post-applypatch":    {},
	"pre-rebase":         {},
	"post-rewrite":       {},
	"post-checkout":      {},
	"pre-merge-commit":   {},
	"post-merge":         {},
	"pre-push":           {},
	"pre-auto-gc":        {},
	"pre-receive":        {},
	"update":             {},
	"post-update":        {},
	"post-receive":       {},
	"fsmonitor-watchman": {},
}

func FetchRepository(repo string) int {
	// Fetch repository from GitHub
	resp, _ := http.Get("https://github.com/" + repo)
	return resp.StatusCode
}

func FetchIntersectingFiles(repo string, branch string) ([]string, int) {
	// Fetch repository from GitHub API
	resp, _ := http.Get("https://api.github.com/repos/" + repo + "/git/trees/" + branch)
	if resp.StatusCode != 200 {
		return nil, resp.StatusCode
	}

	// Read files
	b, _ := io.ReadAll(resp.Body)
	// Body variable
	body := struct {
		Tree []struct {
			Path string
		}
	}{}
	// Unmarshal json
	json.Unmarshal(b, &body)

	// Insert files in slice
	files := []string{}
	for _, file := range body.Tree {
		if _, ok := hooks[file.Path]; ok {
			files = append(files, file.Path)
		}
	}

	return files, resp.StatusCode
}
