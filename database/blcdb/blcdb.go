package blcdb

import (
	"fmt"
	"github.com/jhdriver/UWORLD/common/encode/rlp"
	"github.com/jhdriver/UWORLD/common/hasharry"
	"github.com/jhdriver/UWORLD/core/types"
	"github.com/jhdriver/UWORLD/database/leveldb"
	"strconv"
)

const (
	lastHeight        = "lastHeight"
	headerBucket      = "headerBucket"
	transactionBucket = "transactionBucket"
	heightHash        = "heightHash"
	locationBucket    = "locationBucket"
	stateRoot         = "stateRoot"
	contractRoot      = "contractRoot"
	consensusRoot     = "consensusRoot"
	historyConfirmed  = "historyConfirmed"
	termLastHash      = "termLastHash"
)

type BlockChainStorage struct {
	db *leveldb.Base
}

func (b *BlockChainStorage) Open() error {
	err := b.db.Open()
	if err != nil {
		return err
	}
	return b.initBucket()
}

func NewBlockChainStorage(path string) *BlockChainStorage {
	return &BlockChainStorage{&leveldb.Base{path, nil}}
}

func (b *BlockChainStorage) Close() error {
	return b.db.Close()
}

func (b *BlockChainStorage) initBucket() error {

	return nil
}

func (b *BlockChainStorage) GetHeaderByHeight(height uint64) (*types.Header, error) {
	hash, err := b.GetHashByHeight(height)
	if err != nil {
		return nil, err
	}
	return b.GetHeader(hash)
}

func (b *BlockChainStorage) GetHeader(hash hasharry.Hash) (*types.Header, error) {
	key := leveldb.GetKey(headerBucket, hash.Bytes())
	bytes, err := b.db.GetValue(key)
	if err != nil {
		return nil, err
	}
	header := new(types.Header)
	err = rlp.DecodeBytes(bytes, header)
	return header, err
}

func (b *BlockChainStorage) GetTxLocation(hash hasharry.Hash) (*types.TxLocation, error) {
	var txLoc *types.TxLocation
	key := leveldb.GetKey(locationBucket, hash.Bytes())
	bytes, err := b.db.GetValue(key)
	if err != nil {
		return nil, err
	}
	if bytes == nil || len(bytes) == 0 {
		return nil, fmt.Errorf("transaction %s is not exist", hash.String())
	}
	err = rlp.DecodeBytes(bytes, &txLoc)
	return txLoc, err
}

func (b *BlockChainStorage) GetTransactions(txRoot hasharry.Hash) ([]*types.RlpTransaction, error) {
	var txs []*types.RlpTransaction
	key := leveldb.GetKey(transactionBucket, txRoot.Bytes())
	bytes, err := b.db.GetValue(key)
	if err != nil {
		return nil, err
	}
	err = rlp.DecodeBytes(bytes, &txs)
	return txs, err
}

func (b *BlockChainStorage) GetTransaction(hash hasharry.Hash) (*types.RlpTransaction, error) {
	txLoc, err := b.GetTxLocation(hash)
	if err != nil {
		return nil, err
	}
	txs, err := b.GetTransactions(txLoc.TxRoot)
	if err != nil {
		return nil, err
	}
	return txs[txLoc.TxIndex], nil
}

func (b *BlockChainStorage) GetLastHeight() (uint64, error) {
	bytes, err := b.db.GetValue([]byte(lastHeight))
	if err != nil {
		return 0, err
	}
	return strconv.ParseUint(string(bytes), 10, 64)
}

func (b *BlockChainStorage) GetHashByHeight(height uint64) (hasharry.Hash, error) {
	bytes := leveldb.GetKey(heightHash, []byte(strconv.FormatUint(height, 10)))
	hash, err := b.db.GetValue(bytes)
	if err != nil {
		return hasharry.Hash{}, nil
	}
	return hasharry.BytesToHash(hash), nil
}

func (b *BlockChainStorage) GetHistoryConfirmedHeight(height uint64) (uint64, error) {
	bytes := leveldb.GetKey(historyConfirmed, []byte(strconv.FormatUint(height, 10)))
	heightBytes, err := b.db.GetValue(bytes)
	if err != nil {
		return 0, err
	}
	return strconv.ParseUint(string(heightBytes), 10, 64)
}

func (b *BlockChainStorage) GetStateRoot() (hasharry.Hash, error) {
	rootBytes, err := b.db.GetValue([]byte(stateRoot))
	if err != nil {
		return hasharry.Hash{}, err
	}
	return hasharry.BytesToHash(rootBytes), nil
}

func (b *BlockChainStorage) GetContractRoot() (hasharry.Hash, error) {
	rootBytes, err := b.db.GetValue([]byte(contractRoot))
	if err != nil {
		return hasharry.Hash{}, err
	}
	return hasharry.BytesToHash(rootBytes), nil
}

func (b *BlockChainStorage) GetConsensusRoot() (hasharry.Hash, error) {
	rootBytes, err := b.db.GetValue([]byte(consensusRoot))
	if err != nil {
		return hasharry.Hash{}, err
	}
	return hasharry.BytesToHash(rootBytes), nil
}

func (b *BlockChainStorage) GetTermLastHash(term uint64) (hasharry.Hash, error) {
	bytes := []byte(strconv.FormatUint(term, 10))
	key := leveldb.GetKey(termLastHash, bytes)
	bytes, err := b.db.GetValue(key)
	if err != nil {
		return hasharry.Hash{}, err
	}
	return hasharry.BytesToHash(bytes), nil
}

func (b *BlockChainStorage) UpdateLastHeight(height uint64) {
	bytes := []byte(strconv.FormatUint(height, 10))
	b.db.UpdateValue([]byte(lastHeight), bytes)
}

func (b *BlockChainStorage) UpdateHeader(header *types.Header) {
	bytes, _ := rlp.EncodeToBytes(header)
	key := leveldb.GetKey(headerBucket, header.Hash.Bytes())
	b.db.UpdateValue(key, bytes)
	b.UpdateHeightHash(header.Height, header.Hash)
}

func (b *BlockChainStorage) UpdateTransactions(txRoot hasharry.Hash, iTxs []*types.RlpTransaction) {
	bytes, _ := rlp.EncodeToBytes(iTxs)
	key := leveldb.GetKey(transactionBucket, txRoot.Bytes())
	b.db.UpdateValue(key, bytes)
}

func (b *BlockChainStorage) UpdateTxLocation(txLocs map[hasharry.Hash]*types.TxLocation) {
	for hash, loc := range txLocs {
		locBytes, _ := rlp.EncodeToBytes(loc)
		key := leveldb.GetKey(locationBucket, hash.Bytes())
		b.db.UpdateValue(key, locBytes)
	}
}

func (b *BlockChainStorage) UpdateHeightHash(height uint64, hash hasharry.Hash) {
	bytes := []byte(strconv.FormatUint(height, 10))
	key := leveldb.GetKey(heightHash, bytes)
	b.db.UpdateValue(key, hash.Bytes())
}

func (b *BlockChainStorage) UpdateStateRoot(hash hasharry.Hash) {
	b.db.UpdateValue([]byte(stateRoot), hash.Bytes())
}

func (b *BlockChainStorage) UpdateContractRoot(hash hasharry.Hash) {
	b.db.UpdateValue([]byte(contractRoot), hash.Bytes())
}

func (b *BlockChainStorage) UpdateConsensusRoot(hash hasharry.Hash) {
	b.db.UpdateValue([]byte(consensusRoot), hash.Bytes())
}

func (b *BlockChainStorage) UpdateHistoryConfirmedHeight(height uint64, confirmedHeight uint64) {
	heightBytes := []byte(strconv.FormatUint(height, 10))
	confirmedBytes := []byte(strconv.FormatUint(confirmedHeight, 10))
	key := leveldb.GetKey(historyConfirmed, heightBytes)
	b.db.UpdateValue(key, confirmedBytes)
}

func (b *BlockChainStorage) UpdateTermLastHash(term uint64, hash hasharry.Hash) {
	bytes := []byte(strconv.FormatUint(term, 10))
	key := leveldb.GetKey(termLastHash, bytes)
	b.db.UpdateValue(key, hash.Bytes())
}
