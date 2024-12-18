package logger

import (
	"log"
	"os"
)

var ErrLogger = log.New(os.Stderr, "bing-wall:", 0)
var InfoLogger = log.New(os.Stdout, "bing-wall:", 0)
