package cliutil

import (
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"github.com/rs/xid"
)

// See https://stackoverflow.com/a/40891417
type ctxKey int

// Ctx* constants contain the keys for contexts with values.
const (
	CtxCmdArgs ctxKey = iota
	CtxLogTags
)

// Log* constants represent the log levels as strings for configuration.
const (
	LogNone   = "none"
	LogFatal  = "fatal"
	LogError  = "error"
	LogNotice = "notice"
	LogInfo   = "info"
	LogDebug  = "debug"
)

// LogTag* constants contain common log tags.
const (
	LogTagFile          = "file"
	LogTagTransactionID = "transid"
	LogTagURL           = "url"
)

// logLevels is a map of log level names to Log* constant.
var logLevels map[string]int

// LeveledLogger is a simple leveled logger that writes logs to STDOUT.
type LeveledLogger struct {
	FatalLogger  *log.Logger
	ErrorLogger  *log.Logger
	NoticeLogger *log.Logger
	InfoLogger   *log.Logger
	DebugLogger  *log.Logger
}

// LogLevel* represents log levels as integers for comparrison.
const (
	LogLevelNone = iota
	LogLevelFatal
	LogLevelError
	LogLevelNotice
	LogLevelInfo
	LogLevelDebug
)

// LogLevel returns the log level's integer representation.
func LogLevel(level string) (id int) {
	var ok bool
	level = strings.ToLower(level)
	if id, ok = logLevels[level]; !ok {
		id = LogLevelInfo
	}
	return
}

// NewLoggerWithContext returns a leveled logger with a context that is
// initialized with a unique transaction ID.
func NewLoggerWithContext(tx context.Context, level string) (ctx context.Context, logger *LeveledLogger, transid xid.ID) {
	transid = xid.New()
	ctx = context.WithValue(ctx, CtxLogTags, "")
	ctx = ContextWithLogTag(ctx, LogTagTransactionID, transid.String())
	logger = NewLogger(level)
	return
}

// NewLogger returns a LeveledLogger that writes logs to either os.Stdout or
// ioutil.Discard depending on the passed minimum log level.
func NewLogger(level string) *LeveledLogger {

	// Use os.Stdout for all numbers >= the log level's ID, ioutil.Discard for
	// everything < the log level's ID.
	id := LogLevel(level)
	w := make([]io.Writer, 6)
	for i := range w {
		if i <= id {
			w[i] = os.Stdout
		} else {
			w[i] = ioutil.Discard
		}
	}

	flags := log.Ldate | log.Ltime | log.Lmicroseconds | log.LUTC
	return &LeveledLogger{
		FatalLogger:  log.New(w[LogLevelFatal], "", flags),
		ErrorLogger:  log.New(w[LogLevelError], "", flags),
		NoticeLogger: log.New(w[LogLevelNotice], "", flags),
		InfoLogger:   log.New(w[LogLevelInfo], "", flags),
		DebugLogger:  log.New(w[LogLevelDebug], "", flags),
	}
}

// Fatal writes an fatal level log and exits with a non-zero exit code.
func (l LeveledLogger) Fatal(ctx context.Context, message string, err error) {
	printLog(ctx, l.FatalLogger, "FATAL", message, err)
	os.Exit(1)
}

// FatalIfError writes a fatal level log and exits with a non-zero exit code if
// err != nil. This function is a no-op if err == nil.
func (l LeveledLogger) FatalIfError(ctx context.Context, message string, err error) {
	if err != nil {
		l.Fatal(ctx, message, err)
	}
}

// Error writes an error level log.
func (l LeveledLogger) Error(ctx context.Context, message string, err error) {
	printLog(ctx, l.ErrorLogger, "ERROR", message, err)
}

// ErrorIfError writes an error level log if err is not nil.
func (l LeveledLogger) ErrorIfError(ctx context.Context, message string, err error) {
	if err != nil {
		l.Error(ctx, message, err)
	}
}

// Notice writes an notice level log.
func (l LeveledLogger) Notice(ctx context.Context, message string) {
	printLog(ctx, l.NoticeLogger, "NOTICE", message, nil)
}

// Info writes an info level log.
func (l LeveledLogger) Info(ctx context.Context, message string) {
	printLog(ctx, l.InfoLogger, "INFO", message, nil)
}

// Debug writes a debug level log.
func (l LeveledLogger) Debug(ctx context.Context, message string) {
	printLog(ctx, l.DebugLogger, "DEBUG", message, nil)
}

// ContextWithLogTag returns a new context with log tags appended.
func ContextWithLogTag(ctx context.Context, key string, val string) context.Context {
	if !IsLetters(key) {
		panic(fmt.Errorf("key must only contain letters: %q passed", key))
	}

	var modifier string
	if HasSpace(val) {
		modifier = "%q"
	} else {
		modifier = "%s"
	}

	format := key + "=" + modifier
	tag := fmt.Sprintf(format, val)

	tags := ctx.Value(CtxLogTags)
	if tags != nil {
		tag = tags.(string) + " " + tag
	}

	return context.WithValue(ctx, CtxLogTags, tag)
}

func printLog(ctx context.Context, logger *log.Logger, level string, message string, err error) {
	format := "%s message=%q"

	args := make([]interface{}, 2)
	args[0] = level
	args[1] = message

	// Append the error if there is one.
	if err != nil {
		format = format + " error=%q"
		args = append(args, err)
	}

	// Append the log tags if there are any.
	tags := ctx.Value(CtxLogTags)
	if tags != nil {
		format = format + " %s"
		args = append(args, tags)
	}

	logger.Printf(format, args...)
}

func init() {
	logLevels = make(map[string]int, 6)
	logLevels[LogNone] = LogLevelNone
	logLevels[LogFatal] = LogLevelFatal
	logLevels[LogError] = LogLevelError
	logLevels[LogNotice] = LogLevelNotice
	logLevels[LogInfo] = LogLevelInfo
	logLevels[LogDebug] = LogLevelDebug
}
