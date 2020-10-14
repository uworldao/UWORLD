package blkmgr

import (
	"errors"
	"github.com/jhdriver/UWORLD/consensus"
	"github.com/jhdriver/UWORLD/core"
	"github.com/jhdriver/UWORLD/core/types"
	log "github.com/jhdriver/UWORLD/log/log15"
	"github.com/jhdriver/UWORLD/p2p"
	"github.com/jhdriver/UWORLD/services/reqmgr"
	"github.com/libp2p/go-libp2p-core/network"
	"github.com/libp2p/go-libp2p-core/peer"
	"math/rand"
	"time"
)

const getPeerInterval = 1
const syncInterval = 1000

// Block management, synchronization and sending new blocks
type BlockManager struct {
	blockChain  core.IBlockChain
	peerManager p2p.IPeerManager
	syncPeer    *p2p.PeerInfo
	network     Network
	consensus   consensus.IConsensus
	newStream   ICreateStream
	revBlkCh    chan *types.Block
	genBlkCh    chan *types.Block
	minerWokCh  chan bool
	needHash    []byte
	quitCh      chan bool
	isQuit      chan bool
}

type ICreateStream interface {
	CreateStream(peerId peer.ID) (network.Stream, error)
}

func NewBlockManager(blockChain core.IBlockChain, peerManager p2p.IPeerManager, network Network, consensus consensus.IConsensus,
	revBlkCh chan *types.Block, genBlkCh chan *types.Block, minerWokCh chan bool, createStream ICreateStream) *BlockManager {
	return &BlockManager{
		blockChain:  blockChain,
		peerManager: peerManager,
		network:     network,
		consensus:   consensus,
		newStream:   createStream,
		revBlkCh:    revBlkCh,
		genBlkCh:    genBlkCh,
		minerWokCh:  minerWokCh,
		quitCh:      make(chan bool, 1),
		isQuit:      make(chan bool, 1),
	}
}

func (bm *BlockManager) Start() error {
	go bm.network.Start()

	go bm.handleBlock()

	go bm.syncBlock()

	log.Info("Block manager startup successful")
	return nil
}

func (bm *BlockManager) Stop() error {
	close(bm.quitCh)
	<-bm.isQuit
	log.Info("Stop block manager")
	return nil
}

// Start sync block
func (bm *BlockManager) syncBlock() {
	for {
		select {
		case _, _ = <-bm.quitCh:
			log.Info("Sync block quit")
			bm.isQuit <- true
			return
		default:
			bm.createSyncStream()
			bm.syncBlockFromStream()
		}
		time.Sleep(time.Millisecond * syncInterval)
	}
}

// Create a network channel of the synchronization block, and randomly
// select a new peer node for synchronization every 1s.
func (bm *BlockManager) createSyncStream() {
	for {
		select {
		case _, _ = <-bm.quitCh:
			log.Info("Create sync stream quit")
			return
		default:
			bm.findSyncPeer()

			if err := bm.createStream(); err != nil {
				/*log.Debug("The current peer is closed，choose again", "closed peer", bm.syncPeer.AddrInfo)*/
				if bm.syncPeer != nil {
					bm.peerManager.Remove(bm.syncPeer.AddrInfo.ID.String())
				}
			} else {
				bm.syncPeer.Close()
			}
			return
		}
	}
}

// Replace the new peer node
func (bm *BlockManager) findSyncPeer() {
	t := time.NewTicker(time.Second * getPeerInterval)
	defer t.Stop()

	for {
		select {
		case _, _ = <-bm.quitCh:
			log.Info("Find sync peer quit")
			return
		case _ = <-t.C:
			bm.syncPeer = bm.peerManager.GetPeer()
			if bm.syncPeer == nil {
				//log.Warn("No available peers were found， wait...")
			} else {
				//log.Info("Find an available peer", "peer", bm.syncPeer)
				return
			}
		}
	}
}

// Create node communication stream
func (bm *BlockManager) createStream() error {
	if bm.syncPeer == nil {
		return errors.New("sync peer is nil")
	}
	stream, err := bm.syncPeer.NewStreamFunc(bm.syncPeer.PeerId)
	if err != nil {
		return err
	}
	bm.syncPeer.StreamCreator.Stream = stream
	return nil
}

// Synchronize blocks from the stream and verify storage
func (bm *BlockManager) syncBlockFromStream() error {
	for {
		select {
		case _, _ = <-bm.quitCh:
			log.Info("Sync block from stream quit")
			return nil
		default:
			localHeight := bm.blockChain.GetLastHeight()

			// Get the block of the remote node from the next block height，
			// If the error is that the peer has stopped, delete the peer.
			// If the storage fails locally, the remote block verification
			// is performed, the verification proves that the local block
			// is wrong, and the local chain is rolled back to the valid block.
			blocks, err := bm.network.GetBlocksByHeight(bm.syncPeer.StreamCreator, localHeight+1)
			if err != nil {
				if err == reqmgr.ErrorPeerClose {
					bm.peerManager.Remove(bm.syncPeer.AddrInfo.ID.String())
				}
				return err
			}
			/*if err := bm.blockChain.InsertChain(block); err != nil {
				log.Warn("Insert chain failed!", "error", err, "height", block.Height)
				if bm.remoteValidation(block.Header) {
					bm.fallBack()
					return err
				}
				bm.peerManager.Remove(bm.syncPeer.AddrInfo.ID.String())
				return err
			}*/
			if err := bm.insertBlocksToChain(blocks); err != nil {
				return err
			}
		}
	}
	return nil
}

func (bm *BlockManager) insertBlocksToChain(blocks []*types.Block) error {
	for _, block := range blocks {
		select {
		case _, _ = <-bm.quitCh:
			log.Info("Insert blocks quit")
			return nil
		default:
			if err := bm.blockChain.InsertChain(block); err != nil {
				log.Warn("Insert chain failed!", "error", err, "height", block.Height)
				if bm.remoteValidation(block.Header) {
					bm.fallBack()
					return err
				} else {
					localPreHeader, _ := bm.blockChain.GetHeaderByHeight(block.Height - 1)
					if !bm.remoteValidation(localPreHeader) {
						bm.fallBack()
						return err
					}
				}
				bm.peerManager.Remove(bm.syncPeer.AddrInfo.ID.String())
				return err
			}
		}
	}
	return nil
}

func (bm *BlockManager) remoteValidation(header *types.Header) bool {
	var hashCount = 0
	if header.Height <= bm.consensus.GetConfirmedBlockHeader(bm.blockChain).Height {
		return false
	}
	ids, err := bm.consensus.GetWinnersPeerID(header.Time)
	if err != nil {
		return false
	}
	localHeader, err := bm.blockChain.GetHeaderByHeight(header.Height)
	if err == nil && localHeader.Hash.IsEqual(header.Hash) {
		hashCount = 1
	}

	for _, id := range ids {
		if id != bm.peerManager.LocalPeerInfo().AddrInfo.ID.String() {
			peerId := new(peer.ID)
			if err = peerId.UnmarshalText([]byte(id)); err == nil {
				streamCreator := p2p.StreamCreator{PeerId: *peerId, NewStreamFunc: bm.newStream.CreateStream}
				ok, err := bm.network.ValidationBlockHash(&streamCreator, header)
				if err == nil && ok {
					hashCount++
				}
			}
		}
	}
	if hashCount > len(ids)/2 {
		return true
	}
	return false
}

// Remotely verify the block, if the block height is less than
// the effective block height, then discard the block. If the
// block occupies the majority of the currently started super
// nodes, it means that the block is more likely to be correct,
// and the block verification is successful.
func (bm *BlockManager) remoteValidation1(header *types.Header) bool {
	if header.Height <= bm.consensus.GetConfirmedBlockHeader(bm.blockChain).Height {
		return false
	}
	ids, err := bm.consensus.GetWinnersPeerID(header.Time)
	if err != nil {
		return false
	}
	localHeader, err := bm.blockChain.GetHeaderByHeight(header.Height)
	if err == nil {
		if localHeader.HashString() == header.HashString() {
			return false
		}
	}
	compareMap := make(map[string][]string)
	for _, id := range ids {
		peerId := new(peer.ID)
		if id != bm.peerManager.LocalPeerInfo().AddrInfo.ID.String() {
			if err = peerId.UnmarshalText([]byte(id)); err == nil {
				streamCreator := p2p.StreamCreator{PeerId: *peerId, NewStreamFunc: bm.newStream.CreateStream}
				remoteHeader, err := bm.network.GetHeaderByHeight(&streamCreator, header.Height)
				if err != nil {
					continue
				}
				if _, ok := compareMap[remoteHeader.HashString()]; ok {
					compareMap[remoteHeader.HashString()] = append(compareMap[remoteHeader.HashString()], id)
				} else {
					compareMap[remoteHeader.HashString()] = []string{id}
				}
			}
		} else {
			localHeader, err := bm.blockChain.GetHeaderByHeight(header.Height)
			if err != nil {
				return true
			}
			if _, ok := compareMap[localHeader.HashString()]; ok {
				compareMap[localHeader.HashString()] = append(compareMap[localHeader.HashString()], id)
			} else {
				compareMap[localHeader.HashString()] = []string{id}
			}
		}
	}
	selectedHash := getMaxCountHash(compareMap)
	if header.HashString() != selectedHash {
		return false
	}
	return true
}

// Block chain rolls back to a valid block
func (bm *BlockManager) fallBack() {
	bm.blockChain.FallBack()
}

// Broadcast the block generated by yourself to the super node
func (bm *BlockManager) broadCastBlock(block *types.Block) {
	ids, err := bm.consensus.GetWinnersPeerID(block.Time)
	if err != nil {
		return
	}
	for _, id := range ids {
		if id != bm.peerManager.LocalPeerInfo().AddrInfo.ID.String() {
			peerId := new(peer.ID)
			if err = peerId.UnmarshalText([]byte(id)); err == nil {
				streamCreator := p2p.StreamCreator{PeerId: *peerId, NewStreamFunc: bm.newStream.CreateStream}
				if err := bm.network.SendBlock(&streamCreator, block); err != nil {
					log.Warn("Failed to send block", "height", block.Height, "target", id, "error", err)
				}
			}
		}
	}
}

// Send and receive blocks
func (bm *BlockManager) handleBlock() {
	for {
		select {
		case _, _ = <-bm.quitCh:
			log.Info("Handle block quit")
			return
		case block := <-bm.genBlkCh:
			go bm.broadCastBlock(block)
		case block := <-bm.revBlkCh:
			go bm.dealReceivedBlock(block)
		}
	}
}

// Process blocks received from other super nodes.If the height
// of the block is greater than the local height, the storage is
// directly verified. If the height is less than the local height,
// the remote verification is performed, and the verification is
// passed back to the local block.
func (bm *BlockManager) dealReceivedBlock(block *types.Block) {
	localHeight := bm.blockChain.GetLastHeight()
	if localHeight == block.Height-1 {
		if err := bm.blockChain.InsertChain(block); err != nil {
			log.Warn("Failed to insert received block", "err", err, "height", block.Height, "singer", block.Signer.String())
		}
	} else if block.Height <= localHeight {
		if localHeader, err := bm.blockChain.GetHeaderByHeight(block.Height); err == nil {
			if !localHeader.Hash.IsEqual(block.Hash) {
				if bm.remoteValidation(block.Header) {
					bm.fallBack()
				} else {
					if !bm.remoteValidation(localHeader) {
						bm.fallBack()
					} else {
						log.Warn("Remote validation failed!", "height", block.Height, "signer", block.Signer.String())
					}

				}
			}
		} else {
			if err := bm.blockChain.InsertChain(block); err != nil {
				log.Warn("Failed to insert received block", "err", err, "height", block.Height, "singer", block.Signer.String())
			}
		}
	}
}

func getMaxCountHash(compareMap map[string][]string) string {
	hashes := make([]string, 0)
	var maxCount int
	for h, peers := range compareMap {
		if len(peers) == maxCount {
			hashes = append(hashes, h)
		} else if len(peers) > maxCount {
			maxCount = len(peers)
			hashes = []string{h}
		}
	}
	if len(hashes) > 1 {
		rand.Intn(len(hashes))
		return hashes[rand.Intn(len(hashes))]
	}
	if len(hashes) == 0 {
		return ""
	}
	return hashes[0]
}
