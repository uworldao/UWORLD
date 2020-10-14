package core

import (
	"github.com/jhdriver/UWORLD/common/hasharry"
	"github.com/jhdriver/UWORLD/core/types"
)

type IContractState interface {
	GetContract(contractAddr string) *types.Contract

	VerifyState(tx types.ITransaction) error

	UpdateContract(tx types.ITransaction, blockHeight uint64)

	UpdateConfirmedHeight(height uint64)

	InitTrie(hash hasharry.Hash) error

	RootHash() hasharry.Hash

	ContractTrieCommit() (hasharry.Hash, error)

	Close() error
}
