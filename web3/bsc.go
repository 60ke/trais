package web3

import (
	"fmt"
	"strings"

	"github.com/60ke/trais/tools"
)

var BSCTIMEOUT = 5

func LatestBlockNumber(rpc string) ([]byte, error) {
	var data = strings.NewReader(`{
		"jsonrpc":"2.0",
		"method":"eth_blockNumber",
		"params":[],
		"id":83
	}`)
	return tools.Post(rpc, BSCTIMEOUT, data)
}

func GetBlockByNum(rpc, block string) ([]byte, error) {
	var data = strings.NewReader(fmt.Sprintf(`{
		"jsonrpc":"2.0",
		"method":"eth_getBlockByNumber",
		"params":[
			"%s", 
			true
		],
		"id":1
	}`, block))
	return tools.Post(rpc, BSCTIMEOUT, data)

}

func GetTxByHash(rpc, hash string) {

}

func GetTxrByHash(rpc, hash string) {

}

func GetBalance(rpc string, addr string) ([]byte, error) {
	var data = strings.NewReader(fmt.Sprintf(`{
		"jsonrpc":"2.0",
		"method":"eth_getBalance",
		"params":[
			%s, 
			"latest"
		],
		"id":1
	}`, addr))
	return tools.Post(rpc, BSCTIMEOUT, data)
}
