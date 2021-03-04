package core

import (
	"github.com/uworldao/UWORLD/common/hasharry"
	"github.com/uworldao/UWORLD/core/types"
)

type IBlockChain interface {
	CurrentHeader() (*types.Header, error)

	GetConfirmedHeight() uint64

	GetLastHeight() uint64

	GetHeaderByHeight(height uint64) (*types.Header, error)

	GetHeaderByHash(hash hasharry.Hash) (*types.Header, error)

	GetBlockByHeight(height uint64) (*types.Block, error)

	GetBlockByHash(hash hasharry.Hash) (*types.Block, error)

	GetRlpBlockByHeight(height uint64) (*types.RlpBlock, error)

	GetRlpBlockByHash(hash hasharry.Hash) (*types.RlpBlock, error)

	GetTransaction(hash hasharry.Hash) (types.ITransaction, error)

	GetTransactionIndex(hash hasharry.Hash) (types.ITransactionIndex, error)

	GetAddressVote(address hasharry.Address) uint64

	GetTermLastHash(term uint64) (hasharry.Hash, error)

	InsertChain(block *types.Block) error

	SaveGenesisBlock(block *types.Block) error

	UpdateCurrentBlockHeight(height uint64)

	UpdateConfirmedHeight(height uint64)

	ValidationBlockHash(header *types.Header) (bool, error)

	FallBack()

	FallBackTo(int64 uint64) error

	StateRoot() hasharry.Hash

	ContractRoot() hasharry.Hash

	ConsensusRoot() hasharry.Hash

	TireRoot() (hasharry.Hash, hasharry.Hash, hasharry.Hash)

	CloseStorage() error
}
