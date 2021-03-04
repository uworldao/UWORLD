package types

import (
	"errors"
	"fmt"
	"github.com/uworldao/UWORLD/common/hasharry"
	"github.com/uworldao/UWORLD/param"
)

// Contract structure, issuing a contract with the same
// name is equivalent to reissuing the pass
type Contract struct {
	Contract       string
	CoinName       string
	CoinAbbr       string
	Description    string
	IncreaseSwitch bool
	Records        *RecordList
}

func NewContract() *Contract {
	return &Contract{Records: &RecordList{}}
}

func (c *Contract) IsExist(txHash hasharry.Hash) bool {
	for _, r := range *c.Records {
		if txHash.IsEqual(r.TxHash) {
			return true
		}
	}
	return false
}

func (c *Contract) Verify(tx ITransaction) error {
	txBody := tx.GetTxBody()
	if c.Contract != "" && !c.IncreaseSwitch {
		return errors.New("this contract does not support additional issuance")
	}
	if c.CoinName != txBody.GetName() {
		return errors.New("coin name is not consistent")
	}
	if c.CoinAbbr != txBody.GetAbbr() {
		return errors.New("coin abbr is not consistent")
	}
	if c.Contract != txBody.GetContract().String() {
		return errors.New("contract address is not consistent")
	}
	if c.IsExist(tx.Hash()) {
		return errors.New("duplicate transaction hash")
	}
	if txBody.GetAmount() > param.MaxContractCoin {
		return fmt.Errorf("the amount of money issued at one time shall not exceed %d", param.MaxContractCoin)
	}
	if c.amount()+txBody.GetAmount() > param.MaxAllContractCoin {
		return fmt.Errorf("the total amount of money issued shall not exceed %d", param.MaxAllContractCoin)
	}
	return nil
}

func (c *Contract) AddContract(record *ContractRecord) {
	c.Records.Set(record)
}

func (c *Contract) FallBack(height uint64) error {
	for _, record := range *c.Records {
		if record.Height > height {
			c.Records.Remove(height)
		}
	}
	return nil
}

func (c *Contract) amount() uint64 {
	var sum uint64
	for _, record := range *c.Records {
		sum += record.Amount
	}
	return sum
}

type ContractRecord struct {
	Height   uint64
	TxHash   hasharry.Hash
	Time     uint64
	Amount   uint64
	Receiver string
}

type RecordList []*ContractRecord

func (r *RecordList) Get(height uint64) (*ContractRecord, bool) {
	for _, record := range *r {
		if height == record.Height {
			return record, true
		}
	}
	return &ContractRecord{}, false
}

func (r *RecordList) Set(newRecord *ContractRecord) {
	for i, record := range *r {
		if newRecord.Height == record.Height {
			(*r)[i] = newRecord
			return
		}
	}
	*r = append(*r, newRecord)
}

func (r *RecordList) Remove(height uint64) {
	for i, record := range *r {
		if record.Height == height {
			(*r) = append((*r)[0:i], (*r)[i+1:]...)
			return
		}
	}
}

func (r *RecordList) Len() int {
	return len(*r)
}
