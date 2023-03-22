package main

import (
	"fmt"

	"github.com/60ke/trais/conf"
	"github.com/60ke/trais/db"
	"github.com/60ke/trais/downloader"
	"github.com/60ke/trais/log"
)

func main() {
	log.LogInit("warn", "./test.log")
	log.Logger.Info("test1")
	log.LogReload("info", "./test.log")
	log.Logger.Info("test2")
	conf.ConfInit("/Users/k/dev/8lab/trais/conf/conf.ini")
	fmt.Println(conf.APPSetting, conf.TaskSetting, conf.ServerSetting)
	fmt.Println(downloader.GetBestRpc(*conf.BscHosts))

	db.BscDBInit(conf.DatabaseSetting.User, conf.DatabaseSetting.Password, conf.DatabaseSetting.Host, conf.DatabaseSetting.Bsc)

	downloader.SyncBsc(*conf.BscHosts)
}
