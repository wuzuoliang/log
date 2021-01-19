package log

import (
	"os"
	"time"

	"github.com/go-stack/stack"
)

const (
	timeKey     = "time"
	levelKey    = "level"
	msgKey      = "msg"
	locationKey = "location"
	errorKey    = "error"
)

type Level int

const (
	LvlAll = iota
	LvlFatal
	LvlError
	LvlWarn
	LvlInfo
	LvlDebug

	LvlAllStr = ""
	LvlFatalStr = "fatal"
	LvlErrorStr = "error"
	LvlWarnStr  = "warn"
	LvlInfoStr  = "info"
	LvlDebugStr = "debug"
)

var levelStringMap = map[Level]string{
	LvlAll:LvlAllStr,
	LvlFatal: LvlFatalStr,
	LvlError: LvlErrorStr,
	LvlWarn:  LvlWarnStr,
	LvlInfo:  LvlInfoStr,
	LvlDebug: LvlDebugStr,
}

// Returns the name of a Level
func (l Level) String() string {
	levelStr, ok := levelStringMap[l]
	if !ok {
		// _ErrLogLevelNotMatch
		panic("the log level not match")
	}
	return levelStr
}

// A Record is what a Logger asks its handler to write
type Record struct {
	Time         time.Time
	Level        Level
	Msg          string
	Ctx          []interface{}
	Call         stack.Call
	CustomCaller string
	KeyNames     RecordKeyNames
}

type RecordKeyNames struct {
	Time  string
	Msg   string
	Level string
	Call  string
}

// A Logger writes key/value pairs to a Handler
type Logger interface {
	// New returns a new Logger that has this logger's context plus the given context
	New(ctx ...interface{}) Logger

	// GetHandler gets the handler associated with the logger.
	GetHandler() Handler

	// SetHandler updates the logger to write records to the specified handler.
	SetHandler(h Handler)

	// Set level value. only level below this can be output
	SetOutLevel(l Level)
	GetOutLevel() Level

	// Log a message at the given level with context key/value pairs
	Debug(msg string, fields ...interface{})
	Info(msg string, fields ...interface{})
	Warn(msg string, fields ...interface{})
	Error(msg string, fields ...interface{})
	Fatal(msg string, fields ...interface{})
	Log(msg string,fields ...interface{})
}

type logger struct {
	ctx     []interface{}
	handler *swapHandler
	level   Level
}

type options struct{

}

func (l *logger) write(msg string, level Level, fields []interface{}) {
	if level <= l.level {
		l.handler.Log(&Record{
			Time:  time.Now(),
			Level: level,
			Msg:   msg,
			Ctx:   newContext(l.ctx, fields),
			Call:  stack.Caller(2),
			KeyNames: RecordKeyNames{
				Time:  timeKey,
				Msg:   msgKey,
				Level: levelKey,
				Call:  locationKey,
			},
		})
	}
}

func (l *logger) New(ctx ...interface{}) Logger {
	child := &logger{newContext(l.ctx, ctx), new(swapHandler), LvlDebug}
	child.SetHandler(l.handler)
	return child
}

func newContext(prefix []interface{}, suffix []interface{}) []interface{} {
	normalizedSuffix := normalize(suffix)
	newCtx := make([]interface{}, len(prefix)+len(normalizedSuffix))
	n := copy(newCtx, prefix)
	copy(newCtx[n:], normalizedSuffix)
	return newCtx
}

func (l *logger) SetOutLevel(level Level) {
	if level >= LvlFatal && level <= LvlDebug {
		l.level = level
	}
}

func (l *logger) GetOutLevel() Level {
	return l.level
}

func (l *logger) Debug(msg string, fields ...interface{}) {
	l.write(msg, LvlDebug, fields)
}

func (l *logger) Info(msg string, fields ...interface{}) {
	l.write(msg, LvlInfo, fields)
}

func (l *logger) Warn(msg string, fields ...interface{}) {
	l.write(msg, LvlWarn, fields)
}

func (l *logger) Error(msg string, fields ...interface{}) {
	l.write(msg, LvlError, fields)
}

func (l *logger) Fatal(msg string, fields ...interface{}) {
	l.write(msg, LvlFatal, fields)
	os.Exit(1)
}

func (l *logger) Log(msg string, fields ...interface{}) {
	l.write(msg, LvlAll, fields)
}

func (l *logger) GetHandler() Handler {
	return l.handler.Get()
}

func (l *logger) SetHandler(h Handler) {
	l.handler.Swap(h)
}

func normalize(ctx []interface{}) []interface{} {
	// if the caller passed a Ctx object, then expand it
	if len(ctx) == 1 {
		if ctxMap, ok := ctx[0].(Ctx); ok {
			ctx = ctxMap.toArray()
		}
	}

	// ctx needs to be even because it's a series of key/value pairs
	// no one wants to check for errors on logging functions,
	// so instead of erroring on bad input, we'll just make sure
	// that things are the right length and users can fix bugs
	// when they see the output looks wrong
	if len(ctx)%2 != 0 {
		ctx = append(ctx, nil, errorKey, "Normalized odd number of arguments by adding nil")
	}

	return ctx
}

// Lazy allows you to defer calculation of a logged value that is expensive
// to compute until it is certain that it must be evaluated with the given filters.
//
// Lazy may also be used in conjunction with a Logger's New() function
// to generate a child logger which always reports the current value of changing
// state.
//
// You may wrap any function which takes no arguments to Lazy. It may return any
// number of values of any type.
type Lazy struct {
	Fn interface{}
}

// Ctx is a map of key/value pairs to pass as context to a log function
// Use this only if you really need greater safety around the arguments you pass
// to the logging functions.
type Ctx map[string]interface{}

func (c Ctx) toArray() []interface{} {
	arr := make([]interface{}, len(c)*2)

	i := 0
	for k, v := range c {
		arr[i] = k
		arr[i+1] = v
		i += 2
	}

	return arr
}
