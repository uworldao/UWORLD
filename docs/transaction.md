# Transaction 文档

## 目录

### 工具包
```
github.com/jhdriver/UWORLD/ut/transaction
github.com/jhdriver/UWORLD/ut
github.com/jhdriver/UWORLD/common/hasharry
github.com/jhdriver/UWORLD/param
github.com/jhdriver/UWORLD/rpc
github.com/jhdriver/UWORLD/rpc/rpctypes
```

### 创建交易
```
from := hasharry.StringToAddress("UbQyzkoPBnWMMtzX946eTJiKcRgVpDtaUoe")
to := hasharry.StringToAddress("UbaJeMrs9EKbBjGovzBSSCKQ4BrfWBPt9tu")
token := hasharry.StringToAddress("UtpuFryDEPGLdYmitofgctYW214AsstWrfG")
tx := transation.NewTransaction(from, to, token, "note string", 100000000, 100000, 1)
```

### 创建代币
```
from := hasharry.StringToAddress("UbQyzkoPBnWMMtzX946eTJiKcRgVpDtaUoe")
to := hasharry.StringToAddress("UbQyzkoPBnWMMtzX946eTJiKcRgVpDtaUoe")
coinAbbr := "TC"
coinName := "TEST COIN"
contract := ut.GenerateUBAddress(param.MainNet, from, coinAbbr)
tx := transation.NewContract(from, to, contract, "note string", 10000000000000, 100000, 1, "name", "abbr string", true)
```


### 消息签名
```
tx.SignTx(private)
```

### 发送交易

```
rpcTx := rpctypes.TranslateTxToRpcTx(tx)
jsonBytes, _ := json.Marshal(rpcTx)
req := &rpc.Request{Params: jsonBytes}
ctx, _ := context.WithTimeout(context.TODO(), time.Second*20)
resp, err := rpcClient.SendTransaction(ctx, req)
```
