package database

import (
	"errors"
	"fmt"
	"github.com/boltdb/bolt"
)

const timeout = 10

type Base struct {
	Path string
	Db   *bolt.DB
}

func (b *Base) Open() error {
	if db, err := bolt.Open(b.Path, 0600, &bolt.Options{ReadOnly: false, Timeout: timeout}); err != nil {
		return fmt.Errorf("open database failed! %s", err.Error())
	} else {
		b.Db = db
	}
	return nil
}

func (b *Base) Close() {
	b.Close()
}

func (b *Base) CreateBucket(name string) error {
	err := b.Db.Batch(func(tx *bolt.Tx) error {
		if _, err := tx.CreateBucketIfNotExists([]byte(name)); err != nil {
			return err
		}
		return nil
	})
	return err
}

func (b *Base) UpdateValue(bucket string, key []byte, value []byte) error {
	err := b.Db.Update(func(tx *bolt.Tx) error {
		if tb := tx.Bucket([]byte(bucket)); tb == nil {
			return errors.New(fmt.Sprintf("bucket %s is not exsit", bucket))
		} else {
			return tb.Put(key, value)
		}
	})
	return err
}

func (b *Base) DeleteKey(bucket string, key []byte) error {
	err := b.Db.Update(func(tx *bolt.Tx) error {
		if tb := tx.Bucket([]byte(bucket)); tb == nil {
			return errors.New(fmt.Sprintf("bucket %s is not exsit", bucket))
		} else {
			return tb.Delete(key)
		}
	})
	return err
}

func (b *Base) GetValue(bucket string, key []byte) ([]byte, error) {
	var rs []byte
	err := b.Db.View(func(tx *bolt.Tx) error {
		if tb := tx.Bucket([]byte(bucket)); tb == nil {
			return errors.New(fmt.Sprintf("bucket %s is not exsit", bucket))
		} else {
			rs = tb.Get(key)
			if rs == nil {
				return errors.New("the key is not exsit")
			}
			return nil
		}
	})
	if err != nil {
		return nil, err
	}
	return rs, nil
}

func (b *Base) GetForeach(bucket string, f func(k []byte, v []byte) error) error {
	err := b.Db.View(func(tx *bolt.Tx) error {
		if tb := tx.Bucket([]byte(bucket)); tb == nil {
			return errors.New(fmt.Sprintf("bucket %s is not exsit", bucket))
		} else {
			err := tb.ForEach(f)
			if err != nil {
				return err
			}
			return nil
		}
	})
	return err
}
