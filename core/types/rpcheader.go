package types

import (
	"time"
)

type RpcHeader struct {
	Hash          string    `json:"hash"`
	ParentHash    string    `json:"parenthash"`
	TxRoot        string    `json:"txroot"`
	StateRoot     string    `json:"stateroot"`
	ContractRoot  string    `json:"contractroot"`
	ConsensusRoot string    `json:"consensusroot"`
	Height        uint64    `json:"height"`
	Time          time.Time `json:"time"`
	Term          uint64    `json:"term"`
	Signer        string    `json:"signer"`
}

func TranslateHeaderToRpcHeader(header *Header) *RpcHeader {
	signer := header.Signer.String()
	return &RpcHeader{
		Hash:          header.HashString(),
		ParentHash:    header.ParentHashString(),
		TxRoot:        header.TxRoot.String(),
		StateRoot:     header.StateRoot.String(),
		ContractRoot:  header.ContractRoot.String(),
		ConsensusRoot: header.ConsensusRoot.String(),
		Height:        header.Height,
		Time:          time.Unix(int64(header.Time), 0),
		Term:          header.Term,
		Signer:        signer,
	}
}
