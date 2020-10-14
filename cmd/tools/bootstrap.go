package main

import (
	"errors"
	"flag"
	"fmt"
	"github.com/jhdriver/UWORLD/config"
	"github.com/jhdriver/UWORLD/p2p"
	"github.com/jhdriver/UWORLD/services/peermgr"
	p2pcrypto "github.com/libp2p/go-libp2p-core/crypto"
	"github.com/libp2p/go-libp2p-core/network"
	"github.com/libp2p/go-libp2p-core/peer"
	"os"
	"sync"
)

func main() {
	var (
		port     = flag.String("port", "2211", "the port of start a bootstrap")
		keyFile  = flag.String("k", "", "bootstrap node key file")
		password = flag.String("p", "", "the decryption password for key file")
	)
	flag.Parse()
	StartBootStrap(*port, *keyFile, *password)
}

func newPeerFunc(info *p2p.PeerInfo) {

}

func removePeerFunc(id string) {

}

func hasPeerFunc(id string) bool {
	return false
}

func requestHanle(stream network.Stream) {

}

func StartBootStrap(port, keyFile, password string) {

	wg := sync.WaitGroup{}
	wg.Add(1)

	if keyFile == "" {
		flag.PrintDefaults()
		return
	}
	if password == "" {
		fmt.Println("please enter the password for the keyfile:")
		passWd, err := readPassWd()
		if err != nil {
			fmt.Printf("read password failed! %s\n", err.Error())
			return
		}
		password = string(passWd)
	}
	nodePrivate, err := config.LoadNodePrivate(keyFile, password)
	if err != nil {
		fmt.Printf("failed to load keyfile %s! %s\n", keyFile, err.Error())
		return
	}

	p2pPriKey, err := p2pcrypto.UnmarshalSecp256k1PrivateKey(nodePrivate.PrivateKey.Serialize())
	if err != nil {
		fmt.Printf("generate priavte failed! %v\n", err)
		return
	}
	localPeer := &p2p.PeerInfo{PrivateKey: p2pPriKey, AddrInfo: &peer.AddrInfo{}}
	server, err := p2p.NewBootStrapP2pServer(port, localPeer, peermgr.NewPeerManager(localPeer), nil)
	if err != nil {
		fmt.Printf("create p2p server failed! %v\n", err)
		return
	}

	if err := server.StartBootStrap(); err != nil {
		fmt.Printf("start p2p server failed! %v\n", err)
		return
	}
	wg.Wait()
}

func readPassWd() ([]byte, error) {
	var passWd [33]byte

	n, err := os.Stdin.Read(passWd[:])
	if err != nil {
		return nil, err
	}
	if n <= 1 {
		return nil, errors.New("not read")
	}
	return passWd[:n-1], nil
}
