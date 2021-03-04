package types

import (
	"github.com/uworldao/UWORLD/common/hasharry"
	"github.com/uworldao/UWORLD/param"
	"github.com/uworldao/UWORLD/ut"
)

// Voting transaction body
type VoteTransactionBody struct {
	To hasharry.Address
}

func (vt *VoteTransactionBody) ToAddress() hasharry.Address {
	return vt.To
}

func (vt *VoteTransactionBody) GetAmount() uint64 {
	return 0
}

func (vt *VoteTransactionBody) GetContract() hasharry.Address {
	return param.Token
}

func (vt *VoteTransactionBody) GetName() string {
	return ""
}

func (vt *VoteTransactionBody) GetAbbr() string {
	return ""
}

func (vt *VoteTransactionBody) GetIncreaseSwitch() bool {
	return false
}

func (vt *VoteTransactionBody) GetDescription() string {
	return ""
}

func (vt *VoteTransactionBody) GetPeerId() []byte {
	return nil
}

func (vt *VoteTransactionBody) VerifyBody(from hasharry.Address) error {
	if err := vt.verifyTo(); err != nil {
		return err
	}
	return nil
}

func (vt *VoteTransactionBody) verifyTo() error {
	if !ut.CheckUWDAddress(param.Net, vt.To.String()) {
		return ErrAddress
	}
	return nil
}
