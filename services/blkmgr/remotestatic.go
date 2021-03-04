package blkmgr

import (
	"github.com/libp2p/go-libp2p-core/peer"
	"github.com/uworldao/UWORLD/core/types"
)

type remoteBlockStatic struct {
	hash  string
	count int
	Ids   []peer.ID
}

type staticMap struct {
	reBlkStaticMap map[string]*remoteBlockStatic
}

func newStaticMap() *staticMap {
	return &staticMap{make(map[string]*remoteBlockStatic)}
}

func (sm *staticMap) Add(remoteBlock *types.Block, peerId peer.ID) {
	hash := remoteBlock.HashString()
	if rb, ok := sm.reBlkStaticMap[hash]; ok {
		rb.count++
		rb.Ids = append(rb.Ids, peerId)
	} else {
		rb := &remoteBlockStatic{hash, 1, []peer.ID{peerId}}
		sm.reBlkStaticMap[hash] = rb
	}
}

func (sm *staticMap) Len() int {
	return len(sm.reBlkStaticMap)
}

func (sm *staticMap) FindMaxCountStatics() []*remoteBlockStatic {
	var maxStatics []*remoteBlockStatic
	maxCount := 0
	for _, rb := range sm.reBlkStaticMap {
		if maxCount <= rb.count {
			maxCount = rb.count
		}
	}
	for _, rb := range sm.reBlkStaticMap {
		if rb.count == maxCount {
			maxStatics = append(maxStatics, rb)
		}
	}
	return maxStatics
}
