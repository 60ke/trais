package web3

import "github.com/60ke/trais/db"

type LatestBlockNumberResp struct {
	Jsonrpc string `json:"jsonrpc"`
	ID      int    `json:"id"`
	Result  string `json:"result"`
}

// eth_getBlockByNumber with true
type BscBlockResp struct {
	Jsonrpc string `json:"jsonrpc"`
	ID      int    `json:"id"`
	Result  struct {
		db.BscBlock
		// 数据库未存此字段
		MixHash      string              `json:"mixHash"`
		Transactions []db.BscTransaction `json:"transactions"`
		// CreditData,CreditValue,CreditMax
		TrustNodeScore string `json:"trustNodeScore"`
		// TODO 暂未发现json uncles具体数据
		Uncles []interface{} `json:"uncles"`
	} `json:"result"`
}
