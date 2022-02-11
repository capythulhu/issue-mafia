package util

import (
	"log"
	"os"
)

var (
	InfoLogger    = log.New(os.Stderr, "\u001b[36m[i]\u001b[0m ", 0)
	WarningLogger = log.New(os.Stderr, "\u001b[33m[!]\u001b[0m ", 0)
	ErrorLogger   = log.New(os.Stderr, "\u001b[31m[x]\u001b[0m ", 0)
)
