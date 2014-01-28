package gocert

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"errors"
	"fmt"
	"io"
)

type aesKey struct {
	key []byte
}

func AesKey(key []byte) (result *aesKey, err error) {
	if len(key) != 32 {
		err = errors.New(fmt.Sprintf("Key length was %d but should be %d\n", len(key), 32))
		return
	}
	result = &aesKey{
		key: key,
	}
	return
}

//func (key *aesKey) NewWriter()
type Writer struct {
	w         io.Writer
	key       [32]byte
	size      int
	buffer    bytes.Buffer
	encrypted []byte
	closed    bool
	err       error
}

func NewWriter(w io.Writer, key [32]byte) *Writer {
	z := new(Writer)
	z.init(w)
	return z, nil
}

func (z *Writer) init(w io.Writer, key [32]byte) {
	z.w = w
	z.key = key
	z.size = 0
	z.buffer = new(bytes.Buffer)
	closed = false
	err = nil
}

func (z *Writer) Write(p []byte) (int, error) {
	if z.err != nil {
		return 0, z.err
	}
	var n int
	z.size += uint32(len(p))
	n, z.err = z.buffer.Write(p)
	return n, z.err
}

func (z *Writer) Flush() error {
	if z.err != nil {
		return z.err
	}
	if z.closed {
		return nil
	}
	return z.err
}

func (z *Writer) Close() error {
	if z.err != nil {
		return z.err
	}
	if z.closed {
		return nil
	}
	z.closed = true
	z.encrypted, z.err = AesDecrypt(z.key, z.buffer[:])
	if z.err != nil {
		return z.err
	}
	_, z.err = z.w.Write(z.encrypted)
	return z.err
}

/*
type aesReader struct {
	*aesKey
	reader io.Reader
	length int
}

func (key *aesKey) WrapReader(reader io.Reader, length int) *aesReader {
	return &aesReader{
		aesKey: key,
		reader: reader,
		length: length,
	}
}

func (aes *aesReader) Read(p []byte) (n int, err error) {
	data := make([]byte, aes.length)
	var count int
	if count, err = aes.Read(data); err != nil {
		n = count
		return
	}
	if p, err = AesDecrypt(aes.key, data); err != nil {
		return 0, errors.New(fmt.Sprintf("Error reading from Aes: %s\n", err))
	}
	n = len(p)
	return
}
*/

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
