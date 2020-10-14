package ecp256

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"errors"
	"math/big"
)

func isOdd(a *big.Int) bool {
	return a.Bit(0) == 1
}

// NewPublicKey instantiates a new public key with the given X,Y coordinates.
func NewPublicKey(x *big.Int, y *big.Int) *ecdsa.PublicKey {
	return &ecdsa.PublicKey{elliptic.P256(), x, y}
}

func UnmarshalPubkey(pub []byte) (*ecdsa.PublicKey, error) {
	x, y := elliptic.Unmarshal(elliptic.P256(), pub)
	if x == nil {
		return nil, errors.New("unmarshal failed")
	}

	return &ecdsa.PublicKey{Curve: elliptic.P256(), X: x, Y: y}, nil
}

func FromECDSAPub(pub *ecdsa.PublicKey) []byte {
	if pub == nil || pub.X == nil || pub.Y == nil {
		return nil
	}
	return elliptic.Marshal(elliptic.P256(), pub.X, pub.Y)

}
