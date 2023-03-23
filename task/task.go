package task

import (
	"bytes"
	"encoding/json"
	"sync"
	"time"

	"github.com/60ke/trais/conf"
	"github.com/60ke/trais/db"
	"github.com/60ke/trais/downloader"
	"github.com/60ke/trais/log"
	"github.com/60ke/trais/tools"
	"github.com/60ke/trais/web3"
	"github.com/robfig/cron/v3"
)

func StartTask() {
	log.Logger.Info("StartTask")
	c := cron.New()

	updateBscBalance := conf.TaskSetting.UpdateBscBalance
	handleFailBlock := conf.TaskSetting.HandleFailBlock
	hosts := conf.BscHosts
	c.AddFunc(updateBscBalance, func() {
		UpdateBscAddress(*hosts)
	})
	c.AddFunc(handleFailBlock, func() {
		HandleFailBlock(*hosts)
	})
	c.Start()
}

// 更新bsc账户余额
func UpdateBscAddress(hosts []string) {

	bestRpc := downloader.GetBestRpc(hosts)

	keys := downloader.LatestAddr.Keys()

	for _, key := range keys {
		var resp web3.BscRpcBalance
		if _, ok := downloader.LatestAddr.Get(key); ok {

			addr := key.(string)
			body, err := web3.GetBalance(bestRpc, addr)
			if err != nil {
				log.Logger.Error(err)
				continue
			}

			// 校验返回值
			decoder := json.NewDecoder(bytes.NewReader(body))
			decoder.DisallowUnknownFields()
			err = decoder.Decode(&resp)
			if err != nil {
				log.Logger.Error(err)
				continue
			}

			// 插入数据库
			var address db.BscAddress
			balance := downloader.Hex2int64(resp.Result)
			timestamp := time.Now().Unix()
			address.Balance = balance
			address.Time = timestamp
			address.Address = addr
			db.BscDB.Table(address.TableName()).Where(db.BscAddress{Address: addr}).Assign(db.BscAddress{Time: timestamp, Balance: balance}).FirstOrCreate(&address)
			downloader.LatestAddr.Remove(addr)
		}

	}
}

func HandleFailBlock(hosts []string) {

	var wg sync.WaitGroup

	keys := downloader.FailedBlock.Keys()

	bestRpc := downloader.GetBestRpc(hosts)
	for _, key := range keys {
		num := key.(int64)
		wg.Add(1)
		go func(num int64) {
			defer wg.Done()
			downloader.GetBscBlock(bestRpc, num)
		}(num)
	}

	wg.Wait()
}

/*
当前程序使用并发,可能由于断电等异常造成程序退出
重新处理最近一次的1000(步进)块
*/
func HandleLastTask(hosts []string) {
	var nums []int64
	var wg sync.WaitGroup
	bestRpc := downloader.GetBestRpc(hosts)

	db.BscDB.Table("block").Model(db.BscBlockTable{}).Order("number desc").Limit(int(conf.DownloaderSetting.BscStep)).Select("Number").Find(&nums)
	maxNum, minNum := nums[len(nums)-1], nums[0]
	for i := int64(0); maxNum < minNum; i++ {
		block := nums[len(nums)-1] + i
		if !tools.SliceContain(nums, block) {
			// 获取数据库中未爬取的块
			nums = append(nums, block)
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

}
