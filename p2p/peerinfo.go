package p2p

import (
	"github.com/libp2p/go-libp2p-core/crypto"
	"github.com/libp2p/go-libp2p-core/network"
	"github.com/libp2p/go-libp2p-core/peer"
)

const PeerInfoFile = "localpeer"

type NewStreamFunction func(peerId peer.ID) (network.Stream, error)

type StreamCreator struct {
	Stream        network.Stream
	PeerId        peer.ID
	NewStreamFunc NewStreamFunction
}

func (s *StreamCreator) Close() {
	s.Stream.Reset()
	s.Stream.Close()
	s.Stream = nil
}

type PeerInfo struct {
	PrivateKey crypto.PrivKey
	AddrInfo   *peer.AddrInfo
	*StreamCreator
}

func NewPeerInfo(private crypto.PrivKey, addrInfo *peer.AddrInfo, newStreamFunc NewStreamFunction) *PeerInfo {
	return &PeerInfo{PrivateKey: private, AddrInfo: addrInfo, StreamCreator: &StreamCreator{PeerId: addrInfo.ID, NewStreamFunc: newStreamFunc}}
}

func StringToPeerID(id string) (peer.ID, error) {
	peerId := new(peer.ID)
	return *peerId, peerId.UnmarshalText([]byte(id))
}
