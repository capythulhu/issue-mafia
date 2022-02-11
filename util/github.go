package util

var hooks = []string{
	"pre-commit",
	"prepare-commit-msg",
	"commit-msg",
	"post-commit",
	"applypatch-msg",
	"pre-applypatch",
	"post-applypatch",
	"pre-rebase",
	"post-rewrite",
	"post-checkout",
	"pre-merge-commit",
	"post-merge",
	"pre-push",
	"pre-auto-gc",
	"pre-receive",
	"update",
	"post-update",
	"post-receive",
	"fsmonitor-watchman",
}
