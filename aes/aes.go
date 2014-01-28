package aes

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"errors"
	"fmt"
	. "github.com/vizidrix/crypto"
	"io"
)

type Reader struct {
	r         io.Reader
	key       [32]byte
	size      int
	buffer    bytes.Buffer
	decrypted []byte
	closed    bool
	err       error
}

func NewReader(r io.Reader, key [32]byte) (*Reader, error) {
	z := new(Reader)
	z.init(r)
	return z, nil
}

func (z *Reader) init(r io.Reader, key [32]byte) {
	z.r = r
	z.key = key
	z.size = 0
	z.buffer = new(bytes.Buffer)
	closed = false
	err = nil
}

func (z *Reader) Read(p []byte) (n int, err error) {
	if z.err != nil {
		return 0, z.err
	}
	for {
		n, z.err = z.r.Read(z.buffer)
		//n, err = z.buffer.Read(p)
		z.size += uint32(n)
		if n == 0 || z.err != nil {
			if z.err == io.EOF { // EOF, decrypt and send
				break
			} else { // Should not read zero bytes until EOF
				return z.size, z.err
			}
		}
	}
	p, z.err = AesDecrypt(z.key, z.buffer[:])
	if z.err != nil {
		return z.size, z.err
	}
	return z.size, z.err
}

type Writer struct {
	w         io.Writer
	key       [32]byte
	size      int
	buffer    bytes.Buffer
	encrypted []byte
	closed    bool
	err       error
}

func NewWriter(w io.Writer, key [32]byte) (*Writer, error) {
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
