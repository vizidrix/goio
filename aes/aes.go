package aes

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"errors"
	"fmt"
	"io"
)

func AesEncrypt(key []byte, data []byte) (result []byte, err error) {
	if len(key) != 16 && len(key) != 24 && len(key) != 32 {
		err = errors.New(fmt.Sprintf("Key length was %d but should be 16, 24 or 32\n", len(key)))
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
	if len(key) != 16 && len(key) != 24 && len(key) != 32 {
		err = errors.New(fmt.Sprintf("Key length was %d but should be 16, 24 or 32\n", len(key)))
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

type Reader struct {
	r         io.Reader
	key       []byte
	size      int
	buffer    *bytes.Buffer
	decrypted []byte
	closed    bool
	err       error
}

func NewReader(r io.Reader, key []byte) (*Reader, error) {
	z := new(Reader)
	z.init(r, key)
	if len(key) != 16 && len(key) != 24 && len(key) != 32 {
		z.err = errors.New("AES Key must be 16, 24 or 32 bytes")
		return z, z.err
	}
	return z, nil
}

func (z *Reader) init(r io.Reader, key []byte) {
	z.r = r
	z.key = key
	z.size = 0
	z.buffer = new(bytes.Buffer)
	z.closed = false
	z.err = nil
}

func (z *Reader) Read(p []byte) (n int, err error) {
	if z.err != nil {
		return 0, z.err
	}
	temp := make([]byte, len(p))
	n, z.err = z.r.Read(temp)
	z.size += int(n)
	if n == 0 || z.err != nil {
		if z.err != io.EOF { // EOF, decrypt and send
			return z.size, z.err
		}
	}
	temp, z.err = AesDecrypt(z.key, temp)
	if z.err != nil {
		return z.size, z.err
	}
	copy(p, temp)
	return z.size, z.err
}

type Writer struct {
	w         io.Writer
	key       []byte
	size      int
	buffer    *bytes.Buffer
	encrypted []byte
	closed    bool
	err       error
}

func NewWriter(w io.Writer, key []byte) (*Writer, error) {
	z := new(Writer)
	z.init(w, key)
	if len(key) != 16 && len(key) != 24 && len(key) != 32 {
		z.err = errors.New("AES Key must be 16, 24 or 32 bytes")
		return z, z.err
	}
	return z, nil
}

func (z *Writer) init(w io.Writer, key []byte) {
	z.w = w
	z.key = key
	z.size = 0
	z.buffer = new(bytes.Buffer)
	z.closed = false
	z.err = nil
}

func (z *Writer) Write(p []byte) (int, error) {
	if z.err != nil {
		return 0, z.err
	}
	var n int
	z.size += int(len(p))
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
	var encrypted []byte = make([]byte, z.size)
	fmt.Println("Decrypting")
	encrypted, z.err = AesEncrypt(z.key, z.buffer.Bytes())
	if z.err != nil {
		return z.err
	}
	fmt.Println("Decrypted")
	z.encrypted = encrypted
	_, z.err = z.w.Write(z.encrypted)
	return z.err
}
