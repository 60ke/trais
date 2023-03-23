package db

// TODO 原有的python项目有关数值在mysql中的存储均为bigint,故golang使用了int64
// 目前不存在数值溢出既保留的对原有数据的兼容也简单了计算,但是可能存在隐患
type Info struct {
	ID                  int64  `gorm:"column:id;primary_key;AUTO_INCREMENT" json:"id"`
	HashRate            int64  `gorm:"column:hashRate;NOT NULL" json:"hashRate"`
	TotalDifficulty     string `gorm:"column:totalDifficulty;NOT NULL" json:"totalDifficulty"`
	LastBlockFees       int64  `gorm:"column:lastBlockFees;NOT NULL" json:"lastBlockFees"`
	LastBlock           int64  `gorm:"column:lastBlock;NOT NULL" json:"lastBlock"`
	Addresses           int64  `gorm:"column:addresses;NOT NULL" json:"addresses"`
	Transactions        int64  `gorm:"column:transactions;NOT NULL" json:"transactions"`
	LastTransactionFees int64  `gorm:"column:lastTransactionFees;NOT NULL" json:"lastTransactionFees"`
	Unconfirmed         int64  `gorm:"column:unconfirmed;NOT NULL" json:"unconfirmed"`
	Timestamp           int64  `gorm:"column:timestamp;NOT NULL" json:"timestamp"`
}

func (m *Info) TableName() string {
	return "info"
}

type BscAddress struct {
	ID      int64  `gorm:"column:id;primary_key;AUTO_INCREMENT" json:"id"`
	Address string `gorm:"column:address;NOT NULL" json:"address"`
	Balance int64  `gorm:"column:balance;NOT NULL" json:"balance"`
	Time    int64  `gorm:"column:time;NOT NULL" json:"time"`

	// 原有项目未存下面三行数据,当前亦未存储
	Received string `gorm:"column:received" json:"received"`
	Sent     string `gorm:"column:sent" json:"sent"`
	TxCount  int64  `gorm:"column:txCount" json:"txCount"`
}

func (m *BscAddress) TableName() string {
	return "address"
}

/*
{"jsonrpc":"2.0","id":1,"result":{"blockHash":"0x74a135356bebf58b5da121134f8b77ed2f380c63945f68f2a5991dd8fc097e1f","blockNumber":"0x5248b","from":"0xf554ba6e54ce652654980c2a82de36992ab230f9","gas":"0xd894","gasPrice":"0x3b9aca00","hash":"0x9fdf68abf9e9ff068005c7397e4ba8d952d351841eb79e5c68cace56104a113b","input":"0x","nonce":"0xe5fa","to":"0x76637edf5587b122de4661564ffb4b35d2590ed5","transactionIndex":"0x78","value":"0x38d7ea4c68000","type":"0x0","v":"0x19067d","r":"0xf7e495411557b2895e224b90053623d01179b411fe694573d9f9daf5f9540f88","s":"0x2c64149bd7b5c59e65e8182f56322f5138df3b5aa3191c559429b6815c0b4f82"}}
*/
type BscTransactionTable struct {
	ID int64 `gorm:"column:id;primary_key;AUTO_INCREMENT" json:"id"`
	BscTransactionCommon

	Gas              int64 `gorm:"column:gas" json:"gas"`
	GasPrice         int64 `gorm:"column:gasPrice" json:"gasPrice"`
	Nonce            int64 `gorm:"column:nonce" json:"nonce"`
	V                int64 `gorm:"column:v" json:"v"`
	GasUsed          int64 `gorm:"column:gasUsed" json:"gasUsed"`
	Timestamp        int64 `gorm:"column:timestamp" json:"timestamp"`
	BlockNumber      int64 `gorm:"column:blockNumber;NOT NULL" json:"blockNumber"`
	TransactionIndex int64 `gorm:"column:transactionIndex" json:"transactionIndex"`
	Type             int   `gorm:"column:type" json:"type"`
	// gas * gasPrice
	Fee int64 `gorm:"column:fee" json:"fee"`
}

type BscTransactionCommon struct {
	BlockHash string `gorm:"column:blockHash;NOT NULL" json:"blockHash"`
	// json字段为from,原有的python项目使用source,这里使用已有的数据库字段
	From  string `gorm:"column:source" json:"from"`
	To    string `gorm:"column:to" json:"to"`
	Hash  string `gorm:"column:hash" json:"hash"`
	Value string `gorm:"column:value" json:"value"`
	R     string `gorm:"column:r" json:"r"`
	S     string `gorm:"column:s" json:"s"`
	// json字段为input,原有的python项目使用tx_str,这里使用已有的数据库字段
	Input string `gorm:"column:tx_str" json:"input"`
}

func (m *BscTransactionTable) TableName() string {
	return "transaction"
}

// 与web3 rpc公用的字段
type BscBlockCommon struct {
	// Number           int64  `gorm:"column:number;NOT NULL" json:"number"`
	Hash             string `gorm:"column:hash;NOT NULL" json:"hash"`
	ParentHash       string `gorm:"column:parentHash;NOT NULL" json:"parentHash"`
	Nonce            string `gorm:"column:nonce" json:"nonce"`
	Sha3Uncles       string `gorm:"column:sha3Uncles" json:"sha3Uncles"`
	LogsBloom        string `gorm:"column:logsBloom" json:"logsBloom"`
	TransactionsRoot string `gorm:"column:transactionsRoot" json:"transactionsRoot"`
	StateRoot        string `gorm:"column:stateRoot" json:"stateRoot"`
	ReceiptsRoot     string `gorm:"column:receiptsRoot" json:"receiptsRoot"`
	Miner            string `gorm:"column:miner" json:"miner"`
	Difficulty       string `gorm:"column:difficulty" json:"difficulty"`
	TotalDifficulty  string `gorm:"column:totalDifficulty" json:"totalDifficulty"`
	ExtraData        string `gorm:"column:extraData" json:"extraData"`
	// Size             int64  `gorm:"column:size" json:"size"`
	// GasLimit         int64  `gorm:"column:gasLimit" json:"gasLimit"`
	// GasUsed          int64  `gorm:"column:gasUsed" json:"gasUsed"`
	// Timestamp        int64  `gorm:"column:timestamp" json:"timestamp"`
}

// TransactionsCount,UncleCount,BlockReward
type BscBlockTable struct {
	ID int64 `gorm:"column:id;primary_key;AUTO_INCREMENT" json:"id"`
	BscBlockCommon

	GasLimit  int64 `gorm:"column:gasLimit" json:"gasLimit"`
	GasUsed   int64 `gorm:"column:gasUsed" json:"gasUsed"`
	Timestamp int64 `gorm:"column:timestamp" json:"timestamp"`
	Number    int64 `gorm:"column:number;NOT NULL" json:"number"`

	// BlockReward现有项目未用到,数据库无数据
	BlockReward       int64 `gorm:"column:blockReward" json:"blockReward"`
	TransactionsCount int64 `gorm:"column:transactionsCount" json:"transactionsCount"`
	UncleCount        int64 `gorm:"column:uncleCount" json:"uncleCount"`
	// 下面三字段来源于TrustNodeScore
	CreditData  string `gorm:"column:credit_data" json:"credit_data"`
	CreditValue string `gorm:"column:credit_value" json:"credit_value"`
	// miner的当前的(可信值是递增的/最大)可信值
	CreditMax int64 `gorm:"column:credit_max" json:"credit_max"`
	// TotalFee现有项目未用到,数据库无数据
	TotalFee int64 `gorm:"column:total_fee" json:"total_fee"`
	// IsComputed现有项目未用到,数据库无数据
	IsComputed int `gorm:"column:is_computed" json:"is_computed"`
}

func (m *BscBlockTable) TableName() string {
	return "block"
}
