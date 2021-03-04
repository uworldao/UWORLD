package config

import (
	"encoding/hex"
	"fmt"
	"github.com/uworldao/UWORLD/common/hasharry"
	"github.com/uworldao/UWORLD/common/keystore"
	"github.com/uworldao/UWORLD/common/utils"
	"github.com/uworldao/UWORLD/crypto/ecc/secp256k1"
	"github.com/uworldao/UWORLD/crypto/mnemonic"
	"github.com/uworldao/UWORLD/ut"
	"io/ioutil"
	"os"
)

type NodePrivate struct {
	Address    hasharry.Address
	PrivateKey *secp256k1.PrivateKey
	Mnemonic   string
}

func CreateNewNodePrivate(net string) (*NodePrivate, error) {
	entropy, err := mnemonic.Entropy()
	if err != nil {
		return nil, err
	}
	mnemonicStr, err := mnemonic.Mnemonic(entropy)
	if err != nil {
		return nil, err
	}
	key, err := mnemonic.MnemonicToEc(mnemonicStr)
	if err != nil {
		return nil, err
	}
	address, err := ut.GenerateAddress(net, key.PubKey())
	if err != nil {
		return nil, err
	}
	return &NodePrivate{hasharry.StringToAddress(address), key, mnemonicStr}, nil
}

func LoadNodePrivate(file string, key string) (*NodePrivate, error) {
	if !utils.IsExist(file) {
		return nil, fmt.Errorf("%s is not exsists", file)
	}
	j, err := keystore.ReadJson(file)
	if err != nil {
		return nil, fmt.Errorf("read json file %s failed! %s", file, err.Error())
	}
	privJson, err := keystore.DecryptPrivate([]byte(key), j)
	if err != nil {
		return nil, fmt.Errorf("decrypt priavte failed! %s", err.Error())
	}
	privKey, err := secp256k1.ParseStringToPrivate(privJson.Private)
	if err != nil {
		return nil, fmt.Errorf("parse priavte failed! %s", err.Error())
	}
	return &NodePrivate{hasharry.StringToAddress(j.Address), privKey, privJson.Mnemonic}, nil
}

func GenerateNodeKeyFile(file string) (*secp256k1.PrivateKey, error) {
	private, err := secp256k1.GeneratePrivateKey()
	if err != nil {
		return nil, err
	}
	k := hex.EncodeToString(private.D.Bytes())
	return private, ioutil.WriteFile(file, []byte(k), os.ModePerm)
}
