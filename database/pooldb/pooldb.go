package pooldb

import (
	"github.com/uworldao/UWORLD/common/encode/rlp"
	"github.com/uworldao/UWORLD/core/types"
	"github.com/uworldao/UWORLD/database/leveldb"
	"github.com/uworldao/UWORLD/services/txmgr/list"
)

const (
	futureTxs  = "futureTxs"
	prepareTxs = "prepareTxs"
)

type TxPoolStorage struct {
	db *leveldb.Base
}

func NewTxPoolStorage(path string) *TxPoolStorage {
	return &TxPoolStorage{&leveldb.Base{path, nil}}
}

func (t *TxPoolStorage) Open() error {
	return t.db.Open()
}

func (t *TxPoolStorage) LoadFutureTxs() *list.FutureTxList {
	future := list.NewFutureTxList()
	rs := t.db.Foreach(futureTxs)
	for _, value := range rs {
		var rlpTx *types.RlpTransaction
		rlp.DecodeBytes(value, &rlpTx)
		future.Put(rlpTx.TranslateToTransaction())
	}
	return future
}

func (t *TxPoolStorage) LoadPreparesTxs() *list.TxSortedMap {
	prepare := list.NewTxSortedMap()
	rs := t.db.Foreach(prepareTxs)
	for _, value := range rs {
		var rlpTx *types.RlpTransaction
		rlp.DecodeBytes(value, &rlpTx)
		prepare.Put(rlpTx.TranslateToTransaction())
	}
	return prepare
}

func (t *TxPoolStorage) SaveFutureTxs(future *list.FutureTxList) {
	t.db.ClearBucket(futureTxs)
	for _, tx := range future.Txs {
		rlpTx := tx.TranslateToRlpTransaction()
		bytes, _ := rlp.EncodeToBytes(rlpTx)
		key := leveldb.GetKey(futureTxs, tx.Hash().Bytes())
		t.db.UpdateValue(key, bytes)
	}
}

func (t *TxPoolStorage) SavePreparesTxs(prepare *list.TxSortedMap) {
	t.db.ClearBucket(prepareTxs)
	for _, tx := range prepare.GetAll() {
		rlpTx := tx.TranslateToRlpTransaction()
		bytes, _ := rlp.EncodeToBytes(rlpTx)
		key := leveldb.GetKey(futureTxs, tx.Hash().Bytes())
		t.db.UpdateValue(key, bytes)
	}
}

func (t *TxPoolStorage) Close() error {
	return t.db.Db.Close()
}
