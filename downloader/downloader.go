package downloader

import (
	"bytes"
	"encoding/json"
	"fmt"

	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/60ke/trais/conf"
	"github.com/60ke/trais/log"

	"github.com/60ke/trais/db"

	"github.com/60ke/trais/web3"
	lru "github.com/hnlq715/golang-lru"
	"gorm.io/gorm"
)

var (
	// 获取失败的区块
	FailedBlock, _ = lru.NewARCWithExpire(1000, 0)
	// 最近60s内的交易的账户地址
	LatestAddr, _ = lru.NewARCWithExpire(1000, 10)
)

type GetBlockResult struct {
	// 成功入库
	Succ bool
	// 区块号
	BlockNum int64
}

func SyncBsc(hosts []string) error {
	var wg sync.WaitGroup
	var bscStep = conf.DownloaderSetting.BscStep
	var resp web3.LatestBlockNumberResp

	bestRpc := GetBestRpc(hosts)
	body, err := web3.LatestBlockNumber(bestRpc)
	if err != nil {
		log.Logger.Error(err)
		return err
	}

	decoder := json.NewDecoder(bytes.NewReader(body))
	decoder.DisallowUnknownFields()
	err = decoder.Decode(&resp)
	if err != nil {
		log.Logger.Error(string(body))
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

	// 限制追块最大数量
	if endNum-startNum > bscStep {
		endNum = startNum + bscStep
	}

	// 开始追块
	// results := make(chan GetBlockResult, endNum-startNum+1)

	for i := startNum; i <= endNum; i++ {
		wg.Add(1)
		go func(i int64) {
			defer wg.Done()
			GetBscBlock(bestRpc, i)
		}(i)
	}
	wg.Wait()

	if latestNum-endNum < bscStep {
		// 落后较少时,可以sleep一下减少资源占用
		time.Sleep(1 * time.Second)
	}
	return nil
}

func GetBscBlock(rpc string, num int64) {
	// TODO 解析区块并存入数据库
	// 使用LRU/缓存更新账户地址余额

	var resp web3.BscRpcBlock
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
			FailedBlock.Add(num, nil)
			getBlockResult.Succ = false
		} else {
			insertBscBlock(resp)
			getBlockResult.Succ = true
		}
	}

}

func insertBscBlock(resp web3.BscRpcBlock) {
	var tx db.BscTransactionTable
	block, txs := BscRpcBlock2Table(resp)
	db.BscDB.Table(block.TableName()).Create(&block)
	db.BscDB.Table(tx.TableName()).Create(&txs)
}

func BscRpcBlock2Table(resp web3.BscRpcBlock) (db.BscBlockTable, []db.BscTransactionTable) {
	var block db.BscBlockTable
	// var txs []db.BscTransactionTable
	block.BscBlockCommon = resp.Result.BscBlockCommon
	block.TransactionsCount = int64(len(resp.Result.Transactions))
	block.UncleCount = int64(len(resp.Result.Uncles))
	block.CreditData = resp.Result.TrustNodeScore
	block.CreditValue, block.CreditMax = ParseTrustNodeScore(block.CreditData, resp.Result.Miner)
	block.GasLimit = Hex2int64(resp.Result.GasLimit)
	block.GasUsed = Hex2int64(resp.Result.GasUsed)
	block.Timestamp = Hex2int64(resp.Result.Timestamp)
	log.Logger.Info(resp.Result.Number)
	block.Number = Hex2int64(resp.Result.Number)
	txs := BscRpcTx2Table(resp.Result.Transactions)
	// TODO-----
	return block, txs
}

func BscRpcTx2Table(txs []web3.BscRpcTransaction) []db.BscTransactionTable {
	var txTables []db.BscTransactionTable
	for _, tx := range txs {
		var txTable db.BscTransactionTable
		txTable.BscTransactionCommon = tx.BscTransactionCommon
		txTable.Gas = Hex2int64(tx.Gas)
		txTable.GasPrice = Hex2int64(tx.GasPrice)
		txTable.GasUsed = Hex2int64(tx.GasUsed)
		txTable.Timestamp = Hex2int64(tx.Timestamp)
		txTable.BlockNumber = Hex2int64(tx.BlockNumber)
		txTable.TransactionIndex = Hex2int64(tx.TransactionIndex)
		txTable.Nonce = Hex2int64(tx.Nonce)
		txTable.V = Hex2int64(tx.V)
		txTable.Type = int(Hex2int64(tx.Type))
		LatestAddr.Add(tx.From, nil)
		LatestAddr.Add(tx.To, nil)

		txTables = append(txTables, txTable)
	}
	return txTables
}

func ParseTrustNodeScore(trustNodeScore, miner string) (string, int64) {
	trustNodeScore = strings.TrimPrefix(trustNodeScore, "0x")
	miner = strings.TrimPrefix(miner, "0x")
	trustLen := len(trustNodeScore) / 52
	var creditMax int64
	var creditValue string = "["
	for i := 0; i < trustLen; i++ {
		trustData := trustNodeScore[i*52 : (i+1)*52]
		addr := trustData[:40]
		score, _ := strconv.ParseInt(trustData[44:], 16, 64)
		if strings.EqualFold(addr, miner) {
			creditMax = score
		}
		data := fmt.Sprintf("{\"%s\": %d},", addr, score)
		creditValue += data
	}
	creditValue = creditValue[:len(creditValue)-1] + "]"
	return creditValue, creditMax
}

func Hex2int64(hexStr string) int64 {
	var num int64
	hexStr = strings.TrimPrefix(hexStr, "0x")
	if hexStr != "" {
		var err error
		num, err = strconv.ParseInt(hexStr, 16, 64)
		if err != nil {
			log.Logger.Error(err)
		}
	}

	return num
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
		body, err := web3.LatestBlockNumber(rpc)
		if err != nil {
			log.Logger.Error(err)
			continue
		}

		var resp web3.LatestBlockNumberResp
		decoder := json.NewDecoder(bytes.NewReader(body))
		decoder.DisallowUnknownFields()
		err = decoder.Decode(&resp)
		if err != nil {
			log.Logger.Error(err)
			continue
		}
		elapsed = time.Since(start)

		if latency == 0*time.Second || latency > elapsed {
			latency = elapsed
			bestRpc = rpc
		}
	}
	log.Logger.Infof("bestRpc is %s,latency: %s", bestRpc, latency)
	return bestRpc
}
