package types

type RpcGetBlockByHeightParam struct {
	Height uint64 `json:"height"`
}

type RpcGetBlockByHashParam struct {
	Hash string `json:"hash"`
}
