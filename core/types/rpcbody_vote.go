package types

type RpcVoteTransactionBody struct {
	To string `json:"to"`
}

func (rvt *RpcVoteTransactionBody) ToBytes() []byte {
	return []byte(rvt.To)
}
