# 地址生成 文档
## 目录

### 工具包
```
github.com/uworldao/UWORLD/param
github.com/uworldao/UWORLD/ut
github.com/uworldao/UWORLD/common/hasharry
```


###  生成BIP39助记词

```
e, _ := ut.Entropy()
m, _ := ut.Mnemonic(e)
```

### 生成secp256k1私钥

```
key, _ := ut.MnemonicToEc(m)
```

###  生成地址

```
addr, _ := ut.GenerateUWDAddress(param.TestNet, key.PubKey())
```

### 校验地址
```
ut.CheckUWDAddress(param.TestNet, "UbQQhJ4zmp4wLQ4Li6tm7zigopaeGrWxSvy")
```

### 生成token地址

```
ut.GenerateTokenAddress(param.TestNet, "UbQQhJ4zmp4wLQ4Li6tm7zigopaeGrWxSvy", "HFC")
```

### 校验token地址

```
ut.CheckTokenAddress(param.TestNet, "UbQQhJ4zmp4wLQ4Li6tm7zigopaeGrWxSvy")
```