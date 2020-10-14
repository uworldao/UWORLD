package triedb

import (
	"github.com/jhdriver/UWORLD/database/leveldb"
)

// Implementation of tire tree storage
type TrieDB struct {
	db *leveldb.Base
}

func NewTrieDB(path string) *TrieDB {
	return &TrieDB{&leveldb.Base{path, nil}}
}

func (s *TrieDB) Open() error {
	err := s.db.Open()
	if err != nil {
		return err
	}
	return s.initBucket()
}

func (s *TrieDB) Close() error {
	return s.db.Close()
}

func (s *TrieDB) initBucket() error {

	return nil
}

func (s *TrieDB) Put(key, value []byte) error {
	return s.db.UpdateValue(key, value)
}

func (s *TrieDB) Get(key []byte) (value []byte, err error) {
	return s.db.GetValue(key)
}

func (s *TrieDB) Has(key []byte) (bool, error) {
	_, err := s.db.GetValue(key)
	if err != nil {
		return false, err
	}
	return true, nil
}

func (s *TrieDB) CreateBucket(bucket string) error {
	if err := s.db.CreateBucket(bucket); err != nil {
		return err
	}
	return nil
}

func (s *TrieDB) PutToBucket(bucket string, key, value []byte) error {
	dbKey := leveldb.GetKey(bucket, key)
	return s.db.UpdateValue(dbKey, value)
}

func (s *TrieDB) GetFromBucket(bucket string, key []byte) ([]byte, error) {
	dbKey := leveldb.GetKey(bucket, key)
	return s.db.GetValue(dbKey)
}
