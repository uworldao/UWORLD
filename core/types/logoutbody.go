package types

import (
	"github.com/uworldao/UWORLD/common/hasharry"
	"github.com/uworldao/UWORLD/param"
)

// Withdraw from the candidate trading body and no longer
// participate in elections after success.
type LogoutTransactionBody struct {
	Placeholder bool
}

func (lt *LogoutTransactionBody) ToAddress() hasharry.Address {
	return hasharry.Address{}
}

func (lt *LogoutTransactionBody) GetAmount() uint64 {
	return 0
}

func (lit *LogoutTransactionBody) GetContract() hasharry.Address {
	return param.Token
}

func (lit *LogoutTransactionBody) GetName() string {
	return ""
}

func (lit *LogoutTransactionBody) GetAbbr() string {
	return ""
}

func (lit *LogoutTransactionBody) GetDescription() string {
	return ""
}

func (lit *LogoutTransactionBody) GetIncreaseSwitch() bool {
	return false
}

func (lt *LogoutTransactionBody) GetPeerId() []byte {
	return nil
}

func (lt *LogoutTransactionBody) VerifyBody(from hasharry.Address) error {

	return nil
}
