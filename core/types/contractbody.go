package types

import (
	"errors"
	"fmt"
	"github.com/uworldao/UWORLD/common/hasharry"
	"github.com/uworldao/UWORLD/param"
	"github.com/uworldao/UWORLD/ut"
)

const MaxCoinName = 32
const MaxCoinDescription = 1024

// The contract body includes the contract address,
// contract name, contract amount.
type ContractBody struct {
	Contract       hasharry.Address
	To             hasharry.Address
	Name           string
	Abbr           string
	Description    string
	Amount         uint64
	IncreaseSwitch bool
}

func (c *ContractBody) ToAddress() hasharry.Address {
	return c.To
}

func (c *ContractBody) GetAmount() uint64 {
	return c.Amount
}

func (c *ContractBody) GetContract() hasharry.Address {
	return c.Contract
}

func (c *ContractBody) GetName() string {
	return c.Name
}

func (c *ContractBody) GetAbbr() string {
	return c.Abbr
}

func (c *ContractBody) GetIncreaseSwitch() bool {
	return c.IncreaseSwitch
}

func (c *ContractBody) GetDescription() string {
	return c.Description
}

func (c *ContractBody) GetPeerId() []byte {
	return nil
}

func (c *ContractBody) VerifyBody(from hasharry.Address) error {
	if err := c.verifyAttribute(); err != nil {
		return err
	}
	if err := c.verifyContractAddress(from); err != nil {
		return err
	}
	if err := c.verifyContractTo(from); err != nil {
		return err
	}
	if err := c.verifyAmount(); err != nil {
		return err
	}
	if err := c.verifyIncreaseSwitch(); err != nil {
		return err
	}
	return nil
}

func (c *ContractBody) verifyAttribute() error {
	if len(c.Name) > MaxCoinName {
		return fmt.Errorf("the maximum length of coin name shall not exceed %d", MaxCoinName)
	}
	if len(c.Description) > MaxCoinDescription {
		return fmt.Errorf("the maximum length of coin description shall not exceed %d", MaxCoinDescription)
	}
	if err := ut.CheckAbbr(c.Abbr); err != nil {
		return err
	}
	return nil
}

func (c *ContractBody) verifyContractAddress(from hasharry.Address) error {
	if !ut.CheckContractAddress(param.Net, from.String(), c.Abbr, c.Contract.String()) {
		return errors.New("check contract address failed")
	}
	return nil
}

func (c *ContractBody) verifyContractTo(to hasharry.Address) error {
	if !ut.CheckUWDAddress(param.Net, to.String()) {
		return errors.New("check to address failed")
	}
	return nil
}

func (c *ContractBody) verifyAmount() error {
	if c.Amount > param.MaxContractCoin {
		return fmt.Errorf("the amount of money issued at one time shall not exceed %d", param.MaxContractCoin)
	}
	return nil
}

func (c *ContractBody) verifyIncreaseSwitch() error {
	if c.IncreaseSwitch {
		return fmt.Errorf("additional issuance is not allowed")
	}
	return nil
}
