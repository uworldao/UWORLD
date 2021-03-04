package reqmgr

import (
	"errors"
	"github.com/libp2p/go-libp2p-core/network"
	"github.com/uworldao/UWORLD/common/encode/rlp"
	"github.com/uworldao/UWORLD/common/hasharry"
	"github.com/uworldao/UWORLD/core/types"
	"strconv"
)

var ErrorBlockNotFound = errors.New("block not exist")
var ErrorPeerClose = errors.New("peer is close")

const maxGetBlockCount = 30

type RWRequest struct {
	request *Request
	stream  network.Stream
}

func NewRWRequest(r *Request, stream network.Stream) *RWRequest {
	return &RWRequest{r, stream}
}

func (rm *RequestManager) getLastHeight(request *RWRequest) (*Response, error) {
	var message string
	var body []byte
	code := Success
	height := rm.blockChain.GetLastHeight()
	body, err := rlp.EncodeToBytes(height)
	if err != nil {
		code = DecodeError
		message = err.Error()
	}
	response := NewResponse(code, message, body)
	return response, nil

}

func (rm *RequestManager) getBlockByHeight(request *RWRequest) (*Response, error) {
	var message string
	var body []byte
	code := Success

	var height uint64 = 0
	err := rlp.DecodeBytes(request.request.Body, &height)
	if err != nil {
		code = DecodeError
		message = err.Error()
	} else if rm.blockChain.GetLastHeight() >= height {
		block, err := rm.blockChain.GetRlpBlockByHeight(height)
		if err != nil {
			code = InternalError
			message = err.Error()
		} else {
			body, _ = rlp.EncodeToBytes(block)
		}
	} else {
		code = InternalError
		message = ErrorBlockNotFound.Error()
	}

	response := NewResponse(code, message, body)
	return response, nil
}

func (rm *RequestManager) getBlocksByHeight(request *RWRequest) (*Response, error) {
	var message string
	var body []byte
	var height uint64 = 0
	var count uint64 = 0
	code := Success
	blocks := make(types.RlpBlocks, 0)

	err := rlp.DecodeBytes(request.request.Body, &height)
	if err != nil {
		code = DecodeError
		message = err.Error()
	} else if rm.blockChain.GetLastHeight() >= height {
		for rm.blockChain.GetLastHeight() >= height && count < maxGetBlockCount {
			block, err := rm.blockChain.GetRlpBlockByHeight(height)
			if err != nil {
				code = InternalError
				message = err.Error()
				response := NewResponse(code, message, body)
				return response, nil
			} else {
				blocks = append(blocks, block)
				height++
				count++
			}
		}
		body, _ = rlp.EncodeToBytes(blocks)
	} else {
		code = InternalError
		message = ErrorBlockNotFound.Error()
	}

	response := NewResponse(code, message, body)
	return response, nil
}

func (rm *RequestManager) getBlockByHash(request *RWRequest) (*Response, error) {
	var message string
	var body []byte
	code := Success

	block, err := rm.blockChain.GetRlpBlockByHash(hasharry.BytesToHash(request.request.Body))
	if err != nil {
		code = InternalError
		message = err.Error()
	} else {
		if block.Height > rm.blockChain.GetLastHeight() {
			code = InternalError
			message = "block not exist"
		} else {
			body, _ = rlp.EncodeToBytes(block)
		}
	}
	response := NewResponse(code, message, body)
	return response, nil
}

func (rm *RequestManager) getHeaderByHeight(request *RWRequest) (*Response, error) {
	var message string
	var body []byte
	code := Success

	height, _ := strconv.ParseUint(string(request.request.Body), 10, 64)

	header, err := rm.blockChain.GetHeaderByHeight(height)
	if err != nil {
		code = InternalError
		message = err.Error()
	} else {
		if header.Height > rm.blockChain.GetLastHeight() {
			code = InternalError
			message = "block not exist"
		} else {
			body, _ = rlp.EncodeToBytes(header)
		}
	}
	response := NewResponse(code, message, body)
	return response, nil
}

func (rm *RequestManager) receivedBlock(request *RWRequest) (*Response, error) {
	var block *types.RlpBlock
	var message string
	var body []byte
	code := Success
	err := rlp.DecodeBytes(request.request.Body, &block)
	if err != nil {
		code = InternalError
		message = "failed to encode"
	} else {
		rm.recBlkCh <- block.TranslateToBlock()
	}
	response := NewResponse(code, message, body)
	return response, nil
}

func (rm *RequestManager) receivedTransaction(request *RWRequest) (*Response, error) {
	var tx *types.RlpTransaction
	var message string
	var body []byte
	code := Success
	err := rlp.DecodeBytes(request.request.Body, &tx)
	if err != nil {
		code = InternalError
		message = "failed to decode"
	} else {
		rm.recTx <- tx.TranslateToTransaction()
	}
	response := NewResponse(code, message, body)
	return response, nil
}

func (rm *RequestManager) validationBlockHash(request *RWRequest) (*Response, error) {
	var message string
	var body []byte
	code := Success
	var header *types.Header
	err := rlp.DecodeBytes(request.request.Body, &header)
	if err != nil {
		code = EncodeError
		return NewResponse(code, message, body), nil
	}
	rs, err := rm.blockChain.ValidationBlockHash(header)
	if err != nil {
		code = InternalError
		return NewResponse(code, message, body), nil
	}
	body, _ = rlp.EncodeToBytes(rs)
	return NewResponse(code, message, body), nil
}

func (rm *RequestManager) getNodeInfo(request *RWRequest) (*Response, error) {
	var message string
	var body []byte
	code := Success
	nodeInfo := rm.peers.NodeInfo()
	body, err := rlp.EncodeToBytes(nodeInfo)
	if err != nil {
		code = InternalError
		message = err.Error()
	}
	return NewResponse(code, message, body), nil
}
