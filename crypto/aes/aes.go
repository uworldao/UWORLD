package aes

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"io"
)

// Encrypt string to base64 crypto using AES
func AESCFBEncrypt(key []byte, text string) ([]byte, bool) {
	plaintext := []byte(text)

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, false
	}

	cipherText := make([]byte, aes.BlockSize+len(plaintext))
	iv := cipherText[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return nil, false
	}

	stream := cipher.NewCFBEncrypter(block, iv)
	stream.XORKeyStream(cipherText[aes.BlockSize:], plaintext)

	return cipherText, true
}

// Decrypt from base64 to decrypted string
func AESCFBDecrypt(key []byte, cipherText []byte) (string, bool) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", false
	}

	if len(cipherText) < aes.BlockSize {
		return "", false
	}
	iv := cipherText[:aes.BlockSize]
	cipherText = cipherText[aes.BlockSize:]

	stream := cipher.NewCFBDecrypter(block, iv)

	stream.XORKeyStream(cipherText, cipherText)

	return string(cipherText), true
}
