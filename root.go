package log

import (
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

	root = &logger{[]interface{}{}, new(swapHandler), LvlAll}
	root.SetHandler(StdoutHandler)
}

// New returns a new logger with the given context.
// New is a convenient alias for Root().New
func New(ctx ...interface{}) Logger {
	return root.New(ctx...)
}

// New returns a new logger with the given context.
// New is a convenient alias for Root().New
func NewWithOptions(options []) Logger {
	return root.New(ctx...)
}


// Root returns the root logger
func Root() Logger {
	return root
}

func SetOutLevel(level Level) {
	root.SetOutLevel(level)
}

func GetLogLevel() Level {
	return root.level
}

// The following functions bypass the exported logger methods (logger.Debug,
// etc.) to keep the call depth the same for all paths to logger.write so
// runtime.Caller(2) always refers to the call site in client code.

func IsDebugEnable() bool {
	return root.level >= LvlDebug
}

// Debug is a convenient alias for Root().Debug
func Debug(msg string, kvalues ...interface{}) {
	root.write(msg, LvlDebug, kvalues)
}

func IsInfoEnable() bool {
	return root.level >= LvlInfo
}

// Info is a convenient alias for Root().Info
func Info(msg string, kvalues ...interface{}) {
	root.write(msg, LvlInfo, kvalues)
}

func IsWarnEnable() bool {
	return root.level >= LvlWarn
}

// Warn is a convenient alias for Root().Warn
func Warn(msg string, kvalues ...interface{}) {
	root.write(msg, LvlWarn, kvalues)
}

func IsErrorEnable() bool {
	return root.level >= LvlError
}

// Error is a convenient alias for Root().Error
func Error(msg string, kvalues ...interface{}) {
	root.write(msg, LvlError, kvalues)
}

func IsFatalEnable() bool {
	return root.level >= LvlFatal
}

// Fatal is a convenient alias for Root().Fatal
func Fatal(msg string, kvalues ...interface{}) {
	root.write(msg, LvlFatal, kvalues)
	os.Exit(1)
}

// Log is a convenient alias for Root().Log
func Log(msg string,kvalues ...interface{}) {
	root.write(msg,LvlAll,kvalues)
}
