package types

import (
	"bytes"
	"github.com/uworldao/UWORLD/common/encode/rlp"
	hash2 "github.com/uworldao/UWORLD/common/hasharry"
	"github.com/uworldao/UWORLD/crypto/hash"
)

type Transactions []ITransaction

// Len returns the length of s
func (s Transactions) Len() int { return len(s) }

// Swap swaps the i'th and the j'th element in s
func (s Transactions) Swap(i, j int) { s[i], s[j] = s[j], s[i] }

// GetRlp implements Rlpable and returns the i'th element of s in rlp
func (s Transactions) GetRlp(i int) []byte {
	enc, _ := rlp.EncodeToBytes(s[i])
	return enc
}

func (s Transactions) Hash() hash2.Hash {
	var hashes [][]byte
	for _, tx := range s {
		hashes = append(hashes, tx.Hash().Bytes())
	}
	hashBytes := bytes.Join(hashes, []byte{})
	return hash.Hash(hashBytes)
}

func (s Transactions) SumFees() uint64 {
	var sum uint64
	for _, tx := range s {
		if tx.GetTxType() != ContractTransaction {
			sum += tx.GetFees()
		}
	}
	return sum
}

func (s Transactions) SumConsumption() uint64 {
	var sum uint64
	for _, tx := range s {
		if tx.GetTxType() == ContractTransaction {
			sum += tx.GetFees()
		}
	}
	return sum
}
