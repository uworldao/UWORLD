package types

type RpcBlock struct {
	RpcHeader *RpcHeader `json:"header"`
	RpcBody   *RpcBody   `json:"body"`
	Confirmed bool       `json:"confirmed"`
}

func TranslateBlockToRpcBlock(block *Block, confirmHeight uint64) (*RpcBlock, error) {
	rpcHeader := TranslateHeaderToRpcHeader(block.Header)
	rpcBody, err := TranslateBodyToRpcBody(block.Body)
	if err != nil {
		return nil, err
	}
	return &RpcBlock{RpcHeader: rpcHeader, RpcBody: rpcBody, Confirmed: confirmHeight >= rpcHeader.Height}, nil
}
