package log

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"sync"
	"time"
)

const (
	timeFormat     = "2006-01-02 15:04:05.999"
	termTimeFormat = "15:04:05.999"
	floatFormat    = 'f'
	termMsgJust    = 40
)

type Format interface {
	Format(r *Record) []byte
}

// FormatFunc returns a new Format object which uses
// the given function to perform record formatting.
func FormatFunc(f func(*Record) []byte) Format {
	return formatFunc(f)
}

type formatFunc func(*Record) []byte

func (f formatFunc) Format(r *Record) []byte {
	return f(r)
}

// TerminalFormat formats log records optimized for human readability on
// a terminal with color-coded level output and terser human friendly timestamp.
// This format should only be used for interactive programs or while developing.
//
//     [TIME] [LEVEL] MESAGE key=value key=value ...
//
// Example:
//
//     [May 16 20:58:45] [DBUG] remove route ns=haproxy addr=127.0.0.1:50002
//
func TerminalFormat() Format {
	return FormatFunc(func(r *Record) []byte {
		var color = 0
		switch r.Level {
		case LvlFatal:
			color = 35
		case LvlError:
			color = 31
		case LvlWarn:
			color = 33
		case LvlInfo:
			color = 32
		case LvlDebug:
			color = 36
		case LvlTrace:
			color = 30
		}

		b := &bytes.Buffer{}
		lvl := strings.ToUpper(r.Level.String())
		if color > 0 {
			fmt.Fprintf(b, "\x1b[%dm%s\x1b[0m [%s][%s] %s=%s ", color, lvl, r.Time.Format(termTimeFormat), r.Call.String(), r.KeyNames.Msg, r.Msg)
		} else {
			fmt.Fprintf(b, "[%s][%s][%s] %s=%s ", lvl, r.Call.String(), r.Time.Format(termTimeFormat), r.KeyNames.Msg, r.Msg)
		}

		if r.Ctx != nil && r.Ctx.Value(requestID) != nil {
			requestID := r.Ctx.Value(requestID).(string)
			if len(requestID) > 0 {
				fmt.Fprintf(b, "%s=%s ", r.KeyNames.RequestID, requestID)
			}
		}

		// try to justify the log output for short messages
		// 此处控制的是msg和后续kvalues中间的' '
		//if len(r.KeyValues) > 0 && len(r.Msg) < termMsgJust {
		//	b.Write(bytes.Repeat([]byte{' '}, termMsgJust-len(r.Msg)))
		//}

		// print the keys kvaluesfmt style
		kvaluesfmt(b, r.KeyValues, color)
		return b.Bytes()
	})
}

// LogfmtFormat prints records in kvaluesfmt format, an easy machine-parseable but human-readable
// format for key/value pairs.
//
// For more details see: https://pkg.go.dev/github.com/kr/logfmt
//
func LogfmtFormat() Format {
	return FormatFunc(func(r *Record) []byte {
		var caller string
		if r.CustomCaller == "" {
			caller = r.Call.String()
		} else {
			caller = r.CustomCaller
		}

		common := []interface{}{r.KeyNames.Time, r.Time, r.KeyNames.Level, r.Level, r.KeyNames.Call, caller, r.KeyNames.Msg, r.Msg}

		if r.Ctx != nil && r.Ctx.Value(requestID) != nil {

			requestID := r.Ctx.Value(requestID).(string)
			if len(requestID) > 0 {
				common = append(common, r.KeyNames.RequestID, requestID)
			}
		}
		buf := &bytes.Buffer{}
		kvaluesfmt(buf, append(common, r.KeyValues...), 0)
		return buf.Bytes()
	})
}

func kvaluesfmt(buf *bytes.Buffer, KeyValues []interface{}, color int) {
	for i := 0; i < len(KeyValues); i += 2 {
		if i != 0 {
			buf.WriteByte(' ')
		}

		k, ok := KeyValues[i].(string)
		v := formatLogfmtValue(KeyValues[i+1])
		if !ok {
			//k, v = errorKey, formatLogfmtValue(k)
			k, v = errorKey, fmt.Sprintf("%+v is not a string key", KeyValues[i+1])
		}

		// XXX: we should probably check that all of your key bytes aren't invalid
		if color > 0 {
			fmt.Fprintf(buf, "\x1b[%dm%s\x1b[0m=%s", color, k, v)
		} else {
			if i < 5 {
				buf.WriteByte('[')
			} else {
				buf.WriteString(k)
				buf.WriteByte('=')
			}
			buf.WriteString(v)
			if i < 5 {
				buf.WriteByte(']')
			}
		}
	}

	buf.WriteByte('\n')
}

// JsonFormat formats log records as JSON objects separated by newlines.
// It is the equivalent of JsonFormatEx(false, true).
func JsonFormat() Format {
	return JsonFormatEx(false, true)
}

// JsonFormatEx formats log records as JSON objects. If pretty is true,
// records will be pretty-printed. If lineSeparated is true, records
// will be logged with a new line between each record.
func JsonFormatEx(pretty, lineSeparated bool) Format {
	jsonMarshal := json.Marshal
	if pretty {
		jsonMarshal = func(v interface{}) ([]byte, error) {
			return json.MarshalIndent(v, "", "    ")
		}
	}

	return FormatFunc(func(r *Record) []byte {
		props := make(map[string]interface{})

		props[r.KeyNames.Time] = r.Time
		props[r.KeyNames.Level] = r.Level.String()
		props[r.KeyNames.Msg] = r.Msg

		if r.Ctx != nil && r.Ctx.Value(requestID) != nil {

			requestID := r.Ctx.Value(requestID).(string)
			if len(requestID) > 0 {
				props[r.KeyNames.RequestID] = requestID
			}
		}

		for i := 0; i < len(r.KeyValues); i += 2 {
			k, ok := r.KeyValues[i].(string)
			if !ok {
				props[errorKey] = fmt.Sprintf("%+v is not a string key", r.KeyValues[i])
			}
			props[k] = formatJsonValue(r.KeyValues[i+1])
		}

		b, err := jsonMarshal(props)
		if err != nil {
			b, _ = jsonMarshal(map[string]string{
				errorKey: err.Error(),
			})
			return b
		}

		if lineSeparated {
			b = append(b, '\n')
		}

		return b
	})
}

func formatShared(value interface{}) (result interface{}) {
	defer func() {
		if err := recover(); err != nil {
			if v := reflect.ValueOf(value); v.Kind() == reflect.Ptr && v.IsNil() {
				result = "nil"
			} else {
				panic(err)
			}
		}
	}()

	switch v := value.(type) {
	case time.Time:
		return v.Format(timeFormat)

	case error:
		return v.Error()

	case fmt.Stringer:
		return v.String()

	default:
		return v
	}
}

func formatJsonValue(value interface{}) interface{} {
	value = formatShared(value)
	switch value.(type) {
	case int, int8, int16, int32, int64, float32, float64, uint, uint8, uint16, uint32, uint64, string:
		return value
	default:
		return fmt.Sprintf("%+v", value)
	}
}

// formatValue formats a value for serialization
func formatLogfmtValue(value interface{}) string {
	if value == nil {
		return "nil"
	}

	if t, ok := value.(time.Time); ok {
		// Performance optimization: No need for escaping since the provided
		// timeFormat doesn't have any escape characters, and escaping is
		// expensive.
		return t.Format(timeFormat)
	}
	value = formatShared(value)
	switch v := value.(type) {
	case bool:
		return strconv.FormatBool(v)
	case float32:
		return strconv.FormatFloat(float64(v), floatFormat, 3, 64)
	case float64:
		return strconv.FormatFloat(v, floatFormat, 3, 64)
	case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64:
		return fmt.Sprintf("%d", value)
	case string:
		return escapeString(v)
	default:
		return escapeString(fmt.Sprintf("%+v", value))
	}
}

var stringBufPool = sync.Pool{
	New: func() interface{} { return new(bytes.Buffer) },
}

func escapeString(s string) string {
	needsQuotes := false
	needsEscape := false
	for _, r := range s {
		if r <= ' ' || r == '=' || r == '"' {
			needsQuotes = true
		}
		if r == '\\' || r == '"' || r == '\n' || r == '\r' || r == '\t' {
			needsEscape = true
		}
	}
	if !needsEscape && !needsQuotes {
		return s
	}
	e := stringBufPool.Get().(*bytes.Buffer)
	e.WriteByte('"')
	for _, r := range s {
		switch r {
		case '\\', '"':
			e.WriteByte('\\')
			e.WriteByte(byte(r))
		case '\n':
			e.WriteString("\\n")
		case '\r':
			e.WriteString("\\r")
		case '\t':
			e.WriteString("\\t")
		default:
			e.WriteRune(r)
		}
	}
	e.WriteByte('"')
	var ret string
	if needsQuotes {
		ret = e.String()
	} else {
		ret = string(e.Bytes()[1 : e.Len()-1])
	}
	e.Reset()
	stringBufPool.Put(e)
	return ret
}

// JSON is a helper function, following is its function code.
//
//  data, _ := json.Marshal(v)
//  return string(data)
func JSON(v interface{}) string {
	pool := getBytesBufferPool()
	buffer := pool.Get()
	defer pool.Put(buffer)
	buffer.Reset()

	if err := json.NewEncoder(buffer).Encode(v); err != nil {
		return ""
	}
	data := buffer.Bytes()

	// remove the trailing newline
	i := len(data) - 1
	if i < 0 || i >= len(data) /* BCE */ {
		return ""
	}
	if data[i] == '\n' {
		data = data[:i]
	}
	return string(data)
}

// XML is a helper function, following is its function code.
//
//  data, _ := xml.Marshal(v)
//  return string(data)
func XML(v interface{}) string {
	pool := getBytesBufferPool()
	buffer := pool.Get()
	defer pool.Put(buffer)
	buffer.Reset()

	if err := xml.NewEncoder(buffer).Encode(v); err != nil {
		return ""
	}
	return string(buffer.Bytes())
}
