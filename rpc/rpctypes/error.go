package rpctypes

const (
	RpcSuccess int32 = iota
	RpcErrTxPool
	RpcErrMarshal
	RpcErrBlockChain
	RpcErrDPos
	RpcErrParam
	RpcErrContract
)
