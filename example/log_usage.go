package main

import (
	"github.com/mattn/go-colorable"
	"github.com/wuzuoliang/log"
)

func main(){
	//newLog:=log.New()
	//newLog.Info("sss","22",log.JSON(&struct{
	//	A int
	//	B string
	//}{1,"2"}))
	//return
	//newLog.Error("ss","23",333,"123","1111")
	//newLog.Debug("de","44",12312312)
	//newLog.Warn("warn","wwww","asdasdasdqwejqwiojfoiajoifushiuheqihiquwhouqwh")
	//newLog.Fatal("ffff","asd",123,"!@3123",111,"2",struct{
	//	A int
	//	B string
	//}{3,"4"},"5")
	//return
	//log.SetRotatePara(100, 10, 30, true, true)
	//
	//log.SetRotatePara(100, 10, 30, true, true)
	//h, _ := log.FileHandlerRotate("test.log", log.JsonFormat())
	//h:=log.StreamHandler(colorable.NewColorableStdout(), log.JsonFormat())
	h:=log.StreamHandler(colorable.NewColorableStdout(), log.TerminalFormat())

	log.Root().SetHandler(h)
	logLevel:="111"
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
	log.Debug("a","1",2,"3",4,"5",6,"7",8,"9",10)


	//nCtx:=context.WithValue(context.Background(),"trace","9527")
	//log.FatalContext(nCtx,"12","11","22")
}

