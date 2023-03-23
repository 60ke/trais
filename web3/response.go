package web3

import "github.com/60ke/trais/db"

type LatestBlockNumberResp struct {
	Jsonrpc string `json:"jsonrpc"`
	ID      int    `json:"id"`
	Result  string `json:"result"`
}

type BscRpcBalance struct {
	Jsonrpc string `json:"jsonrpc"`
	ID      int    `json:"id"`
	Result  string `json:"result"`
}

type BscRpcTransaction struct {
	db.BscTransactionCommon
	Gas              string `gorm:"column:gas" json:"gas"`
	GasPrice         string `gorm:"column:gasPrice" json:"gasPrice"`
	Nonce            string `gorm:"column:nonce" json:"nonce"`
	V                string `gorm:"column:v" json:"v"`
	GasUsed          string `gorm:"column:gasUsed" json:"gasUsed"`
	Timestamp        string `gorm:"column:timestamp" json:"timestamp"`
	BlockNumber      string `gorm:"column:blockNumber;NOT NULL" json:"blockNumber"`
	TransactionIndex string `gorm:"column:transactionIndex" json:"transactionIndex"`
	Type             string `gorm:"column:type" json:"type"`
}

// eth_getBlockByNumber with true
type BscRpcBlock struct {
	Jsonrpc string `json:"jsonrpc"`
	ID      int    `json:"id"`
	Result  struct {
		db.BscBlockCommon

		Size      string `gorm:"column:size" json:"size"`
		GasLimit  string `gorm:"column:gasLimit" json:"gasLimit"`
		GasUsed   string `gorm:"column:gasUsed" json:"gasUsed"`
		Timestamp string `gorm:"column:timestamp" json:"timestamp"`
		Number    string `gorm:"column:number;NOT NULL" json:"number"`

		// 数据库未存此字段
		MixHash      string              `json:"mixHash"`
		Transactions []BscRpcTransaction `json:"transactions"`
		// CreditData,CreditValue,CreditMax
		TrustNodeScore string `json:"trustNodeScore"`
		// TODO 暂未发现json uncles具体数据
		Uncles []string `json:"uncles"`
	} `json:"result"`
}
