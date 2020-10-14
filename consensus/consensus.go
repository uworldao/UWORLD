package consensus

import (
	"github.com/jhdriver/UWORLD/common/hasharry"
	"github.com/jhdriver/UWORLD/core/types"
)

// IConsensus consensus interface to handle all matters concerning consensus
type IConsensus interface {
	IConsensusVerify
	IDPos
	IDPosTrie
}

type IDPos interface {
	Init(chain IChain) error

	// Modify the dpos and finally confirm the block header
	SetConfirmedHeader(header *types.Header)

	// Super node signature function
	Sign(block *types.Block) error

	// Get the genesis block
	GetGenesisBlock() *types.Block

	// Get cycle interval
	GetTermInterval() uint64

	// Get the super node id of the time period
	GetWinnersPeerID(time uint64) ([]string, error)

	// Get the dpos and finally confirm the block header
	GetConfirmedBlockHeader(chain IChain) *types.Header

	// Get current candidate
	GetCandidates(chain IChain) []*types.Candidate

	// Get a super node of a certain period
	GetTermWinners(term uint64) *types.Winners

	// Get the number of blocks produced by a super node in a certain period
	GetTermWinnersMntCount(term uint64, address hasharry.Address) (uint64, error)

	// Check whether a block
	CheckWinner(chain IChain, header *types.Header) error

	// Update dpos status
	UpdateConsensus(block *types.Block)
}

// consensus verify
type IConsensusVerify interface {
	VerifyHeader(header, parent *types.Header) error

	VerifySeal(chain IChain, header *types.Header, parents *types.Header) error

	VerifyTx(tx types.ITransaction) error
}

// DPos trie
type IDPosTrie interface {
	// Initialize dpos trie
	InitTrie(consensusRoot hasharry.Hash) error

	// Commit dpos trie
	Commit() (hasharry.Hash, error)

	// Get trie roothash
	RootHash() hasharry.Hash

	// Close the trie database
	Close() error
}

// Get blockchian information interface
type IChain interface {
	// Get the latest block header
	CurrentHeader() (*types.Header, error)

	// Get block header through height
	GetHeaderByHeight(height uint64) (*types.Header, error)

	// Get block header through hash
	GetHeaderByHash(hash hasharry.Hash) (*types.Header, error)

	// Get block through hash
	GetBlockByHash(hash hasharry.Hash) (*types.Block, error)

	// Get block through height
	GetBlockByHeight(height uint64) (*types.Block, error)

	// Get the number of votes cast by the address
	GetAddressVote(address hasharry.Address) uint64

	// Get the hash of the last block of the previous cycle
	GetTermLastHash(term uint64) (hasharry.Hash, error)

	// Update the header of the final block
	UpdateConfirmedHeight(height uint64)
}

//Consensus signature interface
type ISign interface {
	// Sign the hash
	SignHash(hash hasharry.Hash) (*types.SignScript, error)
}
