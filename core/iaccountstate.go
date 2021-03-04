package core

import (
	"github.com/uworldao/UWORLD/common/hasharry"
	"github.com/uworldao/UWORLD/core/types"
)

type IAccountState interface {
	InitTrie(stateRoot hasharry.Hash) error

	GetAccountState(stateKey hasharry.Address) types.IAccount

	GetAccountNonce(stateKey hasharry.Address) (uint64, error)

	UpdateFrom(tx types.ITransaction, blockHeight uint64) error

	UpdateTo(tx types.ITransaction, blockHeight uint64) error

	UpdateFees(fees, blockHeight uint64) error

	UpdateConfirmedHeight(height uint64)

	VerifyState(tx types.ITransaction) error

	StateTrieCommit() (hasharry.Hash, error)

	RootHash() hasharry.Hash

	Print()

	Close() error
}
