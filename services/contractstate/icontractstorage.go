package contractstate

import (
	"github.com/uworldao/UWORLD/common/hasharry"
	"github.com/uworldao/UWORLD/core/types"
)

// Implement storage as contract state
type IContractStorage interface {
	GetContractState(contractAddr string) *types.Contract
	SetContractState(contract *types.Contract)
	InitTrie(contractRoot hasharry.Hash) error
	RootHash() hasharry.Hash
	Commit() (hasharry.Hash, error)
	Close() error
}
