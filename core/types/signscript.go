package types

import (
	"github.com/uworldao/UWORLD/common/hasharry"
	"github.com/uworldao/UWORLD/crypto/ecc/secp256k1"
	"github.com/uworldao/UWORLD/ut"
)

// Signature information, including the result of the
// signature and the public key.
type SignScript struct {
	Signature []byte `json:"signature"`
	PubKey    []byte `json:"pubkey"`
}

// Sign the hash with the private key
func Sign(key *secp256k1.PrivateKey, hash hasharry.Hash) (*SignScript, error) {
	signature, err := key.Sign(hash.Bytes())
	if err != nil {
		return nil, err
	}
	return &SignScript{signature.Serialize(), key.PubKey().SerializeCompressed()}, nil
}

// Verify signature by hash and signature result
func Verify(hash hasharry.Hash, signScript *SignScript) bool {
	if signScript == nil || signScript.PubKey == nil || signScript.Signature == nil {
		return false
	}
	pubkey, err := secp256k1.ParsePubKey(signScript.PubKey)
	if err != nil {
		return false
	}
	signature, err := secp256k1.ParseSignature(signScript.Signature, secp256k1.S256())
	return signature.Verify(hash.Bytes(), pubkey)
}

// Verify whether the signers are consistent through the public key
func VerifySigner(network string, signer hasharry.Address, pubKey []byte) bool {
	key, err := secp256k1.ParsePubKey(pubKey)
	if err != nil {
		return false
	}
	generateAddress, err := ut.GenerateAddress(network, key)
	if err != nil {
		return false
	}
	return generateAddress == signer.String()
}
