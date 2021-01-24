package log

import (
	"context"
	"github.com/inconshreveable/log15/term"
	"github.com/mattn/go-colorable"
	"os"
)

var (
	root          *logger
	StdoutHandler = StreamHandler(os.Stdout, LogfmtFormat())
	StderrHandler = StreamHandler(os.Stderr, LogfmtFormat())
)

func init() {
	if term.IsTty(os.Stdout.Fd()) {
		StdoutHandler = StreamHandler(colorable.NewColorableStdout(), TerminalFormat())
	}

	if term.IsTty(os.Stderr.Fd()) {
		StderrHandler = StreamHandler(colorable.NewColorableStderr(), TerminalFormat())
	}

	root = &logger{[]interface{}{}, new(swapHandler), LvlTrace}
	root.SetHandler(StdoutHandler)
}

// New returns a new logger with the given keyValues.
// New is a convenient alias for Root().New
func New(keyValues ...interface{}) Logger {
	return root.New(keyValues...)
}

// Root returns the root logger
func Root() Logger {
	return root
}

func SetOutLevel(level Level) {
	root.SetOutLevel(level)
}

func GetLogLevel() Level {
	return root.GetOutLevel()
}

// The following functions bypass the exported logger methods (logger.Debug,
// etc.) to keep the call depth the same for all paths to logger.write so
// runtime.Caller(2) always refers to the call site in client code.

// Log is a convenient alias for Root().Log
func Log(msg string, keyValues ...interface{}) {
	root.write(msg, LvlTrace, keyValues)
}

// LogContext is a convenient alias for Root().Log with context.Context
func LogContext(ctx context.Context, msg string, keyValues ...interface{}) {
	root.writeContext(ctx, msg, LvlTrace, keyValues)
}

func IsDebugEnable() bool {
	return root.level >= LvlDebug
}

// Debug is a convenient alias for Root().Debug
func Debug(msg string, keyValues ...interface{}) {
	root.write(msg, LvlDebug, keyValues)
}

// Debug is a convenient alias for Root().Debug
func DebugContext(ctx context.Context, msg string, keyValues ...interface{}) {
	root.writeContext(ctx, msg, LvlDebug, keyValues)
}

func IsInfoEnable() bool {
	return root.level >= LvlInfo
}

// Info is a convenient alias for Root().Info
func Info(msg string, keyValues ...interface{}) {
	root.write(msg, LvlInfo, keyValues)
}

// Info is a convenient alias for Root().Info
func InfoContext(ctx context.Context, msg string, keyValues ...interface{}) {
	root.writeContext(ctx, msg, LvlInfo, keyValues)
}

func IsWarnEnable() bool {
	return root.level >= LvlWarn
}

// Warn is a convenient alias for Root().Warn
func Warn(msg string, keyValues ...interface{}) {
	root.write(msg, LvlWarn, keyValues)
}

// Warn is a convenient alias for Root().Warn
func WarnContext(ctx context.Context, msg string, keyValues ...interface{}) {
	root.writeContext(ctx, msg, LvlWarn, keyValues)
}

func IsErrorEnable() bool {
	return root.level >= LvlError
}

// Error is a convenient alias for Root().Error
func Error(msg string, keyValues ...interface{}) {
	root.write(msg, LvlError, keyValues)
}

// Error is a convenient alias for Root().Error
func ErrorContext(ctx context.Context, msg string, keyValues ...interface{}) {
	root.writeContext(ctx, msg, LvlError, keyValues)
}

func IsFatalEnable() bool {
	return root.level >= LvlFatal
}

// Fatal is a convenient alias for Root().Fatal
func Fatal(msg string, keyValues ...interface{}) {
	root.write(msg, LvlFatal, keyValues)
	os.Exit(1)
}

// Fatal is a convenient alias for Root().Fatal
func FatalContext(ctx context.Context, msg string, keyValues ...interface{}) {
	root.writeContext(ctx, msg, LvlFatal, keyValues)
}
