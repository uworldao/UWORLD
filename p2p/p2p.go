package p2p

import (
	"context"
	"fmt"
	"github.com/jhdriver/UWORLD/config"
	"github.com/jhdriver/UWORLD/crypto/ecc/secp256k1"
	log "github.com/jhdriver/UWORLD/log/log15"
	"github.com/jhdriver/UWORLD/param"
	"github.com/libp2p/go-libp2p"
	core "github.com/libp2p/go-libp2p-core"
	"github.com/libp2p/go-libp2p-core/crypto"
	"github.com/libp2p/go-libp2p-core/network"
	"github.com/libp2p/go-libp2p-core/peer"
	"github.com/libp2p/go-libp2p-core/protocol"
	discovery "github.com/libp2p/go-libp2p-discovery"
	dht "github.com/libp2p/go-libp2p-kad-dht"
	"github.com/multiformats/go-multiaddr"
	"net"
	"sync"
	"time"
)

type P2pServer struct {
	host         core.Host
	dht          *dht.IpfsDHT
	ctx          context.Context
	localPeer    *PeerInfo
	peerManager  IPeerManager
	StreamHandle StreamHandle
	config       *config.Config
	stop         chan bool
	hasStop      chan bool
}

// Manage peer nodes
type IPeerManager interface {
	SetLocalPeerInfo(*PeerInfo)
	LocalPeerInfo() *PeerInfo
	GetPeer() *PeerInfo
	HashPeer(id string) bool
	Peers() map[string]*PeerInfo
	Remove(id string)
	Count() uint32
	Add(info *PeerInfo)
	Check()
}

type StreamHandle interface {
	HandleRequest(network.Stream)
}

func newHost(localPeer *PeerInfo, config *config.Config) (core.Host, error) {
	addressFactory := func(addrs []multiaddr.Multiaddr) []multiaddr.Multiaddr { return []multiaddr.Multiaddr{} }
	ips := GetLocalIp()
	external, _ := multiaddr.NewMultiaddr(fmt.Sprintf("/ip4/%s/tcp/%s", config.ExternalIp, config.P2pPort))
	addressFactory = func(addrs []multiaddr.Multiaddr) []multiaddr.Multiaddr {
		addrs = []multiaddr.Multiaddr{}
		if external != nil {
			addrs = append(addrs, external)
		}
		for _, ip := range ips {
			extMultiAddr, _ := multiaddr.NewMultiaddr(fmt.Sprintf("/ip4/%s/tcp/%s", ip, config.P2pPort))
			if extMultiAddr != nil {
				addrs = append(addrs, extMultiAddr)
			}
		}
		return addrs
	}

	opts := []libp2p.Option{
		libp2p.ListenAddrStrings(fmt.Sprintf("/ip4/0.0.0.0/tcp/%s", config.P2pPort)),
		libp2p.Identity(localPeer.PrivateKey),
		libp2p.DefaultMuxers,
		libp2p.EnableRelay(),
		libp2p.AddrsFactory(addressFactory),
	}
	return libp2p.New(context.Background(), opts...)
}

func NewP2pServer(config *config.Config, localPeer *PeerInfo, peerManager IPeerManager, sh StreamHandle) (*P2pServer, error) {

	var err error
	p2pServer := &P2pServer{
		localPeer:    localPeer,
		peerManager:  peerManager,
		StreamHandle: sh,
		ctx:          context.Background(),
		config:       config,
		stop:         make(chan bool, 1),
		hasStop:      make(chan bool, 1),
	}
	if config.Bootstrap != "" {
		ma, err := multiaddr.NewMultiaddr(config.Bootstrap)
		if err != nil {
			return nil, fmt.Errorf("wrong bootstrap node address; %s", err)
		}
		CustomBootstrapPeers = append(CustomBootstrapPeers, ma)
	}

	host, err := newHost(localPeer, config)
	if err != nil {
		return nil, err
	}
	p2pServer.host = host
	localPeer.AddrInfo.ID = p2pServer.host.ID()
	peerManager.SetLocalPeerInfo(localPeer)
	log.Info("Host created", "id", p2pServer.host.ID(), "address", p2pServer.host.Addrs())
	return p2pServer, nil
}

func NewBootStrapP2pServer(p2pPort string, localPeer *PeerInfo, peerManager IPeerManager, sh StreamHandle) (*P2pServer, error) {

	var err error
	p2pServer := &P2pServer{
		localPeer:    localPeer,
		peerManager:  peerManager,
		StreamHandle: sh,
		ctx:          context.Background(),
		stop:         make(chan bool, 1),
		hasStop:      make(chan bool, 1),
	}
	sourceMultiAddr, _ := multiaddr.NewMultiaddr(fmt.Sprintf("/ip4/%s/tcp/%s", "0.0.0.0", p2pPort))
	opts := []libp2p.Option{
		libp2p.ListenAddrs(sourceMultiAddr),
		libp2p.Identity(localPeer.PrivateKey),
		libp2p.EnableRelay(),
	}

	host, err := libp2p.New(p2pServer.ctx, opts...)
	if err != nil {
		return nil, err
	}
	p2pServer.host = host
	localPeer.AddrInfo.ID = p2pServer.host.ID()
	log.Info("Host created", "id", p2pServer.host.ID(), "address", p2pServer.host.Addrs())
	return p2pServer, nil
}

func (p *P2pServer) Start() error {
	p.setP2pHandleStream()

	if err := p.connectBootstrap(); err != nil {
		log.Error("Connect bootstrap failed!", "error", err)
		return err
	}

	go p.discovery()
	log.Info("P2p startup successful")
	return nil
}

func (p *P2pServer) Restart() error {
	if err := p.Stop(); err != nil {
		return err
	}

	var err error
	p.host, err = newHost(p.localPeer, p.config)
	if err != nil {
		return err
	}
	return p.Start()
}

func (p *P2pServer) Stop() error {
	p.stop <- true
	if err := p.host.Close(); err != nil {
		return err
	}
	log.Info("Stop P2P server")
	return nil
}

func (p *P2pServer) Addr() string {
	addrs := p.host.Addrs()
	var rs string
	for _, addr := range addrs {
		rs += "[" + addr.String() + "]"
	}
	return rs
}

func (p *P2pServer) ID() string {
	return p.host.ID().String()
}

func (p *P2pServer) CreateStream(peerId peer.ID) (network.Stream, error) {
	return p.host.NewStream(context.Background(), peerId, protocol.ID(param.UniqueNetWork))
}

func (p *P2pServer) setP2pHandleStream() {
	p.host.SetStreamHandler(protocol.ID(param.UniqueNetWork), p.StreamHandle.HandleRequest)
}

func (p *P2pServer) connectBootstrap() error {
	var err error
	p.dht, err = dht.New(p.ctx, p.host)
	if err != nil {
		return err
	}

	log.Info("Bootstrapping the DHT")
	if err = p.dht.Bootstrap(p.ctx); err != nil {
		return err
	}

	// Preferably use a custom boot node
	bootstrap := DefaultBootstrapPeers
	if len(CustomBootstrapPeers) > 0 {
		bootstrap = CustomBootstrapPeers
	}
	var wg sync.WaitGroup
	//
	for _, peerAddr := range bootstrap {
		peerInfo, _ := peer.AddrInfoFromP2pAddr(peerAddr)
		wg.Add(1)
		go func() {
			defer wg.Done()
			if err := p.host.Connect(p.ctx, *peerInfo); err != nil {
				log.Warn("Connection established with bootstrap node", "error", err)
			} else {
				p.peerIsLive(peerInfo.ID)
				log.Info("Connection established with bootstrap node", "peer", *peerInfo)
			}
		}()
	}
	wg.Wait()
	return nil
}

// Rediscover new nodes every 10s
func (p *P2pServer) discovery() error {
	log.Info("Announcing ourselves...")
	routingDiscovery := discovery.NewRoutingDiscovery(p.dht)
	discovery.Advertise(p.ctx, routingDiscovery, param.UniqueNetWork)

	for {
		select {
		case _ = <-p.stop:
			p.hasStop <- true
			return nil
		default:
			log.Info("Searching for other peers...")
			peerChan, err := routingDiscovery.FindPeers(p.ctx, param.UniqueNetWork)
			if err != nil {
				log.Error("Failed to find peers", "error", err)
				time.Sleep(time.Second * 60)
				continue
			}
		OUTCHAN:
			for {
				select {
				case addr, ok := <-peerChan:
					if ok {
						if addr.ID == p.host.ID() || IsInBootstrapPeers(addr.ID) {
							continue
						}
						if !p.peerManager.HashPeer(addr.ID.String()) {
							//log.Info("Connecting to:", "peer", addr.String())
							if !p.peerIsLive(addr.ID) {
								//log.Warn("New stream failed!", "addr", addr.String())
								p.peerManager.Remove(addr.ID.String())
								continue
							}
							peerInfo := NewPeerInfo(nil, copyPeerAddr(&addr), p.CreateStream)

							p.peerManager.Add(peerInfo)
						}
					} else {
						break OUTCHAN
					}
				}
			}
		}
		time.Sleep(time.Second * 60)
	}
	return nil
}

// Determine whether the peer is alive
func (p *P2pServer) peerIsLive(id peer.ID) bool {
	stream, err := p.host.NewStream(context.Background(), id, protocol.ID(param.UniqueNetWork))
	if err != nil {
		return false
	}
	stream.Reset()
	stream.Close()

	return true
}

func GenerateP2pId(key *secp256k1.PrivateKey) (peer.ID, error) {
	p2pPriKey, err := crypto.UnmarshalSecp256k1PrivateKey(key.Serialize())
	if err != nil {
		return "", err
	}
	opts := []libp2p.Option{
		libp2p.ListenAddrStrings(fmt.Sprintf("/ip4/0.0.0.0/tcp/%s", "65535")),
		libp2p.Identity(p2pPriKey),
	}
	host, err := libp2p.New(context.Background(), opts...)
	if err != nil {
		return "", err
	}
	defer host.Close()
	return host.ID(), nil
}

func GetLocalIp() []string {
	ips := []string{}
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return ips
	}
	for _, address := range addrs {
		if ipnet, flag := address.(*net.IPNet); flag && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				ips = append(ips, ipnet.IP.String())
			}
		}
	}
	return ips
}

func copyPeerAddr(addr *peer.AddrInfo) *peer.AddrInfo {
	newAddrs := make([]multiaddr.Multiaddr, len(addr.Addrs))
	for i, maddr := range addr.Addrs {
		newAddrs[i] = maddr
	}
	return &peer.AddrInfo{ID: addr.ID, Addrs: newAddrs}
}
