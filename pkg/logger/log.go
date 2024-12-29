package logger

import (
	"io"
	"log"
	"os"
	"strings"

	"github.com/hashicorp/go-retryablehttp"
)

const (
	LogLevelDebug logLevel = iota
	LogLevelInfo
	LogLevelWarn
	LogLevelError
)

var Logger = &logger{logger: log.New(LogLevelInfo, "bing-wall [DEBUG]:\t", log.Ldate|log.Ltime|log.Lmsgprefix)}

var _ retryablehttp.LeveledLogger = logger{}
var _ io.Writer = logLevel(0)

type (
	logLevel int
	logger   struct {
		logger *log.Logger
	}
)

func (l logLevel) String() string {
	return map[logLevel]string{
		LogLevelDebug: "DEBUG",
		LogLevelInfo:  "INFO",
		LogLevelWarn:  "WARN",
		LogLevelError: "ERROR",
	}[l]
}

func (l logLevel) Replace(in string) string {
	return strings.NewReplacer(
		LogLevelDebug.String(), l.String(),
		LogLevelInfo.String(), l.String(),
		LogLevelWarn.String(), l.String(),
		LogLevelError.String(), l.String(),
	).Replace(in)
}

func (l logLevel) Write(p []byte) (n int, err error) {
	if l == LogLevelError {
		return os.Stderr.Write(p)
	}

	return os.Stdout.Write(p)
}

func (l logger) Debug(msg string, args ...any) { l.Logln(LogLevelDebug, 0, toArgs(msg, args)...) }
func (l logger) Error(msg string, args ...any) { l.Logln(LogLevelError, 0, toArgs(msg, args)...) }
func (l logger) Info(msg string, args ...any)  { l.Logln(LogLevelInfo, 0, toArgs(msg, args)...) }
func (l logger) Warn(msg string, args ...any)  { l.Logln(LogLevelWarn, 0, toArgs(msg, args)...) }

func (l logger) Print(args ...any)              { l.Log(LogLevelInfo, 0, args...) }
func (l logger) Printf(msg string, args ...any) { l.Logf(LogLevelInfo, msg, 0, args...) }
func (l logger) Println(args ...any)            { l.Logln(LogLevelInfo, 0, args...) }
func (l logger) Fatal(args ...any)              { l.Log(LogLevelError, 1, args...) }
func (l logger) Fatalf(msg string, args ...any) { l.Logf(LogLevelError, msg, 1, args...) }
func (l logger) Fatalln(args ...any)            { l.Logln(LogLevelError, 1, args...) }
func (l logger) Panic(args ...any)              { l.Log(LogLevelError, 2, args...) }
func (l logger) Panicf(msg string, args ...any) { l.Logf(LogLevelError, msg, 2, args...) }
func (l logger) Panicln(args ...any)            { l.Logln(LogLevelError, 2, args...) }

func (l logger) SetLevel(lvl logLevel) { l.logger.SetOutput(lvl) }

func (l logger) Log(lvl logLevel, code int, args ...any) {
	selectLogAction(l, lvl, code, l.logger.Fatal, l.logger.Panic, l.logger.Print, "", args...)
}

func (l logger) Logf(lvl logLevel, msg string, code int, args ...any) {
	selectLogAction(l, lvl, code, l.logger.Fatalf, l.logger.Panicf, l.logger.Printf, msg, args...)
}

func (l logger) Logln(lvl logLevel, code int, args ...any) {
	selectLogAction(l, lvl, code, l.logger.Fatalln, l.logger.Panicln, l.logger.Println, "", args...)
}

func selectLogAction[T interface {
	func(...any) | func(string, ...any)
}](l logger, lvl logLevel, code int, fatal, panic, print T, msg string, args ...any) {
	if wLvl, ok := l.logger.Writer().(logLevel); wLvl > lvl || !ok {
		return
	}

	var selected T
	switch code {
	case 1:
		selected = fatal
	case 2:
		selected = panic
	default:
		selected = print
	}

	l.logger.SetPrefix(lvl.Replace(l.logger.Prefix()))
	switch action := any(selected).(type) {
	case func(...any):
		action(args...)
	case func(string, ...any):
		action(msg, args...)
	}
}

func toArgs(msg string, args []any) []any {
	return append([]any{msg}, args...)
}
