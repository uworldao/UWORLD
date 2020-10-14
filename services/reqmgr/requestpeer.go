package reqmgr

import (
	"errors"
	"fmt"
	"github.com/jhdriver/UWORLD/common/encode/rlp"
	"github.com/jhdriver/UWORLD/common/hasharry"
	"github.com/jhdriver/UWORLD/core/types"
	"github.com/jhdriver/UWORLD/p2p"
	"strconv"
	"time"
)

func (rm *RequestManager) GetLastHeight(stream *p2p.StreamCreator) (uint64, error) {
	var height uint64 = 0
	s, err := stream.NewStreamFunc(stream.PeerId)
	if err != nil {
		return 0, err
	}
	defer func() {
		s.Reset()
		s.Close()
	}()

	s.SetDeadline(time.Unix(time.Now().Unix()+readTimeOut, 0))
	request := NewRequest(getLastHeight, nil)
	err = sendRequest(request, s)
	if err != nil {
		return 0, err
	}
	response, err := rm.ReadResponse(s)
	if response != nil && response.Code == Success {
		err := rlp.DecodeBytes(response.Body, &height)
		if err != nil {
			return 0, err
		}
	} else {
		return 0, fmt.Errorf("peer error: %v", err)
	}
	return height, nil
}

func (rm *RequestManager) GetBlockByHeight(stream *p2p.StreamCreator, height uint64) (*types.Block, error) {
	var block *types.RlpBlock
	bytes, err := rlp.EncodeToBytes(height)
	if err != nil {
		return nil, err
	}
	s, err := stream.NewStreamFunc(stream.PeerId)
	if err != nil {
		return nil, err
	}
	defer func() {
		s.Reset()
		s.Close()
	}()

	s.SetDeadline(time.Unix(time.Now().Unix()+readTimeOut, 0))
	request := NewRequest(getBlockByHeight, bytes)
	err = sendRequest(request, s)
	if err != nil {
		return nil, ErrorPeerClose
	}
	response, err := rm.ReadResponse(s)
	if response != nil && response.Code == Success {
		err := rlp.DecodeBytes(response.Body, &block)
		if err != nil {
			return nil, err
		}
	} else if response != nil && response.Message == ErrorBlockNotFound.Error() {
		return nil, ErrorBlockNotFound
	} else {
		return nil, ErrorPeerClose
	}
	return block.TranslateToBlock(), nil
}

func (rm *RequestManager) GetBlocksByHeight(stream *p2p.StreamCreator, height uint64) ([]*types.Block, error) {
	blocks := types.RlpBlocks{}
	bytes, err := rlp.EncodeToBytes(height)
	if err != nil {
		return nil, err
	}
	s, err := stream.NewStreamFunc(stream.PeerId)
	if err != nil {
		return nil, err
	}
	defer func() {
		s.Reset()
		s.Close()
	}()

	s.SetDeadline(time.Unix(time.Now().Unix()+60, 0))
	request := NewRequest(getBlocksByHeight, bytes)
	err = sendRequest(request, s)
	if err != nil {
		return nil, ErrorPeerClose
	}
	response, err := rm.ReadResponse(s)
	if response != nil && response.Code == Success {
		err := rlp.DecodeBytes(response.Body, &blocks)
		if err != nil {
			return nil, err
		}
	} else if response != nil && response.Message == ErrorBlockNotFound.Error() {
		return nil, ErrorBlockNotFound
	} else {
		return nil, ErrorPeerClose
	}
	return blocks.TranslateToBlocks(), nil
}

func (rm *RequestManager) GetBlockByHash(stream *p2p.StreamCreator, hash hasharry.Hash) (*types.Block, error) {
	var block *types.RlpBlock
	s, err := stream.NewStreamFunc(stream.PeerId)
	if err != nil {
		return nil, err
	}
	defer func() {
		s.Reset()
		s.Close()
	}()

	s.SetDeadline(time.Unix(time.Now().Unix()+readTimeOut, 0))
	request := NewRequest(getBlockByHash, hash.Bytes())
	err = sendRequest(request, s)
	if err != nil {
		return nil, err
	}
	response, err := rm.ReadResponse(s)
	if response != nil && response.Code == Success {
		err := rlp.DecodeBytes(response.Body, &block)
		if err != nil {
			return nil, err
		}
	} else {
		return nil, fmt.Errorf("peer error: %v", err)
	}
	return block.TranslateToBlock(), nil
}

func (rm *RequestManager) GetHeaderByHeight(stream *p2p.StreamCreator, height uint64) (*types.Header, error) {
	var header *types.Header
	s, err := stream.NewStreamFunc(stream.PeerId)
	if err != nil {
		return nil, err
	}
	defer func() {
		s.Reset()
		s.Close()
	}()

	s.SetDeadline(time.Unix(time.Now().Unix()+readTimeOut, 0))
	sHeight := strconv.FormatUint(height, 10)
	request := NewRequest(getHeaderByHeight, []byte(sHeight))
	err = sendRequest(request, s)
	if err != nil {
		return nil, err
	}
	response, err := rm.ReadResponse(s)
	if response != nil && response.Code == Success {
		err := rlp.DecodeBytes(response.Body, &header)
		if err != nil {
			return nil, err
		}
	} else {
		return nil, fmt.Errorf("peer error: %v", err)
	}
	return header, nil
}

func (rm *RequestManager) SendBlock(stream *p2p.StreamCreator, block *types.Block) error {
	s, err := stream.NewStreamFunc(stream.PeerId)
	if err != nil {
		return err
	}
	defer func() {
		s.Reset()
		s.Close()
	}()

	s.SetDeadline(time.Unix(time.Now().Unix()+readTimeOut, 0))
	bytes, err := rlp.EncodeToBytes(block.TranslateToRlpBlock())
	if err != nil {
		return err
	}

	request := NewRequest(sendBlock, bytes)
	err = sendRequest(request, s)
	if err != nil {
		return ErrorPeerClose
	}

	response, err := rm.ReadResponse(s)
	if response != nil && response.Code != Success {
		return errors.New("send block failed")
	}
	return nil
}

func (rm *RequestManager) SendTransaction(stream *p2p.StreamCreator, tx types.ITransaction) error {
	s, err := stream.NewStreamFunc(stream.PeerId)
	if err != nil {
		return err
	}
	defer func() {
		s.Reset()
		s.Close()
	}()

	s.SetDeadline(time.Unix(time.Now().Unix()+readTimeOut, 0))
	bytes, err := rlp.EncodeToBytes(tx.TranslateToRlpTransaction())
	if err != nil {
		return err
	}
	request := NewRequest(sendTransaction, bytes)
	err = sendRequest(request, s)
	if err != nil {
		return err
	}
	response, err := rm.ReadResponse(s)
	if response != nil && response.Code != Success {
		return errors.New("send transaction failed")
	}
	return nil
}

func (rm *RequestManager) ValidationBlockHash(stream *p2p.StreamCreator, header *types.Header) (bool, error) {
	s, err := stream.NewStreamFunc(stream.PeerId)
	if err != nil {
		return false, err
	}
	defer func() {
		s.Reset()
		s.Close()
	}()

	s.SetDeadline(time.Unix(time.Now().Unix()+readTimeOut, 0))
	bytes, err := rlp.EncodeToBytes(header)
	if err != nil {
		return false, err
	}
	request := NewRequest(validationBlockHash, bytes)
	err = sendRequest(request, s)
	response, err := rm.ReadResponse(s)
	var rs bool
	if response != nil && response.Code == Success {
		err := rlp.DecodeBytes(response.Body, &rs)
		if err != nil {
			return false, err
		}
	} else {
		return false, fmt.Errorf("peer error: %v", err)
	}
	return rs, nil
}

func (rm *RequestManager) GetNodeInfo(stream *p2p.StreamCreator) (*types.NodeInfo, error) {
	s, err := stream.NewStreamFunc(stream.PeerId)
	if err != nil {
		return nil, err
	}
	defer func() {
		s.Reset()
		s.Close()
	}()

	s.SetDeadline(time.Unix(time.Now().Unix()+readTimeOut, 0))
	request := NewRequest(getNodeInfo, nil)
	err = sendRequest(request, s)
	response, err := rm.ReadResponse(s)
	var rs *types.NodeInfo
	if response != nil && response.Code == Success {
		err := rlp.DecodeBytes(response.Body, &rs)
		if err != nil {
			return nil, err
		}
	} else {
		return nil, fmt.Errorf("peer error: %v", err)
	}
	return rs, nil
}
