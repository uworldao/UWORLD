# RPC 文档

## 介绍
GRPC实现
* go语言调用，使用现有的[Client](https://github.com/uworldao/UWORLD/blob/master/rpc/rpcclient.go)
* 其他语言，使用[proto文件](https://github.com/uworldao/UWORLD/blob/master/rpc/rpc.proto)生成rpcclient

```
type customCredential struct {
	OpenTLS  bool
	Password string
}

var opts []grpc.DialOption
opts = append(opts, grpc.WithPerRPCCredentials(&customCredential{Password: "123", OpenTLS: false}))
conn, _ = grpc.Dial("127.0.0.1:19161", opts...)
gc = NewGreeterClient(conn)
```

## 目录

### GetAccount
- info：获取账户信息
- result:
    
```json
{
    "address": "UWDKoLj4mRTKr4SjyyFG4LY3ExZVSZT9dNZv",
    "nonce": 0,
    "time": 0,
    "coins": [
        {
            "contract": "UWD",
            "balance": 3045.0003,
            "lockedout": 3,
            "lockedin": 0
        }
    ],
    "confirmedheight": 11203,
    "confirmednonce": 0,
    "confirmedtime": 0
}
```

### SendTransaction
- info：发送交易

### GetTransaction
- info：获取交易
- result:
    
```json
{
    "txhead": {
        "txhash": "0xbc7c8d4fa7d24915aa877f33a6a3801437df7d209d27528945b4a51488135b9e",
        "txtype": 0,
        "from": "coinbase",
        "nonce": 0,
        "fees": 0,
        "time": 1597130625,
        "note": "",
        "signscript": {
            "signature": "",
            "pubkey": ""
        }
    },
    "normalbody": {
        "contract": "UBC",
        "to": "UWDKoLj4mRTKr4SjyyFG4LY3ExZVSZT9dNZv",
        "amount": 300000000
    },
}
```

### GetBlockByHash
- info：获取block
- result:
    
```json
{
    "header": {
        "hash": "0x10917fa77060fcd1d6bdf0ea2e98c5514fea9dc9c06051d13e46c6f7430f80ec",
        "parenthash": "0x89f05afa3462bec7e5e8d7666b489a3c5820150d06259cd479be7164c99d5bf3",
        "txroot": "0x1b6c8a1596ddc3059cd329f129e7f8789c9899e26941da215aa7433baf79c608",
        "stateroot": "0xae185c9799660361604c4add0190efdb94d6f885e29205e2bdc6b0077cbaf7d9",
        "contractroot": "0x56e81f171bcc55a6ff8345e692c0f86e5b48e01b996cadc001622fb5e363b421",
        "consensusroot": "0xfdaf25615745cdd48157631a25da5ed181c2db0276fa7178638ff3ce1d44ef5e",
        "height": 10,
        "time": "2020-08-11T15:23:45+08:00",
        "term": 0,
        "signer": "UWDKoLj4mRTKr4SjyyFG4LY3ExZVSZT9dNZv"
    },
    "body": {
        "transactions": [
            {
                "txhead": {
                    "txhash": "0xbc7c8d4fa7d24915aa877f33a6a3801437df7d209d27528945b4a51488135b9e",
                    "txtype": 0,
                    "from": "coinbase",
                    "nonce": 0,
                    "fees": 0,
                    "time": 1597130625,
                    "note": "",
                    "signscript": {
                        "signature": "",
                        "pubkey": ""
                    }
                },
                "normalbody": {
                    "contract": "UWD",
                    "to": "UWDKoLj4mRTKr4SjyyFG4LY3ExZVSZT9dNZv",
                    "amount": 300000000
                },
            }
        ]
    },
    "confirmed": true
}
```
### GetBlockHeight
- info：获取block
- result:
    
```json
{
    "header": {
        "hash": "0x10917fa77060fcd1d6bdf0ea2e98c5514fea9dc9c06051d13e46c6f7430f80ec",
        "parenthash": "0x89f05afa3462bec7e5e8d7666b489a3c5820150d06259cd479be7164c99d5bf3",
        "txroot": "0x1b6c8a1596ddc3059cd329f129e7f8789c9899e26941da215aa7433baf79c608",
        "stateroot": "0xae185c9799660361604c4add0190efdb94d6f885e29205e2bdc6b0077cbaf7d9",
        "contractroot": "0x56e81f171bcc55a6ff8345e692c0f86e5b48e01b996cadc001622fb5e363b421",
        "consensusroot": "0xfdaf25615745cdd48157631a25da5ed181c2db0276fa7178638ff3ce1d44ef5e",
        "height": 10,
        "time": "2020-08-11T15:23:45+08:00",
        "term": 0,
        "signer": "UWDKoLj4mRTKr4SjyyFG4LY3ExZVSZT9dNZv"
    },
    "body": {
        "transactions": [
            {
                "txhead": {
                    "txhash": "0xbc7c8d4fa7d24915aa877f33a6a3801437df7d209d27528945b4a51488135b9e",
                    "txtype": 0,
                    "from": "coinbase",
                    "nonce": 0,
                    "fees": 0,
                    "time": 1597130625,
                    "note": "",
                    "signscript": {
                        "signature": "",
                        "pubkey": ""
                    }
                },
                "normalbody": {
                    "contract": "UWD",
                    "to": "UWDKoLj4mRTKr4SjyyFG4LY3ExZVSZT9dNZv",
                    "amount": 300000000
                },
            }
        ]
    },
    "confirmed": true
}
```

### GetLastHeight
- info：获取最高高度
- result: 高度(string bytes)

### GetConfirmedHeight
- info：获取已经确认的最高区块高度
- result: 高度(string bytes)

### GetPoolTxs
- info：获取交易池
- result: 高度(string bytes)
```json
{
    "txscount": 1,
    "preparedcount": 1,
    "futurecount": 0,
    "preparedtxs": [
        {
            "txhead": {
                "txhash": "0x786315263b74fef17b227cb74b940cae456deb33d034fda3f3170a82abfe17b5",
                "txtype": 0,
                "from": "UWDKoLj4mRTKr4SjyyFG4LY3ExZVSZT9dNZv",
                "nonce": 3,
                "fees": 100000,
                "time": 1597730820,
                "note": "1",
                "signscript": {
                    "signature": "30440220472593b3a8cbe98b8487c5b5f5891787dedfaf1857b05a368524b7fa0e6b42a40220277d5fb0f09fb0c69e228118e55722c2f7755ad93a248d49bee3adfdf5fac317",
                    "pubkey": "03ec37e27994fd9c6c12958d2f46d87ec2d0930804a4da6741317eeeca8af5e5a5"
                }
            },
            "normalbody": {
                "contract": "UWD",
                "to": "UWDKoLj4mRTKr4SjyyFG4LY3ExZVSZT9dNZv",
                "amount": 2000000000
            }
        }
    ],
    "futuretxs": null
}
```

### GetContract
- info：获取发币详情
- result:
```json
{
    "contract": "UWTKoLj4mRTKr4SjyyFG4LY3ExZVSZT9dNZv",
    "name": "Test Coin",
    "abbr": "TCC",
    "increase": false,
    "records": [
        {
            "height": 39963,
            "txhash": "0x1a28af0225cda0aa2b36793cb44e892c6679a78c7c80850f4f7852fd8b0fedfe",
            "time": 1597731050,
            "amount": 1000
        }
    ]
}
```

### Peers
- info：获取p2p节点信息
- result:
```json
[
    {
        "version": "v0.3.1",
        "net": "TN",
        "p2pid": "16Uiu2HAm4Wqt9qQbvBAUqfvzV52W5Eett3yzX3ujhzBqxmDXbe7R",
        "p2pAddr": "[/ip4/47.57.100.253/tcp/28000][/ip4/172.31.244.159/tcp/28000]",
        "connections": 11,
        "height": 39957,
        "confirmed": 39949
    }
]
```
### NodeInfo
- info：获取本地节点信息
- result:
```json
{
    "version": "v0.3.1",
    "net": "TN",
    "p2pid": "16Uiu2HAkxQruPEeNsC5adX5PJ66rzZkYreM5q76utx8doQ2Kcrfd",
    "p2pAddr": "[/ip4/0.0.0.0/tcp/30000][/ip4/192.168.31.140/tcp/30000]",
    "connections": 1,
    "height": 39958,
    "confirmed": 39950
}
```