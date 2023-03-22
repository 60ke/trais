package web3

import (
	"bytes"
	"encoding/json"
	"fmt"
	"testing"
)

func TestLatestBlockNumber(t *testing.T) {
	// body, _ := LatestBlockNumber("http://106.3.133.179:8545")
	var resp LatestBlockNumberResp
	// decoder := json.NewDecoder(bytes.NewReader(body))
	decoder := json.NewDecoder(bytes.NewReader([]byte(`{"jsonrpc":"2.0","id":83,"result1":"0x5c76a"}`)))

	decoder.DisallowUnknownFields()
	err := decoder.Decode(&resp)
	fmt.Println(resp, err)

}
