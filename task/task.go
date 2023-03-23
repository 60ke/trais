package main

import (
	"bytes"
	"encoding/json"
	"sync"
	"time"

	"github.com/60ke/trais/conf"
	"github.com/60ke/trais/db"
	"github.com/60ke/trais/downloader"
	"github.com/60ke/trais/log"
	"github.com/60ke/trais/web3"
	"github.com/robfig/cron/v3"
)

func StartTask() {
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
