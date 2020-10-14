package types

type RpcLoginTransactionBody struct {
	PeerId string `json:"peerid"`
}

func (rlt *RpcLoginTransactionBody) PeerIdBytes() []byte {
	return []byte(rlt.PeerId)
}
