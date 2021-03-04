package types

import (
	"github.com/uworldao/UWORLD/common/hasharry"
	"github.com/uworldao/UWORLD/param"
)

const PeerIdLength = 53

type PeerId [PeerIdLength]byte

// Become a candidate trading subject and can participate
// in the next round of elections after success.
type LoginTransactionBody struct {
	PeerId
}

func (lit *LoginTransactionBody) GetPeerId() []byte {
	return lit.PeerId[:]
}

func (lit *LoginTransactionBody) GetAmount() uint64 {
	return 0
}

func (lit *LoginTransactionBody) GetContract() hasharry.Address {
	return param.Token
}

func (lit *LoginTransactionBody) GetName() string {
	return ""
}

func (lit *LoginTransactionBody) GetAbbr() string {
	return ""
}

func (lit *LoginTransactionBody) GetIncreaseSwitch() bool {
	return false
}

func (lit *LoginTransactionBody) GetDescription() string {
	return ""
}

func (lit *LoginTransactionBody) ToAddress() hasharry.Address {
	return hasharry.Address{}
}

func (lit *LoginTransactionBody) VerifyBody(from hasharry.Address) error {
	return nil
}
