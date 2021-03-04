package peermgr

import (
	log "github.com/uworldao/UWORLD/log/log15"
	"github.com/uworldao/UWORLD/p2p"
	"math/rand"
	"sync"
	"time"
)

const maxPeers = 1000000
const checkPeerInterval = 60 * 30

// For the implementation of peer node management, check the
// status of the node every 10s, and delete the peer node if
// the node is shut down.
type PeerManager struct {
	localNode *p2p.PeerInfo
	peersMap  map[string]*p2p.PeerInfo
	removeMap map[string]*p2p.PeerInfo
	mutex     sync.RWMutex
	userId    string
	peerIds   []string
}

func NewPeerManager(localNode *p2p.PeerInfo) *PeerManager {
	return &PeerManager{
		localNode: localNode,
		peersMap:  make(map[string]*p2p.PeerInfo),
		removeMap: make(map[string]*p2p.PeerInfo),
	}
}

func (pm *PeerManager) HashPeer(id string) bool {
	pm.mutex.RLock()
	defer pm.mutex.RUnlock()

	if _, ok := pm.peersMap[id]; !ok {
		return false
	}
	return true
}

func (pm *PeerManager) Add(info *p2p.PeerInfo) {
	pm.mutex.Lock()
	defer pm.mutex.Unlock()

	if len(pm.peersMap) >= maxPeers {
		return
	}
	log.Info("New Peer", "address", info.AddrInfo.String())
	pm.peersMap[info.AddrInfo.ID.String()] = info
	pm.peerIds = append(pm.peerIds, info.AddrInfo.ID.String())
}

func (pm *PeerManager) Remove(reId string) {
	pm.mutex.Lock()
	defer pm.mutex.Unlock()

	for index, id := range pm.peerIds {
		if id == reId {
			pm.peerIds = append(pm.peerIds[0:index], pm.peerIds[index+1:]...)
			if peer, ok := pm.peersMap[reId]; ok {
				delete(pm.peersMap, reId)
				pm.removeMap[reId] = peer
				log.Info("Delete peer", "peer", id)
			}
			break
		}
	}
}

func (pm *PeerManager) Check() {
	t := time.NewTicker(time.Second * checkPeerInterval)
	defer t.Stop()

	for range t.C {
		for id, peer := range pm.peersMap {
			if id != pm.userId {
				if !pm.isPeerLive(peer) {
					pm.Remove(id)
				}
			}
		}
		for id, peer := range pm.removeMap {
			if id != pm.userId {
				if pm.isPeerLive(peer) {
					pm.Add(peer)
				}
			}
		}
	}
}

func (pm *PeerManager) isPeerLive(peer *p2p.PeerInfo) bool {
	stream, err := peer.NewStreamFunc(peer.AddrInfo.ID)
	if err != nil {
		return false
	}
	stream.Reset()
	stream.Close()
	return true
}

func (pm *PeerManager) GetPeer() *p2p.PeerInfo {
	pm.mutex.Lock()
	defer pm.mutex.Unlock()

	if len(pm.peersMap) == 0 {
		return nil
	}
	index := rand.New(rand.NewSource(time.Now().Unix())).Int31n(int32(len(pm.peerIds)))
	peerId := pm.peerIds[index]
	return pm.peersMap[peerId]
}

func (pm *PeerManager) LocalPeerInfo() *p2p.PeerInfo {
	return pm.localNode
}

func (pm *PeerManager) SetLocalPeerInfo(local *p2p.PeerInfo) {
	pm.localNode = local
}

func (pm *PeerManager) Peers() map[string]*p2p.PeerInfo {
	pm.mutex.RLock()
	defer pm.mutex.RUnlock()

	re := make(map[string]*p2p.PeerInfo)
	for key, value := range pm.peersMap {
		re[key] = value
	}
	return re
}

func (pm *PeerManager) Count() uint32 {
	pm.mutex.RLock()
	defer pm.mutex.RUnlock()

	count := uint32(len(pm.peersMap))
	return count
}
