package gocert

import (
	"testing"
)

func Test_Should_encrypt_and_decrypt_bytes(t *testing.T) {
	key := []byte("1a2a3a4a5a 1a2a3a4a5a 1a2a3a4a5a")
	data := []byte("some secret information to encode")
	var encrypted []byte
	var decrypted []byte
	var err error

	//fmt.Printf("Encrypting... \n\t%v\n", data)
	if encrypted, err = AesEncrypt(key, data); err != nil {
		t.Errorf("Error encrypting\n\t- %s\n", err)
	}
	//fmt.Printf("Decrypting... \n\t%v\n", encrypted)
	if decrypted, err = AesDecrypt(key, encrypted); err != nil {
		t.Errorf("Error decrypting\n\t- %s\n", err)
	}
	//fmt.Printf("Comparing...\n")
	if !ByteSliceEqual(data, decrypted) {
		t.Errorf("Decrypted data expected [%v] but was [%v]\n", data, decrypted)
	}
}

func ByteSliceEqual(l, r []byte) bool {
	if len(l) != len(r) {
		return false
	}
	for i, v := range l {
		if v != r[i] {
			return false
		}
	}
	return true
}
