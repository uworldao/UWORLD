package list

import (
	"container/heap"
	"github.com/uworldao/UWORLD/core"
	"github.com/uworldao/UWORLD/core/types"
)

type TxSortedMap struct {
	txs   map[string]types.ITransaction
	cache map[string]types.ITransaction
	index *txInfoList
}

func NewTxSortedMap() *TxSortedMap {
	return &TxSortedMap{
		txs:   make(map[string]types.ITransaction),
		cache: make(map[string]types.ITransaction),
		index: new(txInfoList),
	}
}

func (t *TxSortedMap) Put(tx types.ITransaction) {
	t.txs[tx.From().String()] = tx
	t.cache[tx.From().String()] = tx
	heap.Push(t.index, &txInfo{
		address: tx.From().String(),
		txHash:  tx.Hash().String(),
		fees:    tx.GetFees(),
		nonce:   tx.GetNonce(),
		time:    tx.GetTime(),
	})
}

func (t *TxSortedMap) GetAll() types.Transactions {
	var all types.Transactions
	for _, tx := range t.cache {
		all = append(all, tx)
	}
	return all
}

func (t *TxSortedMap) Gets(count int) types.Transactions {
	var txs types.Transactions
	rIndex := t.index.CopySelf()

	for rIndex.Len() > 0 && count > 0 {
		ti := heap.Pop(rIndex).(*txInfo)
		tx := t.txs[ti.address]
		txs = append(txs, tx)
		count--
	}
	return txs
}

func (t *TxSortedMap) GetByAddress(addr string) types.ITransaction {
	return t.txs[addr]
}

// If the transaction pool is full, delete the transaction with a small fee
func (t *TxSortedMap) PopMin(fees uint64) types.ITransaction {
	if t.Len() > 0 {
		if (*t.index)[0].fees <= fees {
			ti := heap.Remove(t.index, 0).(*txInfo)
			tx := t.txs[ti.address]
			delete(t.txs, ti.address)
			delete(t.cache, ti.address)
			return tx
		}
	}
	return nil
}

func (t *TxSortedMap) Len() int { return len(t.txs) }

func (t *TxSortedMap) IsExist(from string, txHash string) bool {
	tx, ok := t.cache[from]
	if ok {
		return tx.Hash().String() == txHash
	}
	return false
}

func (t *TxSortedMap) Remove(tx types.ITransaction) {
	for i, ti := range *(t.index) {
		if ti.txHash == tx.Hash().String() {
			heap.Remove(t.index, i)
			delete(t.txs, tx.From().String())
			delete(t.cache, tx.From().String())
			return
		}
	}
}

// Delete already packed transactions
func (t *TxSortedMap) RemoveExecuted(state core.IAccountState) {
	for _, tx := range t.cache {
		if err := state.VerifyState(tx); err != nil {
			t.Remove(tx)
		}
	}
}

// Delete expired transactions
func (t *TxSortedMap) RemoveExpiredTx(timeThreshold uint64) {
	for _, tx := range t.cache {
		if tx.GetTime() <= timeThreshold {
			t.Remove(tx)
		}
	}
}

type txInfoList []*txInfo

type txInfo struct {
	address string
	txHash  string
	fees    uint64
	nonce   uint64
	time    uint64
}

func (t txInfoList) Len() int           { return len(t) }
func (t txInfoList) Less(i, j int) bool { return t[i].fees > t[j].fees }
func (t txInfoList) Swap(i, j int)      { t[i], t[j] = t[j], t[i] }

func (t *txInfoList) Push(x interface{}) {
	*t = append(*t, x.(*txInfo))
}

func (t *txInfoList) Pop() interface{} {
	old := *t
	n := len(old)
	x := old[n-1]
	*t = old[0 : n-1]
	return x
}

func (t *txInfoList) CopySelf() *txInfoList {
	reReelList := new(txInfoList)
	for _, nonce := range *t {
		*reReelList = append(*reReelList, nonce)
	}
	return reReelList
}

func (t *txInfoList) FindIndex(addr string, nonce uint64) int {
	for index, ti := range *t {
		if ti.address == addr && ti.nonce <= nonce {
			return index
		}
	}
	return -1
}
