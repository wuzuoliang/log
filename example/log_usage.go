package main

import (
	"context"
	"github.com/wuzuoliang/log"
	"os"
)

func main() {
	newLog := log.New()
	newLog.Info("newLog", "test1", log.JSON(&struct {
		A int
		B string
	}{1, "2"}))

	newLog.Log("newLog", "test2", &struct {
		A int
		B string
	}{3, "4"})

	newLog2 := log.New("withinit", "test2")
	newLog2.Error("sss", "1", "2")
	newLog2.Warn("1111", "fuck", "you")

	newLog3 := log.New()
	newLog3.SetHandler(log.CallerFileHandler(log.StreamHandler(os.Stdout, log.LogfmtFormat())))
	newLog3.Debug("1111", "asd", "aaaa")
	newLog3.SetHandler(log.CallerFuncHandler(log.StreamHandler(os.Stdout, log.LogfmtFormat())))
	newLog3.Debug("2222", "asd", "aaaa")
	newLog3.SetHandler(log.CallerFuncHandler(log.StreamHandler(os.Stdout, log.TerminalFormat())))
	newLog3.Debug("3333", "asd", "aaaa")
	newLog3.SetHandler(log.CallerFileHandler(log.StreamHandler(os.Stdout, log.JsonFormat())))
	newLog3.Debug("4444", "asd", "aaaa")

	log.Root().SetHandler(log.StderrHandler)
	logLevel := "111"
	switch logLevel {
	case "debug":
		log.SetOutLevel(log.LvlDebug)
	case "info":
		log.SetOutLevel(log.LvlInfo)
	case "warn":
		log.SetOutLevel(log.LvlWarn)
	default:
		log.SetOutLevel(log.LvlTrace)
	}
	log.Log("asdasdas")
	log.Debug("a", "1", 2, "3", 4, "5", 6, "7", 8, "9", 10)

	//nCtx:=context.WithValue(context.Background(),"trace","9527")
	//log.FatalContext(nCtx,"12","11","22")

	log.DebugContext(context.Background(), "a", "1", 2, "3", 4, "5", 6, "7", 8, "9", 10)
	ctx := context.WithValue(context.Background(), "request_id", "asdasdkjcxizcjci")
	log.LogContext(ctx, "1", "123", "123123")
	log.FatalContext(ctx, "asdas", "213", "12312", "1111")
}
