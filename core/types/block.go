package types

import (
	hash2 "github.com/jhdriver/UWORLD/common/hasharry"
	"github.com/jhdriver/UWORLD/crypto/hash"
)

// Block structure
type Block struct {
	*Header
	*Body
}

func NewBlock(header *Header, body *Body) *Block {
	block := &Block{header, body}
	return block
}

func (b *Block) SetHash() {
	b.Hash = hash.Hash(b.Header.ToBytes())
}

func (b *Block) VerifyTxRoot() bool {
	return b.TxRoot.IsEqual(b.Body.Transactions.Hash())
}

func (b *Block) TranslateToRlpBlock() *RlpBlock {
	return &RlpBlock{
		Header:  b.Header,
		RlpBody: b.Body.TranslateToRlpBody(),
	}
}

func (b *Block) GetTxsLocations() map[hash2.Hash]*TxLocation {
	mapLocation := make(map[hash2.Hash]*TxLocation)
	for index, tx := range b.Transactions {
		mapLocation[tx.Hash()] = &TxLocation{
			TxRoot:  b.TxRoot,
			TxIndex: uint32(index),
			Height:  b.Header.Height,
		}
	}
	return mapLocation
}
