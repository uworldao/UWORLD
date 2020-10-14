package types

type RpcNormalTransactionBody struct {
	Contract string `json:"contract"`
	To       string `json:"to"`
	Amount   uint64 `json:"amount"`
}

func (rnb *RpcNormalTransactionBody) ToBytes() []byte {
	return []byte(rnb.To)
}
