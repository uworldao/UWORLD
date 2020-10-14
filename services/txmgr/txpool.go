package txmgr

import (
	"errors"
	"github.com/jhdriver/UWORLD/config"
	"github.com/jhdriver/UWORLD/consensus"
	"github.com/jhdriver/UWORLD/core"
	"github.com/jhdriver/UWORLD/core/types"
	"github.com/jhdriver/UWORLD/database/pooldb"
	log "github.com/jhdriver/UWORLD/log/log15"
	"github.com/jhdriver/UWORLD/p2p"
	"github.com/jhdriver/UWORLD/services/blkmgr"
	"github.com/jhdriver/UWORLD/services/txmgr/list"
	"github.com/libp2p/go-libp2p-core/peer"
	"time"
)

// Clear the expired transaction interval
const monitorTxInterval = 20

const txChanLength = 50

// Maximum number of transactions in the transaction pool
const maxPoolTx = 5000

const txPoolStorage = "txpool"

// Manage transactions not packaged into blocks
type TxPool struct {
	accountState  core.IAccountState
	contractState core.IContractState
	consensus     consensus.IConsensus
	txs           *list.TxList
	peerManager   p2p.IPeerManager
	network       blkmgr.Network
	newStream     blkmgr.ICreateStream
	txChan        chan types.ITransaction
	recTx         chan types.ITransaction
	removeTxsCh   chan types.Transactions
	stateUpdateCh chan struct{}
	stop          chan bool
}

func NewTxPool(config *config.Config, accountState core.IAccountState, contractState core.IContractState, consensus consensus.IConsensus, peerManager p2p.IPeerManager, network blkmgr.Network,
	recTx chan types.ITransaction, stateUpdateCh chan struct{}, removeTxsCh chan types.Transactions,
	newStream blkmgr.ICreateStream) *TxPool {

	return &TxPool{
		accountState:  accountState,
		contractState: contractState,
		consensus:     consensus,
		txs:           list.NewTxList(accountState, pooldb.NewTxPoolStorage(config.DataDir+"/"+txPoolStorage)),
		peerManager:   peerManager,
		network:       network,
		recTx:         recTx,
		removeTxsCh:   removeTxsCh,
		stateUpdateCh: stateUpdateCh,
		newStream:     newStream,
		txChan:        make(chan types.ITransaction, txChanLength),
		stop:          make(chan bool, 1),
	}
}

// Start transaction pool
func (tp *TxPool) Start() error {
	if err := tp.txs.Load(); err != nil {
		return err
	}

	go tp.monitorTxTime()
	go tp.dealTx()

	log.Info("Transaction pool startup successful")
	return nil
}

func (tp *TxPool) Stop() error {
	tp.stop <- true
	log.Info("Stop transaction pool")
	return tp.txs.Close()
}

// Monitor transaction time
func (tp *TxPool) monitorTxTime() {
	t := time.NewTicker(time.Second * monitorTxInterval)
	defer t.Stop()

	for range t.C {
		tp.clearExpiredTx()
	}
}

func (tp *TxPool) dealTx() {
	for {
		select {
		case _ = <-tp.stop:
			return
		case tx := <-tp.txChan:
			go tp.broadcastTx(tx)
		case tx := <-tp.recTx:
			go tp.Add(tx, true)
		case txs := <-tp.removeTxsCh:
			go tp.Remove(txs)
		case _ = <-tp.stateUpdateCh:
			go tp.txs.UpdateTxsList()
		}
	}
}

// Broadcast transaction
func (tp *TxPool) broadcastTx(tx types.ITransaction) {
	peers := tp.peerManager.Peers()
	for id, _ := range peers {
		if id != tp.peerManager.LocalPeerInfo().AddrInfo.ID.String() {
			peerId := new(peer.ID)
			if err := peerId.UnmarshalText([]byte(id)); err == nil {
				streamCreator := p2p.StreamCreator{PeerId: *peerId, NewStreamFunc: tp.newStream.CreateStream}
				go tp.network.SendTransaction(&streamCreator, tx)
			}
		}
	}
}

func (tp *TxPool) Add(tx types.ITransaction, isPeer bool) error {
	return tp.AddTransaction(tx, isPeer)
}

// Verify adding transactions to the transaction pool
func (tp *TxPool) AddTransaction(tx types.ITransaction, isPeer bool) error {
	if tp.IsExist(tx) {
		return errors.New("the transaction already exists")
	}

	if err := tp.verifyTx(tx); err != nil {
		return err
	}

	if tp.txs.Len() >= maxPoolTx {
		tp.txs.RemoveMinFeeTx(tx)
	}

	if err := tp.txs.Put(tx); err != nil {
		return err
	}
	log.Info("TxPool receive transaction", "hash", tx.Hash())
	if !isPeer {
		tp.txChan <- tx
	}
	return nil
}

// Get transactions from the transaction pool
func (tp *TxPool) Gets(count int) types.Transactions {
	txs := tp.txs.Gets(count)
	failed := types.Transactions{}
	for i, tx := range txs {
		if err := tp.verifyTx(tx); err != nil {
			failed = append(failed, tx)
			txs = append(txs[0:i], txs[i+1:txs.Len()]...)
		}
	}
	tp.Remove(failed)
	return txs
}

// Get all transactions in the trading pool
func (tp *TxPool) GetAll() (types.Transactions, types.Transactions) {
	prepareTxs, futureTxs := tp.txs.GetAll()
	return prepareTxs, futureTxs
}

func (tp *TxPool) Get() types.ITransaction {
	panic("implement me")
}

// Delete transaction
func (tp *TxPool) Remove(txs types.Transactions) {
	for _, tx := range txs {
		switch tx.GetTxType {
		default:
			tp.txs.Remove(tx)
		}
	}

}

func (tp *TxPool) IsExist(tx types.ITransaction) bool {
	return tp.txs.IsExist(tx.From().String(), tx.Hash().String())
}

// Verify the transaction is legal
func (tp *TxPool) verifyTx(tx types.ITransaction) error {
	if err := tx.VerifyTx(); err != nil {
		return err
	}

	if err := tp.consensus.VerifyTx(tx); err != nil {
		return err
	}

	if err := tp.accountState.VerifyState(tx); err != nil {
		return err
	}

	if err := tp.contractState.VerifyState(tx); err != nil {
		return err
	}

	return nil
}

func (tp *TxPool) clearExpiredTx() {
	timeThreshold := time.Now().Unix() - list.TxLifeTime
	tp.txs.RemoveExpiredTx(uint64(timeThreshold))
}
