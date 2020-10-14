package types

type RpcBody struct {
	Transactions []*RpcTransaction `json:"transactions"`
}

func TranslateBodyToRpcBody(body *Body) (*RpcBody, error) {
	var rpcTxs []*RpcTransaction
	for _, tx := range body.Transactions {
		rpcTx, err := TranslateTxToRpcTx(tx.(*Transaction))
		if err != nil {
			return nil, err
		}
		rpcTxs = append(rpcTxs, rpcTx)
	}
	return &RpcBody{rpcTxs}, nil
}
