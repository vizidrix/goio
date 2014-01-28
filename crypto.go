package gocert

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"errors"
	"fmt"
	"io"
)

func AesEncrypt(key []byte, data []byte) (result []byte, err error) {
	if len(key) != 32 {
		err = errors.New(fmt.Sprintf("Key length was %d but should be %d\n", len(key), 32))
		return
	}
	var block cipher.Block
	if block, err = aes.NewCipher(key); err != nil {
		err = errors.New(fmt.Sprintf("Unable to create encrypt cipher from key\n\t- %s\n", err))
		return
	}
	ciphertext := make([]byte, aes.BlockSize+len(data))
	iv := ciphertext[:aes.BlockSize]
	if _, err = io.ReadFull(rand.Reader, iv); err != nil {
		err = errors.New(fmt.Sprintf("Unable to produce IV\n\t- %s\n", err))
	}
	stream := cipher.NewCFBEncrypter(block, iv)
	stream.XORKeyStream(ciphertext[aes.BlockSize:], data)
	result = ciphertext
	return
}

func AesDecrypt(key []byte, data []byte) (result []byte, err error) {
	if len(key) != 32 {
		err = errors.New(fmt.Sprintf("Key length was %d but should be %d\n", 32, len(key)))
		return
	}
	var block cipher.Block
	if block, err = aes.NewCipher(key); err != nil {
		err = errors.New(fmt.Sprintf("Unable to create decrypt cipher from key\n\t- %s\n", err))
	}
	if len(data) < aes.BlockSize {
		err = errors.New(fmt.Sprintf("Incoming data too short to be valid %d\n", len(data)))
		return
	}
	iv := data[:aes.BlockSize]
	data = data[aes.BlockSize:]
	stream := cipher.NewCFBDecrypter(block, iv)
	stream.XORKeyStream(data, data)
	result = data
	return
}
