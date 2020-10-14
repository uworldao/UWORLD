package contractstate

import (
	"github.com/jhdriver/UWORLD/common/hasharry"
	"github.com/jhdriver/UWORLD/core/types"
	"github.com/jhdriver/UWORLD/database/contractdb"
	"sync"
)

const contractSate = "contract_state"

// Contract status, used to store all published contract information records
type ContractState struct {
	contractDb      IContractStorage
	contractMutex   sync.RWMutex
	confirmedHeight uint64
}

func NewContractState(dataDir string) (*ContractState, error) {
	storage := contractdb.NewContractStorage(dataDir + "/" + contractSate)
	err := storage.Open()
	if err != nil {
		return nil, err
	}
	return &ContractState{
		contractDb: storage,
	}, nil
}

// Initialize the contract state tree
func (cs *ContractState) InitTrie(contractRoot hasharry.Hash) error {
	return cs.contractDb.InitTrie(contractRoot)
}

func (cs *ContractState) RootHash() hasharry.Hash {
	return cs.contractDb.RootHash()
}

// Commit contract status changes
func (cs *ContractState) ContractTrieCommit() (hasharry.Hash, error) {
	return cs.contractDb.Commit()
}

func (c *ContractState) GetContract(contractAddr string) *types.Contract {
	c.contractMutex.RLock()
	defer c.contractMutex.RUnlock()

	contract := c.contractDb.GetContractState(contractAddr)
	return contract
}

func (c *ContractState) UpdateConfirmedHeight(height uint64) {
	c.confirmedHeight = height
}

// Verification contract
func (c *ContractState) VerifyState(tx types.ITransaction) error {
	c.contractMutex.RLock()
	defer c.contractMutex.RUnlock()

	if tx.GetTxType() != types.ContractTransaction {
		return nil
	}
	contractAddr := tx.GetTxBody().GetContract()
	contract := c.contractDb.GetContractState(contractAddr.String())
	if contract != nil {
		return contract.Verify(tx)
	}
	return nil
}

// Update contract status
func (c *ContractState) UpdateContract(tx types.ITransaction, blockHeight uint64) {
	c.contractMutex.Lock()
	defer c.contractMutex.Unlock()

	txBody := tx.GetTxBody()
	contractRecord := &types.ContractRecord{
		Height:   blockHeight,
		TxHash:   tx.Hash(),
		Time:     tx.GetTime(),
		Amount:   txBody.GetAmount(),
		Receiver: txBody.ToAddress().String(),
	}
	contractAddr := txBody.GetContract()
	contract := c.contractDb.GetContractState(contractAddr.String())
	if contract != nil {
		contract.AddContract(contractRecord)
	} else {
		contract = &types.Contract{
			Contract:       contractAddr.String(),
			CoinName:       txBody.GetName(),
			CoinAbbr:       txBody.GetAbbr(),
			Description:    txBody.GetDescription(),
			IncreaseSwitch: txBody.GetIncreaseSwitch(),
			Records: &types.RecordList{
				contractRecord,
			},
		}
	}
	c.contractDb.SetContractState(contract)
}

func (c *ContractState) Close() error {
	return c.contractDb.Close()
}
