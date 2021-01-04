package cliutil

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	"github.com/rs/xid"
)

// See https://stackoverflow.com/a/40891417
type ctxKey int

// Ctx* constants contain the keys for contexts with values.
const (
	CtxLogTags ctxKey = iota
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

// LogTagTransactionID is the log tag that contains the transaction ID.
const LogTagTransactionID = "transid"

// logLevels is a map of log level names to Log* constant.
var logLevels map[string]int

// MessageWriter defines a function the writes the log messages.
type MessageWriter func(ctx context.Context, logger *log.Logger, level string, message string, err error)

// LeveledLogger is a simple leveled logger that writes logs to STDOUT.
type LeveledLogger struct {
	level   int
	loggers []*log.Logger
	writer  MessageWriter
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

// LogLevelValid return true if the log level is valid.
func LogLevelValid(level string) (ok bool) {
	_, ok = logLevels[strings.ToLower(level)]
	return
}

// NewLoggerWithContext returns a leveled logger with a context that is
// initialized with a unique transaction ID.
func NewLoggerWithContext(ctx context.Context, level string) (context.Context, *LeveledLogger, xid.ID) {
	transid := xid.New()
	ctx = ContextWithLogTag(ctx, LogTagTransactionID, transid.String())
	return ctx, NewLogger(level), transid
}

// NewLogger returns a LeveledLogger that writes logs to either os.Stdout or
// ioutil.Discard depending on the passed minimum log level.
func NewLogger(level string) *LeveledLogger {
	logger := &LeveledLogger{
		level:  LogLevel(level),
		writer: DefaultMessageWriter,
	}

	flags := log.Ldate | log.Ltime | log.Lmicroseconds | log.LUTC
	logger.loggers = make([]*log.Logger, 5)
	for i := range logger.loggers {
		logger.loggers[i] = log.New(os.Stdout, "", flags)
	}

	return logger
}

// SetLevel sets the minimum log level. The log level defaults to "info" if the
// passed log level is not valid.
func (l *LeveledLogger) SetLevel(level string) {
	l.level = LogLevel(level)
}

// SetOutput sets the output for all loggers.
func (l *LeveledLogger) SetOutput(w io.Writer) {
	for _, logger := range l.loggers {
		logger.SetOutput(w)
	}
}

// SetFlags sets the output flags for all loggers.
func (l *LeveledLogger) SetFlags(flag int) {
	for _, logger := range l.loggers {
		logger.SetFlags(flag)
	}
}

// SetPrefix sets the prefix for all loggers.
func (l *LeveledLogger) SetPrefix(prefix string) {
	for _, logger := range l.loggers {
		logger.SetPrefix(prefix)
	}
}

// SetMessageWriter sets the MessageWriter for all loggers.
func (l *LeveledLogger) SetMessageWriter(fn MessageWriter) {
	l.writer = fn
}

// Fatal writes an fatal level log and exits with a non-zero exit code.
func (l LeveledLogger) Fatal(ctx context.Context, message string, err error) {
	l.printLog(ctx, l.level < LogLevelFatal, l.loggers[0], "FATAL", message, err)
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
	l.printLog(ctx, l.level < LogLevelError, l.loggers[1], "ERROR", message, err)
}

// ErrorIfError writes an error level log if err is not nil.
func (l LeveledLogger) ErrorIfError(ctx context.Context, message string, err error) {
	if err != nil {
		l.Error(ctx, message, err)
	}
}

// Notice writes an notice level log.
func (l LeveledLogger) Notice(ctx context.Context, message string) {
	l.printLog(ctx, l.level < LogLevelNotice, l.loggers[2], "NOTICE", message, nil)
}

// Info writes an info level log.
func (l LeveledLogger) Info(ctx context.Context, message string) {
	l.printLog(ctx, l.level < LogLevelInfo, l.loggers[3], "INFO", message, nil)
}

// Debug writes a debug level log.
func (l LeveledLogger) Debug(ctx context.Context, message string) {
	l.printLog(ctx, l.level < LogLevelDebug, l.loggers[4], "DEBUG", message, nil)
}

// printLog writes the log message using LeveledLogger.writer.
func (l *LeveledLogger) printLog(ctx context.Context, skip bool, logger *log.Logger, level string, message string, err error) {
	if !skip {
		l.writer(ctx, logger, level, message, err)
	}
}

// ContextWithLogTag returns a new context with log tags appended.
func ContextWithLogTag(ctx context.Context, key string, val interface{}) context.Context {
	if !IsLetters(key) {
		panic(fmt.Errorf("key must only contain letters: %q passed", key))
	}

	s := fmt.Sprintf("%v", val)
	var modifier string
	if HasSpace(s) {
		modifier = "%q"
	} else {
		modifier = "%s"
	}

	format := key + "=" + modifier
	tag := fmt.Sprintf(format, s)

	tags := ctx.Value(CtxLogTags)
	if tags != nil {
		tag = tags.(string) + " " + tag
	}

	return context.WithValue(ctx, CtxLogTags, tag)
}

// DefaultMessageWriter formats log messages according to Splunk's best
// practices. It is the default MessageWriter.
//
// See https://dev.splunk.com/enterprise/docs/developapps/addsupport/logging/loggingbestpractices/
func DefaultMessageWriter(ctx context.Context, logger *log.Logger, level string, message string, err error) {

	// Initialize the log message Printf() format.
	format := "%s message=%q"

	// Initialize the log message Printf() arguments.
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

	// Print the log message.
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
