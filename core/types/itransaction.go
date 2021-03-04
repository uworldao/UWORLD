package types

import (
	"github.com/uworldao/UWORLD/common/hasharry"
	"github.com/uworldao/UWORLD/crypto/ecc/secp256k1"
)

type ITransaction interface {
	Size() uint64
	IsCoinBase() bool
	VerifyTx() error
	VerifyCoinBaseTx(height, sumFees uint64) error
	EncodeToBytes() ([]byte, error)
	SignTx(key *secp256k1.PrivateKey) error
	SetHash() error
	NonceKey() string
	TranslateToRlpTransaction() *RlpTransaction

	Hash() hasharry.Hash
	From() hasharry.Address
	GetFees() uint64
	GetNonce() uint64
	GetTime() uint64
	GetTxType() TransactionType
	GetSignScript() *SignScript
	GetTxHead() *TransactionHead
	GetTxBody() ITransactionBody
}

type ITransactionBody interface {
	GetAmount() uint64
	GetContract() hasharry.Address
	GetName() string
	GetAbbr() string
	GetDescription() string
	GetIncreaseSwitch() bool
	ToAddress() hasharry.Address
	GetPeerId() []byte
	VerifyBody(from hasharry.Address) error
}

type ITransactionIndex interface {
	GetHeight() uint64
}
