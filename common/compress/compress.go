package compress

import (
	"bytes"
	"github.com/ulikunitz/xz"
)

func CompressBytes(in []byte) ([]byte, error) {
	var buf bytes.Buffer
	w, err := xz.NewWriter(&buf)
	if err != nil {
		return nil, err
	}
	_, err = w.Write(in)
	if err != nil {
		return nil, err
	}
	if err := w.Close(); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func DecompressBytes(in []byte) ([]byte, error) {
	var buf bytes.Buffer

	if _, err := buf.Write(in); err != nil {
		return nil, err
	}
	_, err := xz.NewReader(&buf)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
