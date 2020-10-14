package codec

import (
	"bytes"
	"encoding/binary"
	"encoding/gob"
)

func Uint16toBytes(u uint16) []byte {
	var buf = make([]byte, 2)
	binary.BigEndian.PutUint16(buf, u)
	return buf
}

func Uint32toBytes(u uint32) []byte {
	var buf = make([]byte, 4)
	binary.BigEndian.PutUint32(buf, u)
	return buf
}

func Uint64toBytes(u uint64) []byte {
	var buf = make([]byte, 8)
	binary.BigEndian.PutUint64(buf, u)
	return buf
}

func FromBytes(val []byte, obj interface{}) error {
	decoder := gob.NewDecoder(bytes.NewReader(val))
	err := decoder.Decode(obj)
	if err != nil {
		return err
	}
	return nil
}
