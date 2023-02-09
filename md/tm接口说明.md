tm版本：https://github.com/trias-lab/tmware/tree/upgrade-0.35.0-libp2p
- 根据交易hash获取对应的交易(交易哈希需有0x前缀，大小写不敏感)：
```curl
curl -X GET "http://127.0.0.1:46657/tri_tx?hash=0x37b5be61fefd1edea7761c41767c6c4e8bda09d7e7b73002b6ef98b481bb68a5&prove=true" -H "accept: application/json"




{
  "jsonrpc": "2.0",
  "id": -1,
  "result": {
    "hash": "37B5BE61FEFD1EDEA7761C41767C6C4E8BDA09D7E7B73002B6EF98B481BB68A5",
    "height": "23",
    "index": 0,
    "tx_result": {
      "code": 0,
      "data": null,
      "log": "",
      "info": "",
      "gas_wanted": "0",
      "gas_used": "0",
      "events": [
        {
          "type": "app",
          "attributes": [
            {
              "key": "creator",
              "value": "Cosmoshi Netowoko",
              "index": true
            },
            {
              "key": "key",
              "value": "A11",
              "index": true
            },
            {
              "key": "index_key",
              "value": "index is working",
              "index": true
            },
            {
              "key": "noindex_key",
              "value": "index is working",
              "index": false
            }
          ]
        }
      ],
      "codespace": ""
    },
    "tx": "QTEx",
    "proof": {
      "root_hash": "85670B8F34B16DBB9502C2D338EB2EC6B3403A96EA88311D74A9705E45EA1C4A",
      "data": "QTEx",
      "proof": {
        "total": "1",
        "index": "0",
        "leaf_hash": "hWcLjzSxbbuVAsLTOOsuxrNAOpbqiDEddKlwXkXqHEo=",
        "aunts": []
      }
    }
  }
}%


```


- 发送交易

```curl
curl "http://127.0.0.1:46657/tri_broadcast_tx_sync?tx=%22A11%22"

{
  "jsonrpc": "2.0",
  "id": -1,
  "result": {
    "code": 0,
    "data": "",
    "log": "",
    "codespace": "",
    "mempool_error": "",
    "hash": "37B5BE61FEFD1EDEA7761C41767C6C4E8BDA09D7E7B73002B6EF98B481BB68A5"
  }
}

```


- 根据块高获取区块信息
```
curl 'http://127.0.0.1:46657/tri_block_info?height=24'  -H "accept: application/json"
```


- 交易内容解码及哈希计算
```py
b64data = base64.b64decode(tx)
m = hashlib.sha256()
m.update(b64data)
str_hash = m.hexdigest()
tx_data = b64data.decode()
```