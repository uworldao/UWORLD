package core

import "github.com/jhdriver/UWORLD/core/types"

// Transaction pool interface, which is used to manage the transaction pool
type ITxPool interface {
	Start() error
	Stop() error
	Add(tx types.ITransaction, isPeer bool) error
	Gets(count int) types.Transactions
	GetAll() (types.Transactions, types.Transactions)
	Get() types.ITransaction
	Remove(txs types.Transactions)
	IsExist(tx types.ITransaction) bool
}
