package main

import (
	"fmt"

	"github.com/60ke/trais/conf"
	"github.com/60ke/trais/log"
)

func main() {
	log.LogInit("warn", "./test.log")
	log.Logger.Info("test1")
	log.LogReload("info", "./test.log")
	log.Logger.Info("test2")
	conf.ConfInit()
	fmt.Println(conf.APPSetting, conf.TaskSetting, conf.ServerSetting)
}
