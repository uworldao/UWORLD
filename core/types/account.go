package types

import (
	"errors"
	"fmt"
	"github.com/uworldao/UWORLD/common/hasharry"
	"github.com/uworldao/UWORLD/crypto/hash"
	"github.com/uworldao/UWORLD/param"
)

// Account information, not sending a transaction or sending an
// action will increase the nonce value by 1
type Account struct {
	Address         hasharry.Address
	Nonce           uint64
	Time            uint64
	ConfirmedHeight uint64
	ConfirmedNonce  uint64
	ConfirmedTime   uint64
	Coins           *Coins
	JournalIn       *journalIn
	JournalOut      *journalOut
}

// Calculate user status key
/*func AccountStateKey(address []byte, others ...[]byte) hasharry.Hash {
	for _, other := range others {
		address = append(address, other...)
	}
	return hash.Hash(address)
}*/

func AccountStateKeyString(address []byte, others ...[]byte) string {
	for _, other := range others {
		address = append(address, other...)
	}
	hash := hash.Hash(address)
	return hash.String()
}

func NewAccount() *Account {
	coin := &CoinAccount{
		Contract:  param.Token.String(),
		Balance:   0,
		LockedIn:  0,
		LockedOut: 0,
	}
	return &Account{
		Address:         hasharry.Address{},
		Nonce:           0,
		Time:            0,
		ConfirmedHeight: 0,
		ConfirmedNonce:  0,
		ConfirmedTime:   0,
		Coins:           &Coins{coin},
		JournalIn:       newJournalIn(),
		JournalOut:      newJournalOut(),
	}
}

func (a *Account) FallBack(height uint64) error {
	if height < a.ConfirmedHeight {
		return errors.New("too small fall back height")
	}
	if height > a.ConfirmedHeight {
		if err := a.Update(height); err != nil {
			return err
		}
	}
	a.Nonce = a.ConfirmedNonce
	a.Time = a.ConfirmedTime
	amounts := a.JournalIn.Amount()
	for contract, amount := range amounts {
		coinAccount, ok := a.Coins.Get(contract)
		if !ok {
			return errors.New("wrong journal")
		}
		coinAccount.Balance += amount
	}

	for i, coinAccount := range *a.Coins {
		coinAccount.LockedOut = 0
		coinAccount.LockedIn = 0
		(*a.Coins)[i] = coinAccount
	}

	a.JournalIn = newJournalIn()
	a.JournalOut = newJournalOut()
	return nil
}

// Calculate the available balance of the current account based on the current effective block height
func (a *Account) Update(confirmedHeight uint64) error {
	confirmedNonce := a.ConfirmedNonce
	confirmedTime := a.ConfirmedTime
	// Update through the account transfer log information
	for _, in := range a.JournalIn.GetJournalIns(confirmedHeight) {
		coinAccount, ok := a.Coins.Get(in.Contract)
		if !ok {
			return errors.New("wrong journal")
		}
		if coinAccount.LockedIn >= in.Amount {
			coinAccount.LockedIn -= in.Amount
			a.Coins.Set(coinAccount)

			tokenAccount, ok := a.Coins.Get(param.Token.String())
			if !ok {
				return errors.New("wrong journal")
			}
			if tokenAccount.LockedIn >= in.Fees {
				tokenAccount.LockedIn -= in.Fees
				a.Coins.Set(tokenAccount)
			} else {
				return errors.New("locked in amount not enough when update account journal")
			}
			a.JournalIn.Remove(in.Height)

		} else {
			return errors.New("locked in amount not enough when update account journal")
		}

		if in.Nonce > confirmedNonce {
			confirmedNonce = in.Nonce
		}
		if in.Time > confirmedTime {
			confirmedTime = in.Time
		}
	}

	// Update through account transfer log information
	for _, out := range a.JournalOut.GetJournalOuts(confirmedHeight) {
		coinAccount, ok := a.Coins.Get(out.Contract)
		if !ok {
			coinAccount = &CoinAccount{
				Contract:  out.Contract,
				Balance:   0,
				LockedIn:  0,
				LockedOut: 0,
			}
		}
		if coinAccount.LockedOut >= out.Amount {
			coinAccount.Balance += out.Amount
			coinAccount.LockedOut -= out.Amount
			a.Coins.Set(coinAccount)
			a.JournalOut.Remove(out.Height, out.Contract)
		} else {
			return errors.New("locked out amount not enough when update account Journal")
		}
	}
	a.ConfirmedHeight = confirmedHeight
	a.ConfirmedNonce = confirmedNonce
	a.ConfirmedTime = confirmedTime
	return nil
}

func (a *Account) StateKey() hasharry.Address {
	return a.Address
}

func (a *Account) IsExist() bool {
	return !hasharry.EmptyAddress(a.Address)
}

// Determine whether the account needs to be updated. If both
// the transfer-out and transfer-in are 0, no update is required.
func (a *Account) IsNeedUpdate() bool {
	for _, coinContract := range *a.Coins {
		if coinContract.LockedIn != 0 || coinContract.LockedOut != 0 {
			return true
		}
	}
	return false
}

// Change the account status of the party that transferred the transaction
func (a *Account) FromChange(tx ITransaction, blockHeight uint64) error {
	if tx.GetTxType() == ContractTransaction {
		return a.fromContractChange(tx, blockHeight)
	}
	if a.Nonce+1 != tx.GetNonce() {
		return ErrNonce
	}
	contract := tx.GetTxBody().GetContract()
	if contract == param.Token {
		return a.fromTokenChange(tx, blockHeight)
	} else {
		return a.fromCoinChange(tx, blockHeight)
	}
}

// Change the primary account status of one party to the transaction transfer
func (a *Account) fromTokenChange(tx ITransaction, blockHeight uint64) error {
	amount := tx.GetTxBody().GetAmount()
	fees := tx.GetFees()
	if !a.IsExist() {
		a.Address = tx.From()
	}
	coinAccount, ok := a.Coins.Get(param.Token.String())
	if !ok {
		return errors.New("account is not exist")
	}
	if coinAccount.Balance < amount {
		return ErrNotEnoughBalance
	}
	if a.Nonce+1 != tx.GetNonce() {
		return ErrNonce
	}

	coinAccount.Balance -= amount
	coinAccount.LockedIn += amount
	a.Coins.Set(coinAccount)
	a.Nonce = tx.GetNonce()
	a.Time = tx.GetTime()
	a.JournalIn.Add(tx, blockHeight, param.Token, amount-fees, fees)
	return nil
}

// Change the status of the secondary account of the transaction transfer party.
// The transaction of the secondary account needs to consume the fee of the
// primary account.
func (a *Account) fromCoinChange(tx ITransaction, blockHeight uint64) error {
	fees := tx.GetFees()
	txBody := tx.GetTxBody()
	amount := txBody.GetAmount()
	contract := txBody.GetContract()
	tokenAccount, ok := a.Coins.Get(param.Token.String())
	if !ok {
		return errors.New("account is not exist")
	}
	if tokenAccount.Balance < fees {
		return ErrNotEnoughFees
	}

	coinAccount, ok := a.Coins.Get(contract.String())
	if !ok {
		return errors.New("account is not exist")
	}
	if coinAccount.Balance < amount {
		return ErrNotEnoughBalance
	}

	tokenAccount.Balance -= fees
	tokenAccount.LockedIn += fees
	coinAccount.Balance -= amount
	coinAccount.LockedIn += amount
	a.Coins.Set(tokenAccount)
	a.Coins.Set(coinAccount)
	a.Nonce = tx.GetNonce()
	a.Time = tx.GetTime()
	a.JournalIn.Add(tx, blockHeight, contract, amount, fees)
	return nil
}

// Change of contract information
func (a *Account) fromContractChange(tx ITransaction, blockHeight uint64) error {
	fees := tx.GetFees()
	tokenAccount, ok := a.Coins.Get(param.Token.String())
	if !ok {
		return errors.New("account is not exist")
	}
	if tokenAccount.Balance < fees {
		return ErrNotEnoughFees
	}
	tokenAccount.Balance -= fees
	tokenAccount.LockedIn += fees
	a.Coins.Set(tokenAccount)
	a.Nonce = tx.GetNonce()
	a.Time = tx.GetTime()
	a.JournalIn.Add(tx, blockHeight, param.Token, 0, fees)
	return nil
}

// Change of contract information
func (a *Account) toContractChange(tx ITransaction, blockHeight uint64) error {
	txBody := tx.GetTxBody()
	amount := txBody.GetAmount()
	contract := txBody.GetContract()
	coinAccount, ok := a.Coins.Get(contract.String())
	if ok {
		coinAccount.LockedOut += amount
	} else {
		coinAccount = &CoinAccount{
			Contract:  contract.String(),
			Balance:   0,
			LockedOut: amount,
			LockedIn:  0,
		}
	}

	a.Coins.Set(coinAccount)
	a.JournalOut.Add(txBody.GetContract(), txBody.GetAmount(), blockHeight)
	return nil
}

// Change the status of the recipient of the transaction
func (a *Account) ToChange(tx ITransaction, blockHeight uint64) error {
	txBody := tx.GetTxBody()
	if !a.IsExist() {
		a.Address = txBody.ToAddress()
	}
	if tx.GetTxType() == ContractTransaction {
		return a.toContractChange(tx, blockHeight)
	}

	amount := txBody.GetAmount()
	if txBody.GetContract().IsEqual(param.Token) {
		amount = amount - tx.GetFees()
	}

	coinAccount, ok := a.Coins.Get(txBody.GetContract().String())
	if ok {
		coinAccount.LockedOut += amount
	} else {
		coinAccount = &CoinAccount{
			Contract:  txBody.GetContract().String(),
			Balance:   0,
			LockedIn:  0,
			LockedOut: amount,
		}
	}
	a.Coins.Set(coinAccount)
	a.JournalOut.Add(txBody.GetContract(), amount, blockHeight)
	return nil
}

func (a *Account) FeesChange(fees, blockHeight uint64) {
	if !a.IsExist() {
		a.Address = param.FeeAddress
	}
	coinAccount, ok := a.Coins.Get(param.Token.String())
	if ok {
		coinAccount.LockedOut += fees
	} else {
		coinAccount = &CoinAccount{
			Contract:  param.Token.String(),
			Balance:   0,
			LockedIn:  0,
			LockedOut: fees,
		}
	}
	a.Coins.Set(coinAccount)
	a.JournalOut.Add(param.Token, fees, blockHeight)
}

// To verify the transaction status, the nonce value of the transaction
// must be greater than the nonce value of the account of the transferring
// party.
func (a *Account) VerifyTxState(tx ITransaction) error {
	if !a.IsExist() {
		a.Address = tx.GetTxBody().ToAddress()
	}

	/*	if tx.GetTime() <= a.Time {
			return ErrTime
		}
	*/
	if tx.GetNonce() <= a.Nonce {
		return ErrTxNonceRepeat
	}

	// The nonce value cannot be greater than the
	// maximum number of address transactions
	if tx.GetNonce() > a.Nonce+param.MaxAddressTxs {
		return ErrTooBigNonce
	}

	// Verify the balance of the token
	switch tx.GetTxType() {
	case NormalTransaction:
		if tx.GetTxBody().GetContract() == param.Token {
			return a.verifyTokenTxBalance(tx)
		} else {
			return a.verifyCoinTxBalance(tx)
		}
	case ContractTransaction:
		return a.verifyFees(tx)
	default:
		if tx.GetTxBody().GetAmount() != 0 {
			return ErrTxAmount
		}
		return a.verifyFees(tx)
	}

	return nil
}

// Verify the account balance of the primary transaction, the transaction
// value and transaction fee cannot be greater than the balance.
func (a *Account) verifyTokenTxBalance(tx ITransaction) error {
	tokenAccount, ok := a.Coins.Get(param.Token.String())
	if !ok {
		return ErrNotEnoughBalance
	} else if tokenAccount.Balance < tx.GetTxBody().GetAmount() {
		return ErrNotEnoughBalance
	}
	return nil
}

// Verify the account balance of the secondary transaction, the transaction
// value cannot be greater than the balance.
func (a *Account) verifyCoinTxBalance(tx ITransaction) error {
	txBody := tx.GetTxBody()
	if err := a.verifyFees(tx); err != nil {
		return err
	}

	coinAccount, ok := a.Coins.Get(txBody.GetContract().String())
	if !ok {
		return ErrNotEnoughBalance
	} else if coinAccount.Balance < txBody.GetAmount() {
		return ErrNotEnoughBalance
	}
	return nil
}

// Verification fee
func (a *Account) verifyFees(tx ITransaction) error {
	tokenAccount, ok := a.Coins.Get(param.Token.String())
	if !ok {
		return ErrNotEnoughFees
	} else if tokenAccount.Balance < tx.GetFees() {
		return ErrNotEnoughFees
	}
	return nil
}

// The current nonce value of the block transaction must be the
// nonce + 1 of the sender's account.
func (a *Account) VerifyNonce(nonce uint64) error {
	if nonce != a.Nonce+1 {
		return fmt.Errorf("the nonce value must be %d", a.Nonce+1)
	}
	return nil
}

func (a *Account) GetAddress() hasharry.Address {
	return a.Address
}

func (a *Account) GetBalance(contract string) uint64 {
	coinAccount, ok := a.Coins.Get(contract)
	if ok {
		return coinAccount.Balance
	}
	return 0
}

func (a *Account) GetLockedIn(contract string) uint64 {
	coinAccount, ok := a.Coins.Get(contract)
	if ok {
		return coinAccount.LockedIn
	}
	return 0
}

func (a *Account) GetLockedOut(contract string) uint64 {
	coinAccount, ok := a.Coins.Get(contract)
	if ok {
		return coinAccount.LockedOut
	}
	return 0
}

func (a *Account) GetNonce() uint64 {
	return a.Nonce
}

func (a *Account) GetTime() uint64 {
	return a.Time
}

func (a *Account) GetConfirmedHeight() uint64 {
	return a.ConfirmedHeight
}

func (a *Account) GetConfirmedNonce() uint64 {
	return a.ConfirmedNonce
}

func (a *Account) GetConfirmedTime() uint64 {
	return a.ConfirmedTime
}

// Determine whether the account is in the initial state
func (a *Account) IsEmpty() bool {
	if a.Nonce != 0 {
		return false
	}
	if !a.JournalOut.IsEmpty() {
		return false
	}
	if !a.JournalIn.IsEmpty() {
		return false
	}
	for _, coin := range *a.Coins {
		if coin.Balance != 0 || coin.LockedIn != 0 || coin.LockedOut != 0 {
			return false
		}
	}
	return true
}

type CoinAccount struct {
	Contract  string
	Balance   uint64
	LockedIn  uint64
	LockedOut uint64
}

// List of secondary accounts
type Coins []*CoinAccount

func (c *Coins) Get(contract string) (*CoinAccount, bool) {
	for _, coin := range *c {
		if coin.Contract == contract {
			return coin, true
		}
	}
	return &CoinAccount{}, false
}

func (c *Coins) Set(newCoin *CoinAccount) {
	for i, coin := range *c {
		if coin.Contract == newCoin.Contract {
			(*c)[i] = newCoin
			return
		}
	}
	*c = append(*c, newCoin)
}

// Account transfer log
type journalIn struct {
	Ins *TxInList
}

func newJournalIn() *journalIn {
	return &journalIn{Ins: &TxInList{}}
}

func (j *journalIn) Add(tx ITransaction, height uint64, contract hasharry.Address, amount uint64, fees uint64) {
	j.Ins.Set(&txIn{
		Contract: contract.String(),
		Amount:   amount,
		Fees:     fees,
		Nonce:    tx.GetNonce(),
		Time:     tx.GetTime(),
		Height:   height,
	})
}

func (j *journalIn) Get(height uint64) *txIn {
	in, ok := j.Ins.Get(height)
	if ok {
		return in
	}
	return nil
}

func (j *journalIn) Remove(height uint64) uint64 {
	tx, _ := j.Ins.Get(height)
	j.Ins.Remove(height)
	return tx.Amount
}

func (j *journalIn) IsExist(height uint64) bool {
	for _, txIn := range *j.Ins {
		if txIn.Height >= height {
			return true
		}
	}
	return false
}

func (j *journalIn) GetJournalIns(confirmedHeight uint64) []*txIn {
	txIns := make([]*txIn, 0)
	for _, txIn := range *j.Ins {
		if txIn.Height <= confirmedHeight {
			txIns = append(txIns, txIn)
		}
	}
	return txIns
}

func (j *journalIn) Amount() map[string]uint64 {
	amounts := map[string]uint64{}
	for _, txIn := range *j.Ins {
		_, ok := amounts[txIn.Contract]
		if ok {
			amounts[txIn.Contract] += txIn.Amount
		} else {
			amounts[txIn.Contract] = txIn.Amount
		}
	}
	return amounts
}

func (j *journalIn) IsEmpty() bool {
	if j.Ins == nil || len(*j.Ins) == 0 {
		return true
	}
	return false
}

type txIn struct {
	Contract string
	Amount   uint64
	Fees     uint64
	Nonce    uint64
	Time     uint64
	Height   uint64
}

type TxInList []*txIn

func (t *TxInList) Get(height uint64) (*txIn, bool) {
	for _, txIn := range *t {
		if txIn.Height == height {
			return txIn, true
		}
	}
	return &txIn{}, false
}

func (t *TxInList) Set(txIn *txIn) {
	for i, in := range *t {
		if in.Height == txIn.Height {
			(*t)[i] = txIn
			return
		}
	}
	*t = append(*t, txIn)
}

func (t *TxInList) Remove(height uint64) {
	for i, in := range *t {
		if in.Height == height {
			*t = append((*t)[0:i], (*t)[i+1:]...)
			return
		}
	}
}

// Account transfer log
type journalOut struct {
	Outs *OutList
}

func newJournalOut() *journalOut {
	return &journalOut{Outs: &OutList{}}
}

func (j *journalOut) Add(contract hasharry.Address, amount, height uint64) {
	out, ok := j.Outs.Get(height, contract.String())
	if ok {
		out.Amount += amount
	} else {
		out = &OutAmount{}
		out.Amount = amount
		out.Height = height
		out.Contract = contract.String()
	}
	j.Outs.Set(out)
}

func (j *journalOut) Get(height uint64, contract string) *OutAmount {
	txOut, ok := j.Outs.Get(height, contract)
	if ok {
		return txOut
	}
	return &OutAmount{0, "", 0}
}

func (j *journalOut) IsExist(height uint64) bool {
	for _, out := range *j.Outs {
		if out.Height >= height {
			return true
		}
	}
	return false
}

func (j *journalOut) Remove(height uint64, contract string) *OutAmount {
	return j.Outs.Remove(height, contract)
}

func (j *journalOut) GetJournalOuts(confirmedHeight uint64) map[string]*OutAmount {
	txOuts := make(map[string]*OutAmount)
	for _, out := range *j.Outs {
		if out.Height <= confirmedHeight {
			key := fmt.Sprintf("%s_%d", out.Contract, out.Height)
			txOuts[key] = out
		}
	}
	return txOuts
}

func (j *journalOut) IsEmpty() bool {
	if j.Outs == nil || len(*j.Outs) == 0 {
		return true
	}
	return false
}

type OutAmount struct {
	Amount   uint64
	Contract string
	Height   uint64
}

type OutList []*OutAmount

func (o *OutList) Get(height uint64, contract string) (*OutAmount, bool) {
	for _, out := range *o {
		if out.Height == height && out.Contract == contract {
			return out, true
		}
	}
	return &OutAmount{}, false
}

func (o *OutList) Set(outAmount *OutAmount) {
	for i, out := range *o {
		if out.Height == outAmount.Height && out.Contract == outAmount.Contract {
			(*o)[i] = outAmount
			return
		}
	}
	*o = append(*o, outAmount)
}

func (o *OutList) Remove(height uint64, contract string) *OutAmount {
	for i, out := range *o {
		if out.Height == height && out.Contract == contract {
			*o = append((*o)[0:i], (*o)[i+1:]...)
			return out
		}
	}
	return nil
}
