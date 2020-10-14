package types

import (
	"bytes"
	"encoding/gob"
	"github.com/jhdriver/UWORLD/common/encode/rlp"
	hash2 "github.com/jhdriver/UWORLD/common/hasharry"
	"github.com/jhdriver/UWORLD/crypto/hash"
)

const BlockVersion = 1

// Block header structure, used to verify the block
type Header struct {
	// Block header version
	Version uint32
	// Block hash
	Hash hash2.Hash
	// The previous block hash
	ParentHash hash2.Hash
	// Hash of all transactions
	TxRoot hash2.Hash
	// Account status tree root hash
	StateRoot hash2.Hash
	// Contract status tree root hash
	ContractRoot hash2.Hash
	// consensus status tree root hash
	ConsensusRoot hash2.Hash
	// Block height
	Height uint64
	// Blocks of time
	Time uint64
	// The block belongs to the election cycle
	Term uint64
	// Block signature
	SignScript *SignScript
	// Block generator
	Signer hash2.Address
}

func (h *Header) ToBytes() []byte {
	bytes, _ := rlp.EncodeToBytes(h)
	return bytes
}

func (h *Header) SetHash() {
	h.Hash = hash.Hash(h.ToBytes())
}

func (h *Header) Serialize() ([]byte, error) {
	var buff bytes.Buffer
	encode := gob.NewEncoder(&buff)
	err := encode.Encode(h)
	if err != nil {
		return nil, err
	}
	return buff.Bytes(), nil
}

func (h *Header) DeSerialize(data []byte) (*Header, error) {
	decode := gob.NewDecoder(bytes.NewReader(data))
	err := decode.Decode(&h)
	if err != nil {
		return nil, err
	}
	return h, nil
}

func (h *Header) HashString() string {
	return h.Hash.String()
}

func (h *Header) ParentHashString() string {
	return h.ParentHash.String()
}

func (h *Header) TxRootString() string {
	return h.TxRoot.String()
}
