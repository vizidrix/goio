package aes

import (
	"bytes"
	"crypto/sha256"
	"fmt"
	"testing"
)

func Test_Should_encrypt_and_decrypt_bytes(t *testing.T) {
	key := []byte("1a2a3a4a5a 1a2a3a4a5a 1a2a3a4a5a")
	data := []byte("some secret information to encode")
	var encrypted []byte
	var decrypted []byte
	var err error

	if encrypted, err = aes.AesEncrypt(key, data); err != nil {
		t.Errorf("Error encrypting\n\t- %s\n", err)
	}
	if decrypted, err = aes.AesDecrypt(key, encrypted); err != nil {
		t.Errorf("Error decrypting\n\t- %s\n", err)
	}
	if !ByteSliceEqual(data, decrypted) {
		t.Errorf("Decrypted data expected [%v] but was [%v]\n", data, decrypted)
	}
}

func Test_Should_write_and_read(t *testing.T) {
	hash := sha256.Sum256([]byte{118, 105, 122, 105, 100, 114, 105, 120})
	key := hash[:]
	secret := []byte("some things shouldn't be said")

	fmt.Printf("Key:\n%v\nSize: %d\n", key, len(key))

	buffer := new(bytes.Buffer)
	w_handle, _ := NewWriter(buffer, key)
	if count, err := w_handle.Write(secret); count == 0 || err != nil {
		fmt.Printf("%d - %s\n", count, err)
	}
	if err := w_handle.Flush(); err != nil {
		fmt.Printf("%s\n", err)
	}
	fmt.Println("Closing...")
	if err := w_handle.Close(); err != nil {
		fmt.Printf("%s\n", err)
	}

	fmt.Printf("Encrypted:\n%s\n", buffer.Bytes())

	bytes_r_handle := bytes.NewReader(buffer.Bytes())

	fmt.Printf("Byte handle: %d\n", bytes_r_handle.Len())
	r_handle, _ := NewReader(bytes_r_handle, key)

	data := make([]byte, len(buffer.Bytes()))
	fmt.Printf("Data size: %d\n", len(data))
	if count, err := r_handle.Read(data); count == 0 || err != nil {
		fmt.Printf("%d - %s\n", count, err)
	}
	fmt.Printf("Data: %v\n\t%s\n", data, data)

	t.Fail()
}
