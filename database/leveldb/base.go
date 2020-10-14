package leveldb

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/btcsuite/goleveldb/leveldb"
	"github.com/btcsuite/goleveldb/leveldb/opt"
	"github.com/btcsuite/goleveldb/leveldb/util"
)

type Base struct {
	Path string
	Db   *leveldb.DB
}

func (b *Base) Open() error {
	var err error
	opts := &opt.Options{
		OpenFilesCacheCapacity: 16,
		Strict:                 opt.DefaultStrict,
		Compression:            opt.NoCompression,
		BlockCacheCapacity:     8 * opt.MiB,
		WriteBuffer:            4 * opt.MiB,
	}
	if b.Db, err = leveldb.OpenFile(b.Path, opts); err != nil {
		if b.Db, err = leveldb.RecoverFile(b.Path, nil); err != nil {
			return errors.New(fmt.Sprintf(`err while recoverfile %s : %s`, b.Path, err.Error()))
		}

	}
	return nil
}

func (b *Base) Close() error {
	return b.Db.Close()
}

func (b *Base) CreateBucket(name string) error {
	return nil
}

func (b *Base) UpdateValue(key []byte, value []byte) error {
	return b.Db.Put(key, value, nil)
}

func (b *Base) DeleteKey(key []byte) error {
	return b.Db.Delete(key, nil)
}

func (b *Base) GetValue(key []byte) ([]byte, error) {
	return b.Db.Get(key, nil)
}

func (b *Base) ClearBucket(bucket string) {
	rs := b.Foreach(bucket)
	for key, _ := range rs {
		b.Db.Delete([]byte(key), nil)
	}
}

func (b *Base) Foreach(bucket string) map[string][]byte {
	rs := make(map[string][]byte)
	iter := b.Db.NewIterator(util.BytesPrefix(bytes.Join([][]byte{[]byte(bucket), []byte("-")}, []byte{})), nil)
	defer iter.Release()

	// Iter will affect RLP decoding and reallocate memory to value
	for iter.Next() {
		value := make([]byte, len(iter.Value()))
		copy(value, iter.Value())
		rs[string(iter.Key())] = value
	}
	return rs
}

func GetKey(bucket string, key []byte) []byte {
	return bytes.Join([][]byte{
		[]byte(bucket + "-"), key}, []byte{})
}
