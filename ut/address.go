package ut

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/uworldao/UWORLD/common/hasharry"
	"github.com/uworldao/UWORLD/crypto/base58"
	"github.com/uworldao/UWORLD/crypto/ecc/secp256k1"
	"github.com/uworldao/UWORLD/crypto/hash"
	"github.com/uworldao/UWORLD/param"
	"unicode"
)

const addressLength = 36
const addressBytesLength = 27

// Generate address by secp256k1 public key
func GenerateAddress(version string, key *secp256k1.PublicKey) (string, error) {
	addr := GenerateUWDAddress(version, key)
	return addr, nil
}

// Check the secondary account name, it must be letters,
// all uppercase or all lowercase, no more than 10
// characters and no less than 2.
func CheckAbbr(abbr string) error {
	if len(abbr) < 2 || len(abbr) > 10 {
		return errors.New("the coin abbr length must be in the range of 2 and 10")
	}
	for _, c := range abbr {
		if !unicode.IsLetter(c) {
			return errors.New("coin abbr must be letter")
		}
		if !unicode.IsUpper(c) {
			return errors.New("coin abbr must be upper")
		}
	}
	return nil
}

// Generate UWD address
func GenerateUWDAddress(version string, key *secp256k1.PublicKey) string {
	ver := []byte{}
	switch version {
	case param.MainNet:
		ver = append(ver, param.MainPubKeyHashAddrID[0:]...)
	case param.TestNet:
		ver = append(ver, param.TestPubKeyHashAddrID[0:]...)
	default:
		return ""
	}

	hashed1 := hash.Hash(key.SerializeCompressed())
	hashed2, _ := hash.Hash160(hashed1.Bytes())
	addVersion := append(ver, hashed2...)
	addVersionHashed1 := hash.Hash(addVersion)
	addVersionHashed2 := hash.Hash(addVersionHashed1.Bytes())
	checkSum := addVersionHashed2[0:4]
	hashedCheck1 := append(addVersion, checkSum...)
	return base58.Encode(hashedCheck1)
}

// Verify UWD address
func CheckUWDAddress(version string, addr string) bool {
	ver := []byte{}
	switch version {
	case param.MainNet:
		ver = append(ver, param.MainPubKeyHashAddrID[0:]...)
	case param.TestNet:
		ver = append(ver, param.TestPubKeyHashAddrID[0:]...)
	default:
		return false
	}

	if addr == param.EaterAddress.String() {
		return true
	}

	if len(addr) != addressLength {
		return false
	}
	addrBytes := base58.Decode(addr)
	if len(addrBytes) != addressBytesLength {
		return false
	}
	checkSum := addrBytes[len(addrBytes)-4:]
	checkBytes := addrBytes[0 : len(addrBytes)-4]
	checkBytesHashed1 := hash.Hash(checkBytes)
	checkBytesHashed2 := hash.Hash(checkBytesHashed1.Bytes())
	versionBytes := checkBytes[0:3]
	if bytes.Compare(versionBytes, ver) != 0 {
		return false
	}
	return bytes.Compare(checkSum, checkBytesHashed2[0:4]) == 0
}

// Generate contract address
func GenerateContractAddress(net string, address string, abbr string) (string, error) {
	ver := []byte{}
	switch net {
	case param.MainNet:
		ver = append(ver, param.MainPubKeyHashTokenID[0:]...)
	case param.TestNet:
		ver = append(ver, param.TestPubKeyHashTokenID[0:]...)
	default:
		return "", errors.New("wrong network")
	}
	if err := CheckAbbr(abbr); err != nil {
		return "", err
	}
	if !CheckUWDAddress(net, address) {
		return "", errors.New("incorrect address")
	}
	addrBytes := base58.Decode(address)
	buffBytes := append(addrBytes, []byte(abbr)...)
	fmt.Println(buffBytes)
	hashed := hash.Hash(buffBytes)
	hash160, err := hash.Hash160(hashed.Bytes())
	if err != nil {
		return "", err
	}

	addNet := append(ver, hash160...)
	hashed1 := hash.Hash(addNet)
	hashed2 := hash.Hash(hashed1.Bytes())
	checkSum := hashed2[0:4]
	hashedCheck1 := append(addNet, checkSum...)
	code58 := base58.Encode(hashedCheck1)
	return hasharry.StringToAddress(code58).String(), nil
}

// Verify contract address
func CheckContractAddress(net string, address string, abbr string, contractAddress string) bool {
	if !IsValidContractAddress(net, contractAddress) {
		return false
	}
	newAddress, err := GenerateContractAddress(net, address, abbr)
	if err != nil {
		return false
	}
	return newAddress == contractAddress
}

func IsValidContractAddress(net string, address string) bool {
	if address == param.Token.String() {
		return true
	}
	ver := []byte{}
	switch net {
	case param.MainNet:
		ver = append(ver, param.MainPubKeyHashTokenID[0:]...)
	case param.TestNet:
		ver = append(ver, param.TestPubKeyHashTokenID[0:]...)
	default:
		return false
	}
	if len(address) != addressLength {
		return false
	}
	addrBytes := base58.Decode(address)
	if len(addrBytes) != addressBytesLength {
		return false
	}
	checkSum := addrBytes[len(addrBytes)-4:]
	checkBytes := addrBytes[0 : len(addrBytes)-4]
	checkBytesHashed1 := hash.Hash(checkBytes)
	checkBytesHashed2 := hash.Hash(checkBytesHashed1.Bytes())
	netBytes := checkBytes[0:3]
	if bytes.Compare(ver, netBytes) != 0 {
		return false
	}
	return bytes.Compare(checkSum, checkBytesHashed2[0:4]) == 0
}
