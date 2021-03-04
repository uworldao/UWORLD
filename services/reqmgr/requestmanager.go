package reqmgr

import (
	"encoding/json"
	"github.com/libp2p/go-libp2p-core/network"
	"github.com/uworldao/UWORLD/core"
	"github.com/uworldao/UWORLD/core/types"
	log "github.com/uworldao/UWORLD/log/log15"
	"sync"
	"time"
)

type requestfunc func(*RWRequest) (*Response, error)

type RequestManager struct {
	blockChain  core.IBlockChain
	requestChan chan *RWRequest
	recBlkCh    chan *types.Block
	recTx       chan types.ITransaction
	pool        sync.Pool
	peers       Peers
}

type Peers interface {
	PeersInfo() []*types.NodeInfo
	NodeInfo() *types.NodeInfo
}

func NewRequestManger(blockChain core.IBlockChain, recBlkCh chan *types.Block, recTx chan types.ITransaction, peers Peers) *RequestManager {
	return &RequestManager{
		blockChain:  blockChain,
		requestChan: make(chan *RWRequest, 100),
		recBlkCh:    recBlkCh,
		recTx:       recTx,
		pool: sync.Pool{
			New: func() interface{} {
				return make([]byte, maxReadBytes)
			},
		},
		peers: peers,
	}
}

// Listen for message requests
func (rm *RequestManager) Start() {
	for rwRequest := range rm.requestChan {
		var rf requestfunc
		switch rwRequest.request.Method {
		case getLastHeight:
			rf = rm.getLastHeight
		case getBlockByHeight:
			rf = rm.getBlockByHeight
		case getBlocksByHeight:
			rf = rm.getBlocksByHeight
		case getBlockByHash:
			rf = rm.getBlockByHash
		case getHeaderByHeight:
			rf = rm.getHeaderByHeight
		case getNodeInfo:
			rf = rm.getNodeInfo
		case sendBlock:
			rf = rm.receivedBlock
		case sendTransaction:
			rf = rm.receivedTransaction
		case validationBlockHash:
			rf = rm.validationBlockHash
		default:
			rwRequest.stream.Reset()
			rwRequest.stream.Close()
			continue
		}
		go dealRequest(rwRequest, rf)
	}
}

// Handling message requests
func dealRequest(rwRequest *RWRequest, rf requestfunc) {
	if response, err := rf(rwRequest); err != nil {
		log.Error("Handling request error", "method", rwRequest.request.Method, "error", err)
	} else if response != nil {
		if err := sendResponse(response, rwRequest.stream); err != nil {
			log.Error("Send response error", "method", rwRequest.request.Method, "peer", "error", err)
		}
	}
	for isConnect(rwRequest.stream) {
		time.Sleep(time.Second * 1)
	}
	rwRequest.stream.Reset()
	rwRequest.stream.Close()
}

func (rm *RequestManager) HandleRequest(stream network.Stream) {
	err := stream.SetDeadline(time.Unix(time.Now().Unix()+readTimeOut, 0))
	request, err := rm.ReadRequest(stream)
	if err != nil {
		return
	}
	rm.requestChan <- NewRWRequest(request, stream)
}

func sendResponse(response *Response, stream network.Stream) error {
	bytes, err := json.Marshal(response)
	if err != nil {
		return err
	}
	_, err = stream.Write(bytes)
	if err != nil {
		return err
	}
	return nil
}

func sendRequest(request *Request, stream network.Stream) error {
	bytes, err := json.Marshal(request)
	if err != nil {
		return err
	}
	_, err = stream.Write(bytes)
	if err != nil {
		return err
	}

	return nil
}

func isConnect(stream network.Stream) bool {
	bytes := [10]byte{}
	_, err := stream.Read(bytes[:])
	if err != nil {
		return false
	}
	return true
}
