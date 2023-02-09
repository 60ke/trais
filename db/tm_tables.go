package db

// class Transaction(BaseModel):
//     block_hash = CharField(column_name='blockHash')
//     block_number = BigIntegerField(column_name='blockNumber')
//     gas = BigIntegerField()
//     gas_price = BigIntegerField(column_name='gasPrice')
//     fee = BigIntegerField()
//     gas_used = BigIntegerField(column_name='gasUsed')
//     hash = CharField()
//     id = BigAutoField()
//     nonce = BigIntegerField()
//     r = CharField()
//     s = CharField()
//     source = CharField()
//     timestamp = BigIntegerField()
//     to = CharField()
//     transaction_index = BigIntegerField(column_name='transactionIndex')
//     tx_str = TextField()
//     type = IntegerField()
//     v = BigIntegerField()
//     value = CharField()

type Transaction struct {
	ID               int64  `gorm:"primaryKey",column:bsc_current_block`
	Hash             string `gorm:"column:hash"`
	Timestamp        int64  `gorm:"column:timestamp"`
	BlockHash        string `gorm:"column:blockHash"`
	BlockNumber      int64  `gorm:"column:blockNumber"`
	Gas              int64  `gorm:"column:blockNumber"`
	GasPrice         int64
	Fee              int64
	GasUsed          int64
	Nonce            int64
	R                string
	S                string
	V                string
	Source           string
	To               string
	TransactionIndex int64
	TxStr            string `gorm:"type:text"`
	// 从python项目迁移而来暂时不清楚此字段的作用
	Type  int
	Value string
}

type Transaction1 struct {
	ID               int64  `gorm:"primaryKey",column:bsc_current_block`
	BlockHash        string `json:"blockHash"`
	BlockNumber      string `json:"blockNumber"`
	From             string `json:"from"`
	Gas              string `json:"gas"`
	GasPrice         string `json:"gasPrice"`
	Hash             string `json:"hash"`
	Input            string `json:"input"`
	Nonce            string `json:"nonce"`
	To               string `json:"to"`
	TransactionIndex string `json:"transactionIndex"`
	Value            string `json:"value"`
	Type             string `json:"type"`
	V                string `json:"v"`
	R                string `json:"r"`
	S                string `json:"s"`
}

// class Address(BaseModel):
//     address = CharField(unique=True)
//     balance = BigIntegerField()
//     id = BigAutoField()
//     received = CharField()
//     sent = CharField()
//     time = BigIntegerField()
//     tx_count = BigIntegerField(column_name='txCount')

//     class Meta:
//         table_name = 'address'

type Account struct {
	ID        int64  `gorm:"primaryKey",column:bsc_current_block`
	Address   string `gorm:"column:address"`
	Balance   int64  `gorm:"column:balance"`
	Received  string `gorm:"column:received"`
	Sent      int64  `gorm:"column:sent"`
	Timestamp int64  `gorm:"column:timestamp"`
	tx_count  int64
}

// class Block(BaseModel):
//     block_reward = BigIntegerField(column_name='blockReward')
//     credit_data = CharField()
//     credit_max = BigIntegerField()

//     credit_value = JSONField()  # json
//     difficulty = CharField()
//     extra_data = CharField(column_name='extraData')
//     gas_limit = BigIntegerField(column_name='gasLimit')
//     gas_used = BigIntegerField(column_name='gasUsed')
//     hash = CharField(unique=True)
//     id = BigAutoField()
//     logs_bloom = CharField(column_name='logsBloom')
//     miner = CharField()
//     nonce = CharField()
//     number = BigIntegerField(unique=True)
//     parent_hash = CharField(column_name='parentHash', unique=True)
//     receipts_root = CharField(column_name='receiptsRoot')
//     sha3_uncles = CharField(column_name='sha3Uncles')
//     size = BigIntegerField()
//     state_root = CharField(column_name='stateRoot')
//     timestamp = BigIntegerField()
//     total_difficulty = CharField(column_name='totalDifficulty')
//     transactions_count = BigIntegerField(column_name='transactionsCount')
//     transactions_root = CharField(column_name='transactionsRoot')
//     uncle_count = BigIntegerField(column_name='uncleCount')

//     class Meta:
//         table_name = 'block'

type Block struct {
	Difficulty       string        `json:"difficulty"`
	ExtraData        string        `json:"extraData"`
	GasLimit         string        `json:"gasLimit"`
	GasUsed          string        `json:"gasUsed"`
	Hash             string        `json:"hash"`
	LogsBloom        string        `json:"logsBloom"`
	Miner            string        `json:"miner"`
	MixHash          string        `json:"mixHash"`
	Nonce            string        `json:"nonce"`
	Number           string        `json:"number"`
	ParentHash       string        `json:"parentHash"`
	ReceiptsRoot     string        `json:"receiptsRoot"`
	Sha3Uncles       string        `json:"sha3Uncles"`
	Size             string        `json:"size"`
	StateRoot        string        `json:"stateRoot"`
	Timestamp        string        `json:"timestamp"`
	TotalDifficulty  string        `json:"totalDifficulty"`
	Transactions     []string      `json:"transactions"`
	TransactionsRoot string        `json:"transactionsRoot"`
	TrustNodeScore   string        `json:"trustNodeScore"`
	Uncles           []interface{} `json:"uncles"`
}

// class Info(BaseModel):
//     addresses = BigIntegerField()
//     hash_rate = BigIntegerField(column_name='hashRate')
//     id = BigAutoField()
//     last_block = BigIntegerField(column_name='lastBlock')
//     last_block_fees = BigIntegerField(column_name='lastBlockFees')
//     last_transaction_fees = BigIntegerField(column_name='lastTransactionFees')
//     timestamp = BigIntegerField()
//     total_difficulty = CharField(column_name='totalDifficulty')
//     transactions = BigIntegerField()
//     unconfirmed = BigIntegerField()

//     class Meta:
//         table_name = 'info'
