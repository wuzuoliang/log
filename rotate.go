package log

import (
	"gopkg.in/natefinch/lumberjack.v2"
	"sync/atomic"
	"unsafe"
)

// rotateOptions storage rotate file parameters
type rotateOptions struct {
	maxSize       int                // 单个文件大小
	maxSaveDay    int                // 文件最多存储天数
	maxBackup     int                // 最多备份数量
	doesCompress  bool               // 切割后文件是否压缩
	doesDayRotate bool               // 是否日切
	output        *lumberjack.Logger // 文件指针
}

var _defaultRotateOption = &rotateOptions{100, 7, 50, false, true, nil}

func GetDefaultRotateOption() *rotateOptions {
	return _defaultRotateOption
}

type RotateOptions func(*rotateOptions)

func SetMaxSize(max int) RotateOptions {
	return func(o *rotateOptions) {
		o.SetMaxSize(max)
	}
}

func (opts *rotateOptions) SetMaxSize(max int) {
	opts.maxSize = max
}

func SetMaxSaveDay(day int) RotateOptions {
	return func(o *rotateOptions) {
		o.SetMaxSaveDay(day)
	}
}
func (opts *rotateOptions) SetMaxSaveDay(day int) {
	opts.maxSaveDay = day
}

func SetMaxBackup(backup int) RotateOptions {
	return func(o *rotateOptions) {
		o.SetMaxBackup(backup)
	}
}
func (opts *rotateOptions) SetMaxBackup(backup int) {
	opts.maxBackup = backup
}

func SetCompress(doesCompress bool) RotateOptions {
	return func(o *rotateOptions) {
		o.SetCompress(doesCompress)
	}
}
func (opts *rotateOptions) SetCompress(doesCompress bool) {
	opts.doesCompress = doesCompress
}
func SetDayRotate(dayrorate bool) RotateOptions {
	return func(o *rotateOptions) {
		o.SetDayRotate(dayrorate)
	}
}

func (opts *rotateOptions) SetDayRotate(dayRotate bool) {
	opts.doesDayRotate = dayRotate
}

func SetOutput(output *lumberjack.Logger) RotateOptions {
	return func(o *rotateOptions) {
		o.SetOutput(output)
	}
}

func (opts *rotateOptions) SetOutput(output *lumberjack.Logger) {
	opts.output = output
}

func newRotateOptions(opts []RotateOptions) *rotateOptions {
	var o rotateOptions
	for _, opt := range getDefaultRotateOptions() {
		if opt == nil {
			continue
		}
		opt(&o)
	}
	for _, opt := range opts {
		if opt == nil {
			continue
		}
		opt(&o)
	}
	return &o
}

var _defaultOptionsPtr unsafe.Pointer // *[]RotateOption

func SetDefaultRotateOptions(opts []RotateOptions) {
	if opts == nil {
		atomic.StorePointer(&_defaultOptionsPtr, nil)
		return
	}
	atomic.StorePointer(&_defaultOptionsPtr, unsafe.Pointer(&opts))
}

func getDefaultRotateOptions() []RotateOptions {
	ptr := (*[]RotateOptions)(atomic.LoadPointer(&_defaultOptionsPtr))
	if ptr == nil {
		return nil
	}
	return *ptr
}
