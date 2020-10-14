# UWorld chain

## How to build

####  Prerequisites

- Update Go to version at least 1.13 (required >= **1.13**)

Check your golang version

```bash
~ go version
go version go1.13 darwin/amd64
```

```bash
cd ub_chain
go build

cd cmd/wallet
go build
```

## How to use


##### Copy configuration file for reconfiguration

```bash
 cp config.toml.example config.toml
```

##### Modify configuration file

* set RpcPass
* set ExternalIp


##### Start the ub_chain

```bash

./UWorld --config config.toml
```

##### Copy wallet configuration file for reconfiguration

```
cd cmd/wallet
cp wallet.toml.example wallet.toml
```

##### Modify wallet configuration file

* set RpcIp
* set RpcPass
* If the node has the RpcTLS switch turned on, you need to configure the node's server.pem path to RpcCert and set RpcTLS in wallet.config to true

##### Use wallet

```bash
./wallet --help
```
##### Create an account or set password at create

```bash
./wallet CreateAccount 

./wallet CreateAccount 123456
```
##### Send transaction

./wallet SendTransaction from to contract amount fee [password]

```bash
./wallet SendTransaction 3ajDe9zSANwuTBL6xBEj5ZWjjbWYQyzBohv1  3ajHhfRK5ZDz9TvjrXqhq2deLo8qk37zakxq  UWD 1000 0.0003

./wallet SendTransaction 3ajDe9zSANwuTBL6xBEj5ZWjjbWYQyzBohv1  3ajHhfRK5ZDz9TvjrXqhq2deLo8qk37zakxq  UWD 1000 0.0003 123456
```

##### Get account balance

```bash
./wallet GetAccount 3ajDe9zSANwuTBL6xBEj5ZWjjbWYQyzBohv1
```
  
    
## Documents


- [RPC](https://github.com/jhdriver/UWORLD/tree/master/docs/rpc.md)
- [Address](https://github.com/jhdriver/UWORLD/tree/master/docs/address.md)
- [Transaction](https://github.com/jhdriver/UWORLD/tree/master/docs/transaction.md)
