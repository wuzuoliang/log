package main

import (
	"context"
	"github.com/wuzuoliang/log"
)

func main(){
	newLog:=log.New()
	newLog.Info("sss","22",log.JSON(&struct{
		A int
		B string
	}{1,"2"}))
	//return
	newLog.Error("ss","23",333,"123","1111")
	newLog.Debug("de","44",12312312)
	newLog.Warn("warn","wwww","asdasdasdqwejqwiojfoiajoifushiuheqihiquwhouqwh")
	//newLog.Fatal("ffff","asd",123,"!@3123",111,"2",struct{
	//	A int
	//	B string
	//}{3,"4"},"5")
	//return
	log.SetRotatePara(100, 10, 30, true, true)
	//
	h, _ := log.FileHandlerRotate("test.log", log.JsonFormat())

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

	log.Debug("a","1",2,"3",4,"5",6,7,8,9,10)


	nCtx:=context.WithValue(context.Background(),"trace","9527")
	log.FatalContext(nCtx,"12","11","22")
}

