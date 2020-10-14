package types

import (
	"github.com/jhdriver/UWORLD/common/hasharry"
	"github.com/jhdriver/UWORLD/param"
	"github.com/jhdriver/UWORLD/ut"
)

const maxContractLength = 50

// Ordinary transfer transaction body
type NormalTransactionBody struct {
	Contract hasharry.Address
	To       hasharry.Address
	Amount   uint64
}

func (nt *NormalTransactionBody) ToAddress() hasharry.Address {
	return nt.To
}

func (nt *NormalTransactionBody) GetAmount() uint64 {
	return nt.Amount
}

func (nt *NormalTransactionBody) GetContract() hasharry.Address {
	return nt.Contract
}

func (nt *NormalTransactionBody) GetName() string {
	return ""
}

func (nt *NormalTransactionBody) GetAbbr() string {
	return ""
}

func (nt *NormalTransactionBody) GetIncreaseSwitch() bool {
	return false
}

func (nt *NormalTransactionBody) GetDescription() string {
	return ""
}

func (nt *NormalTransactionBody) GetPeerId() []byte {
	return nil
}

func (nt *NormalTransactionBody) VerifyBody(from hasharry.Address) error {
	if err := nt.verifyContract(); err != nil {
		return err
	}
	if err := nt.verifyTo(); err != nil {
		return err
	}
	return nil
}

func (nt *NormalTransactionBody) verifyTo() error {
	if !ut.CheckUWDAddress(param.Net, nt.To.String()) {
		return ErrAddress
	}
	return nil
}

func (nt *NormalTransactionBody) verifyContract() error {
	if !ut.IsValidContractAddress(param.Net, nt.Contract.String()) {
		return ErrContractAddr
	}
	return nil
}
