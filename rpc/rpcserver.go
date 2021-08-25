package rpc

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/uworldao/UWORLD/common/hasharry"
	"github.com/uworldao/UWORLD/common/utils"
	"github.com/uworldao/UWORLD/config"
	"github.com/uworldao/UWORLD/consensus"
	"github.com/uworldao/UWORLD/core"
	coreTypes "github.com/uworldao/UWORLD/core/types"
	"github.com/uworldao/UWORLD/crypto/certgen"
	log "github.com/uworldao/UWORLD/log/log15"
	"github.com/uworldao/UWORLD/p2p"
	"github.com/uworldao/UWORLD/param"
	"github.com/uworldao/UWORLD/rpc/rpctypes"
	"github.com/uworldao/UWORLD/services/reqmgr"
	"github.com/uworldao/UWORLD/ut"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/reflection"
	"net"
	"os"
	"strconv"
)

type Server struct {
	config        *config.RpcConfig
	txPool        core.ITxPool
	accountState  core.IAccountState
	contractState core.IContractState
	consensus     consensus.IConsensus
	chain         core.IBlockChain
	grpcServer    *grpc.Server
	peerManager   p2p.IPeerManager
	peers         reqmgr.Peers
}

func NewServer(config *config.RpcConfig, txPool core.ITxPool, state core.IAccountState, contractState core.IContractState,
	consensus consensus.IConsensus, chain core.IBlockChain, peerManager p2p.IPeerManager, peers reqmgr.Peers) *Server {
	return &Server{config: config, txPool: txPool, accountState: state, contractState: contractState,
		consensus: consensus, chain: chain, peerManager: peerManager, peers: peers}
}

func (rs *Server) Start() error {
	lis, err := net.Listen("tcp", ":"+rs.config.RpcPort)
	if err != nil {
		return err
	}
	rs.grpcServer, err = rs.NewServer()
	if err != nil {
		return err
	}

	RegisterGreeterServer(rs.grpcServer, rs)
	reflection.Register(rs.grpcServer)
	go func() {
		if err := rs.grpcServer.Serve(lis); err != nil {
			log.Info("Rpc startup failed!", "err", err)
			os.Exit(1)
			return
		}

	}()
	if rs.config.RpcTLS {
		log.Info("Rpc startup successful", "port", rs.config.RpcPort, "pem", rs.config.RpcCert)
	} else {
		log.Info("Rpc startup successful", "port", rs.config.RpcPort)
	}
	return nil
}

func (rs *Server) Close() {
	rs.grpcServer.Stop()
	log.Info("GRPC server closed")
}

func (rs *Server) NewServer() (*grpc.Server, error) {
	var opts []grpc.ServerOption
	var interceptor grpc.UnaryServerInterceptor
	interceptor = rs.interceptor
	opts = append(opts, grpc.UnaryInterceptor(interceptor))

	// If tls is configured, generate tls certificate
	if rs.config.RpcTLS {
		if err := rs.generateCertFile(); err != nil {
			return nil, err
		}
		transportCredentials, err := credentials.NewServerTLSFromFile(rs.config.RpcCert, rs.config.RpcCertKey)
		if err != nil {
			return nil, err
		}
		opts = append(opts, grpc.Creds(transportCredentials))

	}

	// Set the maximum number of bytes received and sent
	opts = append(opts, grpc.MaxRecvMsgSize(reqmgr.MaxRequestBytes))
	opts = append(opts, grpc.MaxSendMsgSize(reqmgr.MaxRequestBytes))
	return grpc.NewServer(opts...), nil
}

func (rs *Server) SendTransaction(_ context.Context, req *Bytes) (*Response, error) {
	fmt.Println(string(req.String()))
	var rpcTx *coreTypes.RpcTransaction
	if err := json.Unmarshal(req.Bytes, &rpcTx); err != nil {
		return NewResponse(rpctypes.RpcErrParam, nil, err.Error()), nil
	}
	tx, err := coreTypes.TranslateRpcTxToTx(rpcTx)
	if err != nil {
		return NewResponse(rpctypes.RpcErrParam, nil, err.Error()), nil
	}
	if err := rs.txPool.Add(tx, false); err != nil {
		return NewResponse(rpctypes.RpcErrTxPool, nil, err.Error()), nil
	}
	return NewResponse(rpctypes.RpcSuccess, []byte(fmt.Sprintf("send transaction %s success", tx.Hash().String())), ""), nil
}

func (rs *Server) GetAccount(_ context.Context, req *Address) (*Response, error) {
	if !ut.CheckUWDAddress(param.Net, req.Address) {
		return NewResponse(rpctypes.RpcErrParam, nil, fmt.Sprintf("%s address check failed", req.Address)), nil
	}
	addr := hasharry.StringToAddress(req.Address)
	account := rs.accountState.GetAccountState(addr)
	rpcAccount := rpctypes.TranslateAccountToRpcAccount(account.(*coreTypes.Account))
	bytes, err := json.Marshal(rpcAccount)
	if err != nil {
		return NewResponse(rpctypes.RpcErrMarshal, nil, fmt.Sprintf("%s address not exsit", req.Address)), nil
	}
	return NewResponse(rpctypes.RpcSuccess, bytes, ""), nil
}

func (rs *Server) GetAccounts(_ context.Context, req *Null) (*Response, error) {
	accounts := rs.accountState.GetAccounts()
	var rpcAccounts []*rpctypes.Account
	for _, account := range accounts{
		rpcAccounts = append(rpcAccounts, rpctypes.TranslateAccountToRpcAccount(account.(*coreTypes.Account)))
	}
	bytes, _ := json.Marshal(rpcAccounts)
	return NewResponse(rpctypes.RpcSuccess, bytes, ""), nil
}

func (rs *Server) GetTransaction(ctx context.Context, req *Hash) (*Response, error) {
	hash, err := hasharry.StringToHash(req.Hash)
	if err != nil {
		return NewResponse(rpctypes.RpcErrParam, nil, "hash error"), nil
	}
	tx, err := rs.chain.GetTransaction(hash)
	if err != nil {
		return NewResponse(rpctypes.RpcErrBlockChain, nil, err.Error()), nil
	}
	confirmed := rs.chain.GetConfirmedHeight()
	index, err := rs.chain.GetTransactionIndex(hash)
	if err != nil {
		return NewResponse(rpctypes.RpcErrBlockChain, nil, fmt.Sprintf("%s is not exist", hash.String())), nil
	}
	height := index.GetHeight()
	rpcTx, _ := coreTypes.TranslateTxToRpcTx(tx.(*coreTypes.Transaction))
	rsMsg := &coreTypes.RpcTransactionConfirmed{
		TxHead:    rpcTx.TxHead,
		TxBody:    rpcTx.TxBody,
		Height:    height,
		Confirmed: confirmed >= height,
	}
	bytes, _ := json.Marshal(rsMsg)
	return NewResponse(rpctypes.RpcSuccess, bytes, ""), nil
}

func (rs *Server) GetBlockByHash(ctx context.Context, req *Hash) (*Response, error) {
	hash, err := hasharry.StringToHash(req.Hash)
	if err != nil {
		return NewResponse(rpctypes.RpcErrParam, nil, "hash error"), nil
	}
	block, err := rs.chain.GetBlockByHash(hash)
	if err != nil {
		return NewResponse(rpctypes.RpcErrBlockChain, nil, err.Error()), nil
	}
	rpcBlock, _ := coreTypes.TranslateBlockToRpcBlock(block, rs.chain.GetConfirmedHeight())
	bytes, err := json.Marshal(rpcBlock)
	if err != nil {
		return NewResponse(rpctypes.RpcErrMarshal, nil, err.Error()), nil
	}
	return NewResponse(rpctypes.RpcSuccess, bytes, ""), nil

}

func (rs *Server) GetBlockByHeight(_ context.Context, req *Height) (*Response, error) {
	block, err := rs.chain.GetBlockByHeight(req.Height)
	if err != nil {
		return NewResponse(rpctypes.RpcErrBlockChain, nil, err.Error()), nil
	}
	rpcBlock, _ := coreTypes.TranslateBlockToRpcBlock(block, rs.chain.GetConfirmedHeight())
	bytes, err := json.Marshal(rpcBlock)
	if err != nil {
		return NewResponse(rpctypes.RpcErrMarshal, nil, err.Error()), nil
	}
	return NewResponse(rpctypes.RpcSuccess, bytes, ""), nil
}

func (rs *Server) GetPoolTxs(context.Context, *Null) (*Response, error) {
	preparedTxs, futureTxs := rs.txPool.GetAll()
	txPoolTxs, _ := coreTypes.TranslateTxsToRpcTxPool(preparedTxs, futureTxs)
	bytes, err := json.Marshal(txPoolTxs)
	if err != nil {
		return NewResponse(rpctypes.RpcErrMarshal, nil, err.Error()), nil
	}
	return NewResponse(rpctypes.RpcSuccess, bytes, ""), nil
}

func (rs *Server) GetCandidates(context.Context, *Null) (*Response, error) {
	candidates := rs.consensus.GetCandidates(rs.chain)
	if candidates == nil || len(candidates) == 0 {
		return NewResponse(rpctypes.RpcErrDPos, nil, "no candidates"), nil
	}
	bytes, err := json.Marshal(coreTypes.TranslateCandidatesToRpcCandidates(candidates))
	if err != nil {
		return NewResponse(rpctypes.RpcErrMarshal, nil, err.Error()), nil
	}
	return NewResponse(rpctypes.RpcSuccess, bytes, ""), nil
}

func (rs *Server) GetLastHeight(context.Context, *Null) (*Response, error) {
	height := rs.chain.GetLastHeight()
	sHeight := strconv.FormatUint(height, 10)
	return NewResponse(rpctypes.RpcSuccess, []byte(sHeight), ""), nil
}

func (rs *Server) GetContract(ctx context.Context, req *Address) (*Response, error) {
	contract := rs.contractState.GetContract(req.Address)
	if contract == nil {
		return NewResponse(rpctypes.RpcErrContract, nil, fmt.Sprintf("contract address %s is not exist", req.Address)), nil
	}
	bytes, err := json.Marshal(coreTypes.TranslateContractToRpcContract(contract))
	if err != nil {
		return NewResponse(rpctypes.RpcErrMarshal, nil, err.Error()), nil
	}
	return NewResponse(rpctypes.RpcSuccess, bytes, ""), nil
}

func (rs *Server) GetConfirmedHeight(context.Context, *Null) (*Response, error) {
	height := rs.chain.GetConfirmedHeight()
	sHeight := strconv.FormatUint(height, 10)
	return NewResponse(rpctypes.RpcSuccess, []byte(sHeight), ""), nil
}

func (rs *Server) Peers(context.Context, *Null) (*Response, error) {
	peers := rs.peers.PeersInfo()
	peersJson, _ := json.Marshal(peers)
	return NewResponse(rpctypes.RpcSuccess, peersJson, ""), nil
}

func (rs *Server) NodeInfo(context.Context, *Null) (*Response, error) {
	node := rs.peers.NodeInfo()
	nodeJson, _ := json.Marshal(node)
	return NewResponse(rpctypes.RpcSuccess, nodeJson, ""), nil
}

func NewResponse(code int32, result []byte, err string) *Response {
	return &Response{Code: code, Result: result, Err: err}
}

// Authenticate rpc users
func (rs *Server) auth(ctx context.Context) error {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return errors.New("no token authentication information")
	}
	var (
		password string
	)

	if val, ok := md["password"]; ok {
		password = val[0]
	}

	if password != rs.config.RpcPass {
		return fmt.Errorf("the Token authentication information is invalid: password=%s", password)
	}
	return nil
}

func (rs *Server) interceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
	err = rs.auth(ctx)
	if err != nil {
		return
	}
	return handler(ctx, req)
}

func (rs *Server) generateCertFile() error {
	if rs.config.RpcCert == "" {
		rs.config.RpcCert = rs.config.DataDir + "/server.pem"
	}
	if rs.config.RpcCertKey == "" {
		rs.config.RpcCertKey = rs.config.DataDir + "/server.key"
	}
	if !utils.IsExist(rs.config.RpcCert) || !utils.IsExist(rs.config.RpcCertKey) {
		return certgen.GenCertPair(rs.config.RpcCert, rs.config.RpcCertKey)
	}
	return nil
}
