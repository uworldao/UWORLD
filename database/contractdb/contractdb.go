package contractdb

import (
	"github.com/jhdriver/UWORLD/common/codec"
	"github.com/jhdriver/UWORLD/common/hasharry"
	"github.com/jhdriver/UWORLD/core/types"
	"github.com/jhdriver/UWORLD/database/triedb"
	"github.com/jhdriver/UWORLD/trie"
)

type ContractStorage struct {
	trieDB       *triedb.TrieDB
	contractTrie *trie.Trie
}

func NewContractStorage(path string) *ContractStorage {
	trieDB := triedb.NewTrieDB(path)
	return &ContractStorage{trieDB, nil}
}

func (c *ContractStorage) InitTrie(contractRoot hasharry.Hash) error {
	contractTrie, err := trie.New(contractRoot, c.trieDB)
	if err != nil {
		return err
	}
	c.contractTrie = contractTrie
	return nil
}

func (c *ContractStorage) Commit() (hasharry.Hash, error) {
	return c.contractTrie.Commit()
}

func (c *ContractStorage) RootHash() hasharry.Hash {
	return c.contractTrie.Hash()
}

func (c *ContractStorage) Open() error {
	return c.trieDB.Open()
}

func (c *ContractStorage) Close() error {
	return c.trieDB.Close()
}

func (c *ContractStorage) GetContractState(contractAddr string) *types.Contract {
	contract := types.NewContract()
	bytes := c.contractTrie.Get([]byte(contractAddr))
	err := codec.FromBytes(bytes, &contract)
	if err != nil {
		return nil
	}
	return contract
}

func (c *ContractStorage) SetContractState(contract *types.Contract) {
	bytes, err := codec.ToBytes(contract)
	if err != nil {
		return
	}
	c.contractTrie.Update([]byte(contract.Contract), bytes)
}
