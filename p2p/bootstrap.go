package p2p

import (
	log "github.com/jhdriver/UWORLD/log/log15"
	"github.com/libp2p/go-libp2p-core/peer"
	dht "github.com/libp2p/go-libp2p-kad-dht"
	"github.com/multiformats/go-multiaddr"
	"strings"
)

// Default boot node list
var DefaultBootstrapPeers []multiaddr.Multiaddr

// Custom boot node list
var CustomBootstrapPeers []multiaddr.Multiaddr

func init() {
	for _, s := range []string{
		/*"/ip4/103.15.132.180/tcp/11389/ipfs/16Uiu2HAmL2w4iB1i6PiRP8mBLbYdApkpWVktrUup22Fjm33usXSh",
		"/ip4/47.74.248.20/tcp/11360/ipfs/16Uiu2HAmHoiQRyTMbPvMrxueu1uF9wnTYZ6q3MEaPPwJMXrW1oQf",
		"/ip4/123.207.153.137/tcp/11369/ipfs/16Uiu2HAmKeXr5HNsEu2xaLyPPyHuxSeNQUHXCVgWCjieVGWddwWo",*/
		//"/ip4/47.57.100.253/tcp/2211/ipfs/16Uiu2HAkwQ1tmB5WrVgT83nD5KYZKanugXNVDM51vaTJw8TtxLa6",
		//"/ip4/8.210.1.245/tcp/2211/ipfs/16Uiu2HAm51UkR3V2V2zJt6GHnLjyjze7UxHYfGBDKZwKCCHdJ7c4",
		// UWDZPxLoUPGyYy3Y1Boeyyvk4LKy6n1NDT2n
		"/ip4/8.210.23.69/tcp/2211/ipfs/16Uiu2HAmSyvYS7YiANCzL2MkV3UEemABFqN6XNhH4k14ecSVkWmE",
		//"/ip4/127.0.0.1/tcp/2211/ipfs/16Uiu2HAmSyvYS7YiANCzL2MkV3UEemABFqN6XNhH4k14ecSVkWmE",
	} {
		ma, err := multiaddr.NewMultiaddr(s)
		if err != nil {
			panic(err)
		}
		DefaultBootstrapPeers = append(DefaultBootstrapPeers, ma)
	}
}

func IsInBootstrapPeers(id peer.ID) bool {
	bootstrap := DefaultBootstrapPeers
	if len(CustomBootstrapPeers) > 0 {
		bootstrap = CustomBootstrapPeers
	}
	for _, bootstrap := range bootstrap {
		if id.String() == strings.Split(bootstrap.String(), "/")[6] {
			return true
		}
	}
	return false
}

// Start boot node
func (p *P2pServer) StartBootStrap() error {
	var err error
	p.dht, err = dht.New(p.ctx, p.host)
	if err != nil {
		return err
	}
	log.Info("Bootstrapping the DHT")
	if err = p.dht.Bootstrap(p.ctx); err != nil {
		return err
	}
	return nil
}
