package core

import (
	"github.com/uworldao/UWORLD/core/types"
)

type ErrTxs struct {
	errTxs chan types.Transactions
}

func NewErrTxs() *ErrTxs {
	return &ErrTxs{errTxs: make(chan types.Transactions, 50)}
}

func (e *ErrTxs) Add(txs types.Transactions) {
	e.errTxs <- txs
}
