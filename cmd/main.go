package main

import (
	"os"
	"sync"

	"github.com/60ke/trais/conf"
	"github.com/60ke/trais/db"
	"github.com/60ke/trais/downloader"
	"github.com/60ke/trais/log"
	"github.com/60ke/trais/task"
	"github.com/60ke/trais/tools"
)

func Usage() {

}

var once sync.Once

func onceDo() {
	/*
		程序使用并发,可能由于断电等异常造成程序退出
		重新处理最近一次的1000(步进)块
	*/
	once.Do(func() {
		var nums []int64
		var wg sync.WaitGroup
		bestRpc := downloader.GetBestRpc(*conf.BscHosts)

		db.BscDB.Table("block").Model(db.BscBlockTable{}).Order("number desc").Limit(int(conf.DownloaderSetting.BscStep)).Select("Number").Find(&nums)
		maxNum, minNum := nums[len(nums)-1], nums[0]
		for i := minNum; i < maxNum; i++ {
			if !tools.SliceContain(nums, i) {
				// 获取数据库中未爬取的块
				nums = append(nums, i)
			}
		}

		for _, num := range nums {
			wg.Add(1)
			go func(num int64) {
				defer wg.Done()
				downloader.GetBscBlock(bestRpc, num)
			}(num)
		}

		wg.Wait()
		log.Logger.Info("HandleLastTask Done")
	})

}

func main() {
	// log.LogInit("warn", "./test.log")
	// log.Logger.Info("test1")
	// log.LogReload("info", "./test.log")
	// log.Logger.Info("test2")
	if len(os.Args) == 1 {
		conf.ConfInit("./conf.ini")
	} else {
		conf.ConfInit(os.Args[1])
	}

	log.LogInit(conf.LogSetting.Level, conf.LogSetting.File)

	db.BscDBInit(conf.DatabaseSetting.User, conf.DatabaseSetting.Password, conf.DatabaseSetting.Host, conf.DatabaseSetting.Bsc)

	onceDo()
	go task.StartTask()
	downloader.SyncBsc(*conf.BscHosts)

}
