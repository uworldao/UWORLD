package types

type TxPool struct {
	TxsCount      int               `json:"txscount"`
	PreparedCount int               `json:"preparedcount"`
	FutureCount   int               `json:"futurecount"`
	PreparedTxs   []*RpcTransaction `json:"preparedtxs"`
	FutureTxs     []*RpcTransaction `json:"futuretxs"`
}

func TranslateTxsToRpcTxPool(preparedTxs Transactions, futureTxs Transactions) (*TxPool, error) {
	var preparedRpcTxs, futureRpcTxs []*RpcTransaction
	for _, tx := range preparedTxs {
		t, err := TranslateTxToRpcTx(tx.(*Transaction))
		if err != nil {
			return nil, err
		}
		preparedRpcTxs = append(preparedRpcTxs, t)
	}

	for _, tx := range futureTxs {
		t, err := TranslateTxToRpcTx(tx.(*Transaction))
		if err != nil {
			return nil, err
		}
		futureRpcTxs = append(futureRpcTxs, t)
	}

	preparedTxCount := preparedTxs.Len()
	futureTxCount := futureTxs.Len()

	return &TxPool{
		TxsCount:      preparedTxCount + futureTxCount,
		PreparedCount: preparedTxCount,
		FutureCount:   futureTxCount,

		PreparedTxs: preparedRpcTxs,
		FutureTxs:   futureRpcTxs,
	}, nil
}
