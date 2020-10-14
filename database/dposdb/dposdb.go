package dposdb

import (
	"bytes"
	"github.com/jhdriver/UWORLD/common/encode/rlp"
	hash2 "github.com/jhdriver/UWORLD/common/hasharry"
	"github.com/jhdriver/UWORLD/core/types"
	"github.com/jhdriver/UWORLD/crypto/hash"
	"github.com/jhdriver/UWORLD/database/triedb"
	"github.com/jhdriver/UWORLD/trie"
	"strconv"
)

const (
	dposBucket = "dposBucket"
)

type DPosStorage struct {
	trieDB   *triedb.TrieDB
	dposTrie *trie.Trie
}

func NewDPosStorage(path string) *DPosStorage {
	trieDB := triedb.NewTrieDB(path)
	return &DPosStorage{trieDB, nil}
}

func (c *DPosStorage) InitTrie(contractRoot hash2.Hash) error {
	contractTrie, err := trie.New(contractRoot, c.trieDB)
	if err != nil {
		return err
	}
	c.dposTrie = contractTrie
	return nil
}

func (c *DPosStorage) Commit() (hash2.Hash, error) {
	return c.dposTrie.Commit()
}

func (c *DPosStorage) RootHash() hash2.Hash {
	return c.dposTrie.Hash()
}

func (c *DPosStorage) Open() error {
	if err := c.trieDB.Open(); err != nil {
		return err
	}
	return c.trieDB.CreateBucket(dposBucket)
}

func (c *DPosStorage) Close() error {
	return c.trieDB.Close()
}

func (d *DPosStorage) SetTermWinners(term uint64, winners *types.Winners) error {
	termBytes, err := rlp.EncodeToBytes(term)
	if err != nil {
		return err
	}
	bytes, err := rlp.EncodeToBytes(winners)
	if err != nil {
		return err
	}
	d.trieDB.PutToBucket(dposBucket, termBytes, bytes)
	return nil
}

func (d *DPosStorage) GetTermWinners(term uint64) (*types.Winners, error) {
	var winners *types.Winners
	termBytes, err := rlp.EncodeToBytes(term)
	if err != nil {
		return nil, err
	}
	bytes, err := d.trieDB.GetFromBucket(dposBucket, termBytes)
	if err := rlp.DecodeBytes(bytes, &winners); err != nil {
		return nil, err
	}
	return winners, nil
}

func CandidatesHash() hash2.Hash {
	return hash.Hash([]byte("candidates"))
}

func (dps *DPosStorage) SetCandidate(newCan *types.Candidate) error {
	var candidates *types.Candidates
	keyHash := CandidatesHash().Bytes()
	bytes := dps.dposTrie.Get(keyHash)
	if err := rlp.DecodeBytes(bytes, &candidates); err != nil {
		candidates = types.NewCandidates()
	}
	candidates.Set(newCan)
	if bytes, err := rlp.EncodeToBytes(candidates); err != nil {
		return err
	} else {
		dps.dposTrie.Update(keyHash, bytes)
		return nil
	}
}

func (dps *DPosStorage) GetCandidates() (*types.Candidates, error) {
	var candidates *types.Candidates
	keyHash := CandidatesHash().Bytes()
	bytes := dps.dposTrie.Get(keyHash)
	if err := rlp.DecodeBytes(bytes, &candidates); err != nil {
		return nil, err
	}
	return candidates, nil
}

func (dps *DPosStorage) DeleteCandidate(can *types.Candidate) error {
	var candidates *types.Candidates
	keyHash := CandidatesHash().Bytes()
	bytes := dps.dposTrie.Get(keyHash)
	if err := rlp.DecodeBytes(bytes, &candidates); err != nil {
		return err
	}
	candidates.Remove(can)
	if bytes, err := rlp.EncodeToBytes(candidates); err != nil {
		return err
	} else {
		dps.dposTrie.Update(keyHash, bytes)
		return nil
	}
}

type Voter struct {
	Voter hash2.Address
	To    hash2.Address
}

type VoteList []Voter

func NewVoterInfo() *VoteList {
	return &VoteList{}
}

func (v *VoteList) Set(from, to hash2.Address) {
	newVoter := Voter{
		Voter: from,
		To:    to,
	}
	for i, voter := range *v {
		if voter.Voter.IsEqual(from) {
			(*v)[i] = newVoter
			return
		}
	}
	*v = append(*v, newVoter)
}

func (v *VoteList) GetVoters(addr hash2.Address) []hash2.Address {
	var voters = make([]hash2.Address, 0)
	for _, voter := range *v {
		if voter.To.IsEqual(addr) {
			voters = append(voters, voter.Voter)
		}
	}
	return voters
}

func VoteInfoHash() hash2.Hash {
	return hash.Hash([]byte("voter info"))
}

func (dps *DPosStorage) SetVoter(from, to hash2.Address) error {
	var voterInfo *VoteList
	var err error
	voterHash := VoteInfoHash().Bytes()
	bytes := dps.dposTrie.Get(voterHash)
	if err := rlp.DecodeBytes(bytes, &voterInfo); err != nil {
		voterInfo = NewVoterInfo()
	}
	voterInfo.Set(from, to)
	bytes, err = rlp.EncodeToBytes(voterInfo)
	if err != nil {
		return err
	}
	dps.dposTrie.Update(voterHash, bytes)
	return nil

}

func (dps *DPosStorage) GetCandidateVoters(addr hash2.Address) []hash2.Address {
	var voterInfo *VoteList
	voterHash := VoteInfoHash().Bytes()
	bytes := dps.dposTrie.Get(voterHash)
	if err := rlp.DecodeBytes(bytes, &voterInfo); err != nil {
		return []hash2.Address{}
	}
	return voterInfo.GetVoters(addr)
}

func ConfirmedHash() hash2.Hash {
	return hash.Hash([]byte("confirmed block hash"))
}

func (dps *DPosStorage) GetConfirmedBlockHash() (hash2.Hash, error) {
	bytes := dps.dposTrie.Get(ConfirmedHash().Bytes())
	return hash2.BytesToHash(bytes), nil
}

func (dps *DPosStorage) SetConfirmedBlockHash(hash hash2.Hash) {
	dps.dposTrie.Update(ConfirmedHash().Bytes(), hash.Bytes())
}

func (dps *DPosStorage) GetTermWinnerMintCnt(term uint64, address hash2.Address) (uint64, error) {
	hash := termWinnerMintCntHash(term, address).Bytes()
	bytes := dps.dposTrie.Get(hash)
	return strconv.ParseUint(string(bytes), 10, 64)
}
func (dps *DPosStorage) SetTermWinnerMintCnt(term uint64, address hash2.Address) {
	hash := termWinnerMintCntHash(term, address).Bytes()
	cnt, err := dps.GetTermWinnerMintCnt(term, address)
	if err != nil {
		cnt = 0
	}
	cnt++
	dps.dposTrie.Update(hash, []byte(strconv.FormatUint(cnt, 10)))
}

func termWinnerMintCntHash(term uint64, address hash2.Address) hash2.Hash {
	bytes := bytes.Join([][]byte{[]byte(strconv.FormatUint(term, 10)), address.Bytes()}, []byte{})
	return hash.Hash(bytes)
}
