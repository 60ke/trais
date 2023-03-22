package downloader

import (
	"bytes"
	"encoding/json"
	"fmt"

	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/60ke/trais/log"

	"github.com/60ke/trais/db"

	"github.com/60ke/trais/web3"
	"gorm.io/gorm"
)

type GetBlockResult struct {
	// 成功入库
	Succ bool
	// 区块号
	BlockNum int64
}

func SyncBsc(hosts []string) error {
	var wg sync.WaitGroup
	bestRpc := GetBestRpc(hosts)
	body, err := web3.LatestBlockNumber(bestRpc)
	if err != nil {
		log.Logger.Error(err)
		return err
	}
	var resp web3.LatestBlockNumberResp
	decoder := json.NewDecoder(bytes.NewReader(body))
	decoder.DisallowUnknownFields()
	err = decoder.Decode(&resp)
	if err != nil {
		log.Logger.Error(err)
		return err
	}
	startNum := getBscStartNum()
	resp.Result = strings.TrimPrefix(resp.Result, "0x")
	latestNum, err := strconv.ParseInt(resp.Result, 16, 64)
	endNum := latestNum
	if err != nil {
		log.Logger.Error(err)
		return err
	}
	if startNum > latestNum {
		err := fmt.Errorf("数据库区块号:%d,大于当前链的最新区块:%d", startNum, latestNum)
		log.Logger.Error(err)
		return err
	}

	// 限制追块最大数量为10
	if endNum-startNum > 10 {
		endNum = startNum + 10
	}

	// 开始追块
	results := make(chan GetBlockResult, endNum-startNum+1)

	for i := startNum; i <= endNum; i++ {
		wg.Add(1)
		go func(i int64, results chan GetBlockResult) {
			defer wg.Done()
			getBscBlock(bestRpc, i, results)
		}(i, results)
	}
	wg.Wait()
	close(results)

	// 处理失败的块号
	for result := range results {
		if !result.Succ {
			log.Logger.Warnf("get block failed: %d", result.BlockNum)
			// TODO addFailBlock,removeFail
			// addFailBlock(result.BlockNum)
		}
	}
	if latestNum-endNum < 10 {
		// 落后较少时,可以sleep一下减少资源占用
		time.Sleep(1 * time.Second)
	}
	return nil
}

func getBscBlock(rpc string, num int64, results chan GetBlockResult) {
	// TODO 解析区块并存入数据库
	// 使用LRU/缓存更新账户地址余额

	var resp web3.BscBlockResp
	var getBlockResult GetBlockResult
	getBlockResult.BlockNum = num

	data, err := web3.GetBlockByNum(rpc, fmt.Sprintf("0x%x", num))
	if err != nil {
		log.Logger.Error(err)
		getBlockResult.Succ = false
	} else {
		decoder := json.NewDecoder(bytes.NewReader(data))
		decoder.DisallowUnknownFields()
		err = decoder.Decode(&resp)
		if err != nil {
			log.Logger.Error(err)
			getBlockResult.Succ = false
		} else {
			insertBscBlock(resp)
			getBlockResult.Succ = true
		}
	}

	results <- getBlockResult

}

func insertBscBlock(resp web3.BscBlockResp) {
	var block db.BscBlockTable
	block.BscBlock = resp.Result.BscBlock
	block.TransactionsCount = int64(len(resp.Result.Transactions))
	block.UncleCount = int64(len(resp.Result.Uncles))
	db.BscDB.Table(block.TableName()).Create(block)
	// TODO parse trustScore,insert bsc txs
}

func getBscStartNum() int64 {
	// var startNum int64
	var bscBlock db.BscBlockTable
	// 从数据最大区块开始,防止上个区块未完全获取
	return getLatestBlockFromDb(db.BscDB, bscBlock.TableName())
}

func getLatestBlockFromDb(sql *gorm.DB, table string) int64 {
	var block db.BscBlockTable
	// sql.Table(table).Select("max(number)").Scan(&block)
	// row := db.BscDB.Table(table).Select("max(number)").Row()
	// sql.Table(table).First(&block)
	sql.Table(table).Order("number desc").First(&block)
	return block.Number
}

// 获取延迟最低的rpc接口
func GetBestRpc(hosts []string) string {
	var bestRpc string
	var latency time.Duration
	for _, host := range hosts {
		var elapsed time.Duration
		rpc := fmt.Sprintf("http://%s:8545", host)
		start := time.Now()
		_, err := web3.LatestBlockNumber(rpc)
		if err == nil {
			elapsed = time.Since(start)
		} else {
			log.Logger.Error(err)
		}

		if latency == 0*time.Second || latency > elapsed {
			latency = elapsed
			bestRpc = rpc
		}
	}
	log.Logger.Infof("bestRpc is %s,latency: %s", bestRpc, latency)
	return bestRpc
}
