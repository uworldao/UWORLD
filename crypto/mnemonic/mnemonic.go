package mnemonic

import (
	"encoding/hex"
	"github.com/jhdriver/UWORLD/crypto/bip32"
	"github.com/jhdriver/UWORLD/crypto/bip39"
	"github.com/jhdriver/UWORLD/crypto/ecc/secp256k1"
	"github.com/jhdriver/UWORLD/crypto/seed"
)

func Entropy() (string, error) {
	s, err := seed.GenerateSeed(seed.DefaultSeedBytes)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(s), nil
}

func Mnemonic(entropyStr string) (string, error) {
	entropy, err := hex.DecodeString(entropyStr)
	if err != nil {
		return "", err
	}
	return bip39.NewMnemonic(entropy)
}

func MnemonicToEc(mnemonic string) (*secp256k1.PrivateKey, error) {
	seed, err := bip39.NewSeedWithErrorChecking(mnemonic, "")
	if err != nil {
		return nil, err
	}
	masterKey, err := bip32.NewMasterKey(seed)
	if err != nil {
		return nil, err
	}
	return secp256k1.ParseStringToPrivate(hex.EncodeToString(masterKey.Key))
}
