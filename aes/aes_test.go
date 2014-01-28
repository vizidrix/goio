package aes

import (
	"bytes"
	"crypto/sha256"
	"fmt"
	"testing"
)

func Test_Should_write_and_read(t *testing.T) {
	hash := sha256.Sum256([]byte{118, 105, 122, 105, 100, 114, 105, 120})
	key := hash[:]
	secret := []byte("some things shouldn't be said")
	//key := make([]byte, 256)

	//copy(key[:32], hash[:])
	fmt.Printf("Key:\n%v\nSize: %d\n", key, len(key))

	/*for i := byte(0); i <= 255; i++ {
		key[i] = i
	}*/
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

	//r_buffer := new(bytes.Buffer)
	//copy(r_buffer, buffer.Bytes())
	bytes_r_handle := bytes.NewReader(buffer.Bytes())

	fmt.Printf("Byte handle: %d\n", bytes_r_handle.Len())
	//r_buffer, _ := bytes.NewReader(buffer.Bytes())
	r_handle, _ := NewReader(bytes_r_handle, key)

	data := make([]byte, len(buffer.Bytes()))
	fmt.Printf("Data size: %d\n", len(data))
	if count, err := r_handle.Read(data); count == 0 || err != nil {
		fmt.Printf("%d - %s\n", count, err)
	}
	fmt.Printf("Data: %v\n\t%s\n", data, data)

	t.Fail()
}
