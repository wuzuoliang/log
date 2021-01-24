package log

import (
	"context"
	"os"
	"time"

	"github.com/go-stack/stack"
)

const (
	/**
	日志关键字段，Key信息
	*/
	timeKey     = "time"
	levelKey    = "level"
	msgKey      = "msg"
	locationKey = "location"
	errorKey    = "error"
	requestID   = "request_id"
)

// Level 日志级别
type Level int

const (
	LvlFatal = iota
	LvlError
	LvlWarn
	LvlInfo
	LvlDebug
	LvlTrace

	LvlFatalStr = "fatal"
	LvlErrorStr = "error"
	LvlWarnStr  = "warn"
	LvlInfoStr  = "info"
	LvlDebugStr = "debug"
	LvlTraceStr = "trace"
)

var levelStringMap = map[Level]string{
	LvlTrace: LvlTraceStr,
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
	KeyValues    []interface{}
	Ctx          context.Context
	Call         stack.Call
	CustomCaller string
	KeyNames     RecordKeyNames
}

// RecordKeyNames 日志记录规则字段名
type RecordKeyNames struct {
	Time      string
	Msg       string
	Level     string
	Call      string
	RequestID string
}

// A Logger writes key/value pairs to a Handler
type Logger interface {
	// New returns a new Logger that has this logger's context plus the given context
	New(keyValues ...interface{}) Logger

	// GetHandler gets the handler associated with the logger.
	GetHandler() Handler

	// SetHandler updates the logger to write records to the specified handler.
	SetHandler(h Handler)

	// Set level value. only level below this can be output
	SetOutLevel(l Level)
	GetOutLevel() Level

	// Log a message at the given level with context key/value pairs
	Log(msg string, fields ...interface{})
	Debug(msg string, fields ...interface{})
	Info(msg string, fields ...interface{})
	Warn(msg string, fields ...interface{})
	Error(msg string, fields ...interface{})
	Fatal(msg string, fields ...interface{})

	// LogContext a message at the given level with context key/value pairs
	LogContext(ctx context.Context, msg string, fields ...interface{})
	DebugContext(ctx context.Context, msg string, fields ...interface{})
	InfoContext(ctx context.Context, msg string, fields ...interface{})
	WarnContext(ctx context.Context, msg string, fields ...interface{})
	ErrorContext(ctx context.Context, msg string, fields ...interface{})
	FatalContext(ctx context.Context, msg string, fields ...interface{})
}

type logger struct {
	KeyValues []interface{}
	handler   *swapHandler
	level     Level
}

func (l *logger) write(msg string, level Level, fields []interface{}) {
	if level <= l.level {
		l.handler.Log(&Record{
			Time:      time.Now(),
			Level:     level,
			Msg:       msg,
			KeyValues: newKeyValues(l.KeyValues, fields),
			Call:      stack.Caller(2),
			KeyNames: RecordKeyNames{
				Time:      timeKey,
				Msg:       msgKey,
				Level:     levelKey,
				Call:      locationKey,
				RequestID: requestID,
			},
		})
	}
}

func (l *logger) writeContext(ctx context.Context, msg string, level Level, fields []interface{}) {
	if level <= l.level {
		l.handler.Log(&Record{
			Ctx:       ctx,
			Time:      time.Now(),
			Level:     level,
			Msg:       msg,
			KeyValues: newKeyValues(l.KeyValues, fields),
			Call:      stack.Caller(2),
			KeyNames: RecordKeyNames{
				Time:      timeKey,
				Msg:       msgKey,
				Level:     levelKey,
				Call:      locationKey,
				RequestID: requestID,
			},
		})
	}
}

func (l *logger) New(keyValues ...interface{}) Logger {
	child := &logger{newKeyValues(l.KeyValues, keyValues), new(swapHandler), LvlTrace}
	child.SetHandler(l.handler)
	return child
}

func newKeyValues(prefix []interface{}, suffix []interface{}) []interface{} {
	normalizedSuffix := normalize(suffix)
	newCtx := make([]interface{}, len(prefix)+len(normalizedSuffix))
	n := copy(newCtx, prefix)
	copy(newCtx[n:], normalizedSuffix)
	return newCtx
}

func (l *logger) SetOutLevel(level Level) {
	if level >= LvlFatal && level <= LvlTrace {
		l.level = level
	}
}

func (l *logger) GetOutLevel() Level {
	return l.level
}

func (l *logger) Log(msg string, fields ...interface{}) {
	l.write(msg, LvlTrace, fields)
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

func (l *logger) LogContext(ctx context.Context, msg string, fields ...interface{}) {
	l.writeContext(ctx, msg, LvlTrace, fields)
}

func (l *logger) DebugContext(ctx context.Context, msg string, fields ...interface{}) {
	l.writeContext(ctx, msg, LvlDebug, fields)
}

func (l *logger) InfoContext(ctx context.Context, msg string, fields ...interface{}) {
	l.writeContext(ctx, msg, LvlInfo, fields)
}

func (l *logger) WarnContext(ctx context.Context, msg string, fields ...interface{}) {
	l.writeContext(ctx, msg, LvlWarn, fields)
}

func (l *logger) ErrorContext(ctx context.Context, msg string, fields ...interface{}) {
	l.writeContext(ctx, msg, LvlError, fields)
}

func (l *logger) FatalContext(ctx context.Context, msg string, fields ...interface{}) {
	l.writeContext(ctx, msg, LvlFatal, fields)
	os.Exit(1)
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
