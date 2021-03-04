# Transaction 文档

## 目录

### 工具包
```
github.com/uworldao/UWORLD/ut/transaction
github.com/uworldao/UWORLD/ut
github.com/uworldao/UWORLD/common/hasharry
github.com/uworldao/UWORLD/param
github.com/uworldao/UWORLD/rpc
github.com/uworldao/UWORLD/rpc/rpctypes
```

### 创建交易
```
from := "UWDb62iQKvD4z6QqJW4rYobvLbfmPEBskg5h"
to := "UWDb62iQKvD4z6QqJW4rYobvLbfmPEBskg5h"
token := "UWD"
tx := transation.NewTransaction(from, to, token, "note string", 100000000, 1)
```

### 创建代币
```
from := "UbQyzkoPBnWMMtzX946eTJiKcRgVpDtaUoe"
to := "UbQyzkoPBnWMMtzX946eTJiKcRgVpDtaUoe"
coinAbbr := "TC"
coinName := "TEST COIN"
decription := "Test"
contract := ut.GenerateUWDAddress(param.MainNet, from, coinAbbr)
contract := transation.NewContract(from, to, contract, "note string", 10000000000000, 1, "name", "abbr string", true, decription)
```


### 消息签名
```
tx.SignTx(private)
```

### 发送交易

```
rpcTx, err := types.TranslateTxToRpcTx(tx)
jsonBytes, err := json.Marshal(rpcTx)
ctx, _ := context.WithTimeout(context.TODO(), time.Second*20)
resp, err := rpcClient.SendTransaction(ctx, rpc.Bytes{Bytes: jsonBytes})
```
