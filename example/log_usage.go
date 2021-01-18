package main

import (
	"github.com/wuzuoliang/log"
)

func main(){
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