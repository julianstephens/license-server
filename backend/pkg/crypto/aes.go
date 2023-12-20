package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"fmt"
	"io"
)

// GetBlock creates a new cipher.Block using an AES key of size 16, 24, or 32 bytes
func GetBlock(key []byte) (block cipher.Block, err error) {
	block, err = aes.NewCipher(key)
	if err != nil {
		return
	}

	return
}

// Encrypt uses the provided cipher.Block to encrypt text
func Encrypt(block cipher.Block, text string) (ciphertext []byte, err error) {
	byteMsg := []byte(text)

	ciphertext = make([]byte, aes.BlockSize+len(byteMsg))
	iv := ciphertext[:aes.BlockSize]
	if _, err = io.ReadFull(rand.Reader, iv); err != nil {
		return
	}

	stream := cipher.NewCFBEncrypter(block, iv)
	stream.XORKeyStream(ciphertext[aes.BlockSize:], byteMsg)

	return
}

// Decrypt uses the provided cipher.Block to decrypt text
func Decrypt(block cipher.Block, ciphertext []byte) (decoded []byte, err error) {
	if len(ciphertext) < aes.BlockSize {
		err = fmt.Errorf("invalid ciphertext block size")
		return
	}

	iv := ciphertext[:aes.BlockSize]
	ciphertext = ciphertext[aes.BlockSize:]

	stream := cipher.NewCFBDecrypter(block, iv)
	stream.XORKeyStream(ciphertext, ciphertext)

	decoded = ciphertext

	return
}
