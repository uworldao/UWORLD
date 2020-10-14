package types

type RlpBody struct {
	Transactions []*RlpTransaction
}

func (rb *RlpBody) TranslateToBody() *Body {
	txs := make([]ITransaction, 0)
	for _, rlpTx := range rb.Transactions {
		txs = append(txs, rlpTx.TranslateToTransaction())
	}
	return &Body{Transactions: txs}
}
