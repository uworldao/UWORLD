package dpos

import (
	"github.com/jhdriver/UWORLD/common/hasharry"
	"github.com/jhdriver/UWORLD/core/types"
)

// DPos storage data interface
type IDPosStorage interface {
	// Store a period of super nodes
	SetTermWinners(term uint64, winners *types.Winners) error

	// Get a period of super nodes
	GetTermWinners(term uint64) (*types.Winners, error)

	// Add new candidate
	SetCandidate(can *types.Candidate) error

	// Get all candidates
	GetCandidates() (*types.Candidates, error)

	// Remove from candidates
	DeleteCandidate(can *types.Candidate) error

	// Get votes for candidates
	GetCandidateVoters(addr hasharry.Address) []hasharry.Address

	// vote from address to address
	SetVoter(from, to hasharry.Address) error

	// Read the last confirmed block header
	GetConfirmedBlockHash() (hasharry.Hash, error)

	// Store the last confirmed block header
	SetConfirmedBlockHash(hash hasharry.Hash)

	// Add the number of blocks to the super node address of a certain period
	SetTermWinnerMintCnt(term uint64, address hasharry.Address)

	// Read the number of blocks to the super node address of a certain period
	GetTermWinnerMintCnt(term uint64, address hasharry.Address) (uint64, error)

	// Initialize dpos trie root
	InitTrie(contractRoot hasharry.Hash) error

	//Commit dpos trie information
	Commit() (hasharry.Hash, error)

	// Get root hash
	RootHash() hasharry.Hash

	// Close storage
	Close() error
}
