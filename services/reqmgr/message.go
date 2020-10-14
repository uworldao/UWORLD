package reqmgr

import (
	"encoding/json"
	"fmt"
	"github.com/libp2p/go-libp2p-core/network"
)

const (
	Success MessageCode = iota
	EncodeError
	DecodeError
	JsonError
	InternalError
)

const (
	getLastHeight       Method = "getLastHeight"
	getBlockByHeight    Method = "getBlockByHeight"
	getBlocksByHeight   Method = "getBlocksByHeight"
	getBlockByHash      Method = "getBlockByHash"
	getHeaderByHeight   Method = "getHeaderByHeight"
	getNodeInfo         Method = "getNodeInfo"
	sendBlock           Method = "sendBlock"
	sendTransaction     Method = "sendTransaction"
	validationBlockHash Method = "validationBlockHash"
)

const maxReadBytes = 1024 * 10
const MaxRequestBytes = maxReadBytes * 1000
const readTimeOut = 30

type MessageCode int

type Method string

type RawMessage []byte

// Peer communication request body
type Request struct {
	Method Method     `json:"method"`
	Body   RawMessage `json:"body"`
}

func NewRequest(method Method, body RawMessage) *Request {
	return &Request{Method: method, Body: body}
}

// Peer node communication response message
type Response struct {
	Code    MessageCode `json:"code"`
	Message string      `json:"message"`
	Body    RawMessage  `json:"body"`
}

func NewResponse(code MessageCode, message string, body []byte) *Response {
	return &Response{code, message, body}
}

// Read from request
func (rm *RequestManager) ReadRequest(stream network.Stream) (*Request, error) {
	reBytes, err := rm.readBytes(stream)
	request := &Request{}

	err = json.Unmarshal(reBytes, request)
	if err != nil {
		return nil, err
	}
	return request, nil
}

// Read from response
func (rm *RequestManager) ReadResponse(stream network.Stream) (*Response, error) {
	reBytes, err := rm.readBytes(stream)
	if err != nil {
		return nil, err
	}
	request := &Response{}

	err = json.Unmarshal(reBytes, request)
	if err != nil {
		return nil, err
	}
	return request, nil
}

// Read message bytes
func (rm *RequestManager) readBytes(stream network.Stream) ([]byte, error) {
	var reBytes []byte
	var err error
	var n int
	bytesLen := 0
	byteArr := rm.pool.Get().([]byte)
	for bytesLen < MaxRequestBytes {
		resetBytes(byteArr)
		n, err = stream.Read(byteArr)
		if err != nil {
			break
		}
		reBytes = append(reBytes, byteArr...)
		bytesLen += n
		if n < maxReadBytes {
			break
		}
	}
	rm.pool.Put(byteArr)
	if bytesLen > MaxRequestBytes {
		return nil, fmt.Errorf("request data must be less than %d", MaxRequestBytes)
	}
	return reBytes[0:bytesLen], err
}

func resetBytes(bytes []byte) {
	for i, _ := range bytes {
		bytes[i] = 0
	}
}
