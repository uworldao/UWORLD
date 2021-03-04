package types

import (
	"bytes"
	hash2 "github.com/uworldao/UWORLD/common/hasharry"
	"github.com/uworldao/UWORLD/crypto/hash"
)

// Block structure, including all transaction information
type Body struct {
	Transactions Transactions
}

func NewBody(txs Transactions) *Body {
	return &Body{
		Transactions: txs,
	}
}

func (body *Body) HashTransactions() hash2.Hash {
	var txsHash [][]byte
	for _, tx := range body.Transactions {
		txsHash = append(txsHash, tx.Hash().Bytes())
	}
	hashBytes := bytes.Join(txsHash, []byte{})
	return hash.Hash(hashBytes)
}

func (body *Body) TranslateToRlpBody() *RlpBody {
	rTxs := make([]*RlpTransaction, 0)
	for _, tx := range body.Transactions {
		rTxs = append(rTxs, tx.TranslateToRlpTransaction())
	}
	return &RlpBody{Transactions: rTxs}
}
