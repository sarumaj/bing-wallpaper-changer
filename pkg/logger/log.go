package logger

import (
	"log"
	"os"
)

var ErrLogger = log.New(os.Stderr, "bing-wall:", 0)
