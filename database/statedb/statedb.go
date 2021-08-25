package statedb

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/uworldao/UWORLD/common/encode/rlp"
	"github.com/uworldao/UWORLD/common/hasharry"
	"github.com/uworldao/UWORLD/core/types"
	"github.com/uworldao/UWORLD/database/triedb"
	"github.com/uworldao/UWORLD/trie"
)

type StateStorage struct {
	trieDB    *triedb.TrieDB
	stateTrie *trie.Trie
}

func NewStateStorage(path string) *StateStorage {
	trieDB := triedb.NewTrieDB(path)

	return &StateStorage{trieDB, nil}
}

func (s *StateStorage) InitTrie(stateRoot hasharry.Hash) error {
	stateTrie, err := trie.New(stateRoot, s.trieDB)
	if err != nil {
		return err
	}
	s.stateTrie = stateTrie
	return nil
}

func (s *StateStorage) Open() error {
	return s.trieDB.Open()
}

func (s *StateStorage) Close() error {
	return s.trieDB.Close()
}

func (s *StateStorage) GetAccounts() []types.IAccount{
	iter := s.stateTrie.PrefixIterator([]byte{})
	var accounts []types.IAccount
	for iter.Next(true) {
		if iter.Leaf() {
			account := types.NewAccount()
			if err := rlp.DecodeBytes(iter.LeafBlob(), &account);err != nil{
				continue
			}
			accounts = append(accounts, account)
		}
	}
	return accounts
}

func (s *StateStorage) GetAccountState(stateKey hasharry.Address) types.IAccount {
	account := types.NewAccount()
	bytes := s.stateTrie.Get(stateKey.Bytes())
	err := rlp.DecodeBytes(bytes, &account)
	if err != nil {
		return types.NewAccount()
	}
	return account
}

func (s *StateStorage) SetAccountState(account types.IAccount) {
	bytes, err := rlp.EncodeToBytes(account.(*types.Account))
	if err != nil {
		return
	}
	s.stateTrie.Update(account.StateKey().Bytes(), bytes)
}

func (s *StateStorage) GetAccountBalance(stateKey hasharry.Address, contract string) uint64 {
	account := types.NewAccount()
	bytes := s.stateTrie.Get(stateKey.Bytes())
	err := rlp.DecodeBytes(bytes, &account)
	if err != nil {
		return 0
	}
	return account.GetBalance(contract)
}

func (s *StateStorage) GetAccountNonce(stateKey hasharry.Address) uint64 {
	account := types.NewAccount()
	bytes := s.stateTrie.Get(stateKey.Bytes())
	err := rlp.DecodeBytes(bytes, &account)
	if err != nil {
		return 0
	}
	return account.GetNonce()
}

func (s *StateStorage) DeleteAccount(stateKey hasharry.Address) {
	s.stateTrie.Delete(stateKey.Bytes())
}

func (s *StateStorage) Commit() (hasharry.Hash, error) {
	return s.stateTrie.Commit()
}

func (s *StateStorage) RootHash() hasharry.Hash {
	return s.stateTrie.Hash()
}

func (s *StateStorage) Print() {
	account1 := s.GetAccountState(hasharry.StringToAddress("3ajPAQyobsVaDVAwhpeLo8vouirRrEJvDqZ2"))
	account2 := s.GetAccountState(hasharry.StringToAddress("3ajDJUnMYDyzXLwefRfNp7yLcdmg3ULb9ndQ"))
	account3 := s.GetAccountState(hasharry.StringToAddress("3ajF4MdbBYE2UPESEyhQbdUj2Y28CNwGDCWA"))

	rpca1 := account1.(*types.Account)
	rpca2 := account2.(*types.Account)
	rpca3 := account3.(*types.Account)
	jsonAcc1, _ := json.Marshal(rpca1)
	jsonAcc2, _ := json.Marshal(rpca2)
	jsonAcc3, _ := json.Marshal(rpca3)

	var str1 bytes.Buffer
	var str2 bytes.Buffer
	var str3 bytes.Buffer
	fmt.Print(s.stateTrie.Hash().String())
	_ = json.Indent(&str1, jsonAcc1, "", "    ")
	fmt.Println(str1.String())
	_ = json.Indent(&str2, jsonAcc2, "", "    ")
	fmt.Println(str2.String())
	_ = json.Indent(&str3, jsonAcc3, "", "    ")
	fmt.Println(str3.String())
}
