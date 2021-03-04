package types

import (
	"github.com/uworldao/UWORLD/common/hasharry"
)

// Account Status
type IAccount interface {
	GetBalance(contract string) uint64
	GetNonce() uint64
	Update(confirmedHeight uint64) error
	StateKey() hasharry.Address
	IsExist() bool
	IsNeedUpdate() bool
	FromChange(tx ITransaction, blockHeight uint64) error
	ToChange(tx ITransaction, blockHeight uint64) error
	FeesChange(fees, blockHeight uint64)
	VerifyTxState(tx ITransaction) error
	VerifyNonce(nonce uint64) error
	IsEmpty() bool
}
