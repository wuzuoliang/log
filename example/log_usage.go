package main

import (
	"github.com/wuzuoliang/log"
)

func main(){
	newLog:=log.New()
	newLog.Info("sss","22",log.JSON(&struct{
		A int
		B string
	}{1,"2"}))
	return

	log.SetRotatePara(100, 10, 30, true, true)
	//
	h, _ := log.FileHandlerRotate("test.log", log.LogfmtFormat())
	log.Root().SetHandler(h)
	logLevel:="debug"
	switch logLevel {
	case "debug":
		log.SetOutLevel(log.LvlDebug)
	case "info":
		log.SetOutLevel(log.LvlInfo)
	case "warn":
		log.SetOutLevel(log.LvlWarn)
	default:
		log.SetOutLevel(log.LvlError)
	}

	log.Debug("a","1",2,"3",4)
}

