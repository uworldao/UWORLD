package keystore

import (
	"bytes"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/uworldao/UWORLD/common/utils"
	"github.com/uworldao/UWORLD/crypto/aes"
	"github.com/uworldao/UWORLD/crypto/ecc/secp256k1"
	"github.com/uworldao/UWORLD/ut"
	"io"
	"io/ioutil"
	"os"
	"strings"
)

const boot = "/.boot"

type Json struct {
	Address string  `json:"address"`
	Crypto  *crypto `json:"crypto"`
	P2pId   string  `json:"p2pid"`
}

type Private struct {
	Private  string `json:"private"`
	Mnemonic string `json:"mnemonic"`
}

type crypto struct {
	Cipher             string `json:"cipher"`
	CipherText         string `json:"ciphertext"`
	MnemonicCipherText string `json:"mnemonic_ciphertext"`
	Salt               string `json:"salt"`
}

func GenerateKeyJson(net string, dir string, private *secp256k1.PrivateKey, mnemonicStr string, passWd []byte) (*Json, error) {
	if !utils.IsExist(dir) {
		if err := os.Mkdir(dir, os.ModePerm); err != nil {
			return nil, fmt.Errorf("create folder %s failed! %s", dir, err.Error())
		}
	}

	j, err := PrivateToJson(net, private, mnemonicStr, passWd)
	if err := saveJson(dir, j); err != nil {
		return nil, err
	}
	return j, err
}

func PrivateToJson(net string, priv *secp256k1.PrivateKey, mnemonicStr string, passWd []byte) (*Json, error) {
	salt, err := getRandSalt(32 - len(passWd))
	if err != nil {
		return nil, fmt.Errorf("get rand salt failed! %s", err.Error())
	}

	address, err := ut.GenerateAddress(net, priv.PubKey())
	if err != nil {
		return nil, fmt.Errorf("generate address failed! %s", err.Error())
	}
	cipherText, ok := aes.AESCFBEncrypt(bytes.Join([][]byte{passWd, salt}, []byte{}), priv.String())
	if !ok {
		return nil, errors.New("aes encrypt failed")
	}
	mnemonicText, ok := aes.AESCFBEncrypt(bytes.Join([][]byte{passWd, salt}, []byte{}), mnemonicStr)
	if !ok {
		return nil, errors.New("aes encrypt failed")
	}
	j := &Json{
		Address: address,
		Crypto: &crypto{
			Cipher:             "aes-128-cfb",
			CipherText:         hex.EncodeToString(cipherText),
			MnemonicCipherText: hex.EncodeToString(mnemonicText),
			Salt:               hex.EncodeToString(salt),
		},
	}
	return j, nil
}

func ReadAllAccount(dir string) ([]string, error) {
	accountList, err := utils.ReadLine(dir + boot)
	if err != nil {
		return nil, errors.New("no account exists")
	}
	rs := make([]string, len(accountList))
	for i, filePath := range accountList {
		rs[i] = strings.Replace(filePath, ".json", "", 1)
	}
	return rs, nil
}

func ReadAddressJson(dir, addr string) (*Json, error) {
	return ReadJson(dir + "/" + addr + ".json")
}

func DecryptPrivate(passWd []byte, j *Json) (*Private, error) {
	var mnemonicStr string
	var ok bool
	salt, err := hex.DecodeString(j.Crypto.Salt)
	if err != nil {
		return nil, fmt.Errorf("decode salt failed! %s", err.Error())
	}
	cipherText, err := hex.DecodeString(j.Crypto.CipherText)
	if err != nil {
		return nil, fmt.Errorf("decode cipherText failed! %s", err.Error())
	}
	privKeyStr, ok := aes.AESCFBDecrypt(bytes.Join([][]byte{passWd, salt}, []byte{}), cipherText)
	if !ok {
		return nil, errors.New("aes decrypt failed")
	}
	privKey, err := secp256k1.ParseStringToPrivate(privKeyStr)
	if err != nil {
		return nil, fmt.Errorf("parse private string failed! %s", err.Error())
	}

	if j.Crypto.MnemonicCipherText != "" {
		mnemonicText, err := hex.DecodeString(j.Crypto.MnemonicCipherText)
		if err != nil {
			return nil, fmt.Errorf("decode mnemonicText failed! %s", err.Error())
		}
		mnemonicStr, ok = aes.AESCFBDecrypt(bytes.Join([][]byte{passWd, salt}, []byte{}), mnemonicText)
		if !ok {
			return nil, errors.New("aes decrypt failed")
		}
	}
	privateJson := &Private{Private: privKey.String(), Mnemonic: mnemonicStr}
	return privateJson, nil
}

func ReadJson(jsonFile string) (*Json, error) {
	var j *Json
	if bytes, err := ioutil.ReadFile(jsonFile); err != nil {
		return nil, err
	} else if err := json.Unmarshal(bytes, &j); err != nil {
		return nil, err
	}
	return j, nil
}

func getRandSalt(n int) ([]byte, error) {
	mainBuff := make([]byte, n)
	_, err := io.ReadFull(rand.Reader, mainBuff)
	if err != nil {
		return nil, fmt.Errorf("reading from crypto/rand failed: " + err.Error())
	}
	return mainBuff, nil
}

func saveJson(dir string, j *Json) error {
	bytes, err := json.Marshal(j)
	if err != nil {
		return err
	}
	path := dir + "/" + j.Address + ".json"
	err = ioutil.WriteFile(path, bytes, os.ModePerm)
	if err != nil {
		return err
	}
	return utils.WriteLineAppendFile(dir+boot, j.Address+".json")
}
