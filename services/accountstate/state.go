package accountstate

import (
	"errors"
	"github.com/uworldao/UWORLD/common/hasharry"
	"github.com/uworldao/UWORLD/core/types"
	"github.com/uworldao/UWORLD/database/statedb"
	"github.com/uworldao/UWORLD/param"
	"sync"
	"time"
)

const accountSate = "account_state"

type AccountState struct {
	stateDb         IAccountStorage
	accountMutex    sync.RWMutex
	contractMutex   sync.RWMutex
	confirmedHeight uint64
}

func NewAccountState(dataDir string) (*AccountState, error) {
	storage := statedb.NewStateStorage(dataDir + "/" + accountSate)
	err := storage.Open()
	if err != nil {
		return nil, err
	}
	return &AccountState{
		stateDb: storage,
	}, nil
}

// Initialize account balance root hash
func (cs *AccountState) InitTrie(stateRoot hasharry.Hash) error {
	return cs.stateDb.InitTrie(stateRoot)
}

// Get account status, if the account status needs to be updated
// according to the effective block height, it will be updated,
// but not stored.
func (cs *AccountState) GetAccountState(stateKey hasharry.Address) types.IAccount {
	cs.accountMutex.RLock()
	account := cs.stateDb.GetAccountState(stateKey)
	cs.accountMutex.RUnlock()

	if account.IsNeedUpdate() {
		account = cs.updateAccountLocked(stateKey)
	}
	return account
}

func (cs *AccountState) GetAccountNonce(stateKey hasharry.Address) (uint64, error) {
	cs.accountMutex.RLock()
	defer cs.accountMutex.RUnlock()

	return cs.stateDb.GetAccountNonce(stateKey), nil
}

func (cs *AccountState) setAccountState(account types.IAccount) {
	cs.stateDb.SetAccountState(account)
}

// Update sender account status based on transaction information
func (cs *AccountState) UpdateFrom(tx types.ITransaction, blockHeight uint64) error {
	if tx.IsCoinBase() {
		return nil
	}

	cs.accountMutex.Lock()
	defer cs.accountMutex.Unlock()

	fromAccount := cs.stateDb.GetAccountState(tx.From())
	err := fromAccount.Update(cs.confirmedHeight)
	if err != nil {
		return err
	}

	err = fromAccount.FromChange(tx, blockHeight)
	if err != nil {
		return err
	}

	cs.setAccountState(fromAccount)
	return nil
}

// Update the receiver's account status based on transaction information
func (cs *AccountState) UpdateTo(tx types.ITransaction, blockHeight uint64) error {
	cs.accountMutex.Lock()
	defer cs.accountMutex.Unlock()

	var toAccount types.IAccount

	toAccount = cs.stateDb.GetAccountState(tx.GetTxBody().ToAddress())
	err := toAccount.Update(cs.confirmedHeight)
	if err != nil {
		return err
	}
	err = toAccount.ToChange(tx, blockHeight)
	if err != nil {
		return err
	}

	cs.setAccountState(toAccount)

	return nil
}

func (cs *AccountState) UpdateFees(fees, blockHeight uint64) error {
	cs.accountMutex.Lock()
	defer cs.accountMutex.Unlock()

	var account types.IAccount

	account = cs.stateDb.GetAccountState(param.FeeAddress)
	err := account.Update(cs.confirmedHeight)
	if err != nil {
		return err
	}
	account.FeesChange(fees, blockHeight)
	cs.setAccountState(account)
	return nil
}

func (cs *AccountState) UpdateConsumption(fees, blockHeight uint64) error {
	cs.accountMutex.Lock()
	defer cs.accountMutex.Unlock()

	var account types.IAccount

	account = cs.stateDb.GetAccountState(param.EaterAddress)
	err := account.Update(cs.confirmedHeight)
	if err != nil {
		return err
	}
	account.ConsumptionChange(fees, blockHeight)
	cs.setAccountState(account)
	return nil
}

// Update the locked balance of an account
func (cs *AccountState) updateAccountLocked(stateKey hasharry.Address) types.IAccount {
	account := cs.stateDb.GetAccountState(stateKey)
	account.Update(cs.confirmedHeight)
	return account
}

func (cs *AccountState) UpdateConfirmedHeight(height uint64) {
	cs.confirmedHeight = height
}

// Verify the status of the trading account
func (cs *AccountState) VerifyState(tx types.ITransaction) error {
	switch tx.GetTxType() {
	default:
		return cs.verifyTxState(tx)
	}
}

func (cs *AccountState) verifyTxState(tx types.ITransaction) error {
	if tx.GetTime() > uint64(time.Now().Unix()) {
		return errors.New("incorrect transaction time")
	}

	account := cs.GetAccountState(tx.From())
	return account.VerifyTxState(tx)
}

func (cs *AccountState) StateTrieCommit() (hasharry.Hash, error) {
	return cs.stateDb.Commit()
}

func (cs *AccountState) RootHash() hasharry.Hash {
	//cs.Print()
	return cs.stateDb.RootHash()
}

func (cs *AccountState) Print() {
	cs.stateDb.Print()
}

func (cs *AccountState) Close() error {
	return cs.stateDb.Close()
}
