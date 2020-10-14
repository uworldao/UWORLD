package codec

import (
	"bytes"
	"encoding/binary"
	"encoding/gob"
)

func BytesToUint16(val []byte) uint16 {
	return binary.BigEndian.Uint16(val)
}

func BytesToUint32(val []byte) uint32 {
	return binary.BigEndian.Uint32(val)
}

func BytesToUint64(val []byte) uint64 {
	return binary.BigEndian.Uint64(val)
}

func ToBytes(s interface{}) ([]byte, error) {
	var result bytes.Buffer
	encoder := gob.NewEncoder(&result)
	err := encoder.Encode(s)
	if err != nil {
		return nil, err
	}
	return result.Bytes(), nil
}
