package ecp256

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"io"
	"math/big"
)

// NewPrivateKey instantiates a new private key from a scalar encoded as a
// big integer.
func NewPrivateKey(d *big.Int) *ecdsa.PrivateKey {
	b := make([]byte, 0, PrivKeyBytesLen)
	dB := paddedAppend(PrivKeyBytesLen, b, d.Bytes())
	priv, _ := PrivKeyFromBytes(dB)
	return priv
}

func paddedAppend(size uint, dst, src []byte) []byte {
	for i := 0; i < int(size)-len(src); i++ {
		dst = append(dst, 0)
	}
	return append(dst, src...)
}

// PrivKeyFromBytes returns a private and public key for `curve' based on the
// private key passed as an argument as a byte slice.
func PrivKeyFromBytes(pk []byte) (*ecdsa.PrivateKey,
	*ecdsa.PublicKey) {
	curve := elliptic.P256()
	x, y := curve.ScalarBaseMult(pk)

	priv := &ecdsa.PrivateKey{
		PublicKey: ecdsa.PublicKey{
			Curve: elliptic.P256(),
			X:     x,
			Y:     y,
		},
		D: new(big.Int).SetBytes(pk),
	}

	return priv, &priv.PublicKey
}

// PrivKeyFromScalar is the same as PrivKeyFromBytes in secp256k1.
func PrivKeyFromScalar(s []byte) (*ecdsa.PrivateKey,
	*ecdsa.PublicKey) {
	return PrivKeyFromBytes(s)
}

// GeneratePrivateKey is a wrapper for ecdsa.GenerateKey that returns a PrivateKey
// instead of the normal ecdsa.PrivateKey.
func GeneratePrivateKey() (*ecdsa.PrivateKey, error) {
	key, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		return nil, err
	}
	return key, nil
}

// GenerateKey generates a key using a random number generator, returning
// the private scalar and the corresponding public key points.
func GenerateKey(rand io.Reader) (priv []byte, x,
	y *big.Int, err error) {
	key, err := ecdsa.GenerateKey(elliptic.P256(), rand)
	priv = key.D.Bytes()
	x = key.PublicKey.X
	y = key.PublicKey.Y

	return
}

// PrivKeyBytesLen defines the length in bytes of a serialized private key.
const PrivKeyBytesLen = 32
