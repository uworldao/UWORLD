package blkmgr

import (
	"github.com/libp2p/go-libp2p-core/network"
	"github.com/uworldao/UWORLD/common/hasharry"
	"github.com/uworldao/UWORLD/core/types"
	"github.com/uworldao/UWORLD/p2p"
)

// Peer node communication interface
type Network interface {
	// Communication monitoring start
	Start()

	// Handling peer node requests
	HandleRequest(stream network.Stream)

	// Get peer height
	GetLastHeight(stream *p2p.StreamCreator) (uint64, error)

	// Get a block of a certain height of the peer node
	GetBlockByHeight(stream *p2p.StreamCreator, height uint64) (*types.Block, error)

	// Get multiple blocks after a certain height of the peer node
	GetBlocksByHeight(stream *p2p.StreamCreator, height uint64) ([]*types.Block, error)

	// Obtain the peer block through hash
	GetBlockByHash(stream *p2p.StreamCreator, hash hasharry.Hash) (*types.Block, error)

	// Obtain the block header of the peer node by height
	GetHeaderByHeight(stream *p2p.StreamCreator, height uint64) (*types.Header, error)

	// Send blocks to peer nodes
	SendBlock(stream *p2p.StreamCreator, block *types.Block) error

	// Send transactions to peer nodes
	SendTransaction(stream *p2p.StreamCreator, tx types.ITransaction) error

	// Remotely verify whether a block is consistent
	ValidationBlockHash(stream *p2p.StreamCreator, header *types.Header) (bool, error)

	// Get peer information
	GetNodeInfo(stream *p2p.StreamCreator) (*types.NodeInfo, error)
}
