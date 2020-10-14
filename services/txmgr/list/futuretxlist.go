package list

import (
	"fmt"
	"github.com/jhdriver/UWORLD/core/types"
)

type FutureTxList struct {
	Txs        map[string]types.ITransaction
	nonceKeMap map[string]string
}

func NewFutureTxList() *FutureTxList {
	return &FutureTxList{
		Txs:        make(map[string]types.ITransaction),
		nonceKeMap: make(map[string]string),
	}
}

func (f *FutureTxList) Put(tx types.ITransaction) error {
	if f.IsExist(tx.Hash().String()) {
		return fmt.Errorf("transation hash %s exsit", tx.Hash())
	}
	if oldTxHash := f.GetNonceKeyHash(tx.NonceKey()); oldTxHash != "" {
		oldTx := f.Txs[oldTxHash]
		if oldTx.GetFees() > tx.GetFees() {
			return fmt.Errorf("transation nonce %d exist, the fees must biger than before %d", tx.GetNonce, oldTx.GetFees())
		}
		f.Remove(oldTx)
	}
	f.Txs[tx.Hash().String()] = tx
	f.nonceKeMap[tx.NonceKey()] = tx.Hash().String()
	return nil
}

func (f *FutureTxList) Remove(tx types.ITransaction) {
	delete(f.Txs, tx.Hash().String())
	delete(f.nonceKeMap, tx.NonceKey())
}

func (f *FutureTxList) IsExist(txHash string) bool {
	_, ok := f.Txs[txHash]
	return ok
}

func (f *FutureTxList) GetNonceKeyHash(nonceKey string) string {
	return f.nonceKeMap[nonceKey]
}

func (f *FutureTxList) Len() int {
	return len(f.Txs)
}

func (f *FutureTxList) GetAll() types.Transactions {
	var all types.Transactions
	for _, tx := range f.Txs {
		all = append(all, tx)
	}
	return all
}
