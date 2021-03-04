package core

import (
	"github.com/uworldao/UWORLD/common/hasharry"
	"github.com/uworldao/UWORLD/core/types"
)

type IBlockChainStorage interface {
	GetLastHeight() (uint64, error)

	GetHeaderByHeight(height uint64) (*types.Header, error)

	GetHeader(hash hasharry.Hash) (*types.Header, error)

	GetTransactions(txRoot hasharry.Hash) ([]*types.RlpTransaction, error)

	GetTransaction(hash hasharry.Hash) (*types.RlpTransaction, error)

	GetTxLocation(hash hasharry.Hash) (*types.TxLocation, error)

	GetStateRoot() (hasharry.Hash, error)

	GetContractRoot() (hasharry.Hash, error)

	GetConsensusRoot() (hasharry.Hash, error)

	GetHistoryConfirmedHeight(height uint64) (uint64, error)

	GetTermLastHash(term uint64) (hasharry.Hash, error)

	UpdateLastHeight(height uint64)

	UpdateHeader(header *types.Header)

	UpdateTransactions(txRoot hasharry.Hash, txs []*types.RlpTransaction)

	UpdateTxLocation(txLocs map[hasharry.Hash]*types.TxLocation)

	UpdateHeightHash(height uint64, hash hasharry.Hash)

	UpdateStateRoot(hash hasharry.Hash)

	UpdateContractRoot(hash hasharry.Hash)

	UpdateConsensusRoot(hash hasharry.Hash)

	UpdateHistoryConfirmedHeight(height uint64, confirmedHeight uint64)

	UpdateTermLastHash(term uint64, hash hasharry.Hash)

	Close() error
}
