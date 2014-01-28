package goio

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"crypto/sha256"
	"fmt"
	"github.com/vizidrix/goio/aes"
	"io"
	"os"
	"testing"
)

func temp() { fmt.Println("asdf") }

func Test_Should_pack_contents_and_then_read_back(t *testing.T) {
	makeTheThings()
	//readTheThings()
	t.Fail()
}

func makeTheThings() {
	//file_w_handle, _ := os.Create("testing/testfile.tar.gz")
	//defer file_w_handle.Close()
	hash := sha256.Sum256([]byte{118, 105, 122, 105, 100, 114, 105, 120})
	key := hash[:]

	buffer := new(bytes.Buffer)
	aes_w_handle, _ := aes.NewWriter(buffer, key)
	//zip_w_handle := gzip.NewWriter(file_w_handle)
	zip_w_handle, _ := gzip.NewWriterLevel(aes_w_handle, gzip.BestCompression) // NewWriterLevel
	//defer zip_w_handle.Close()
	tar_w_handle := tar.NewWriter(zip_w_handle)
	//defer tar_w_handle.Close()
	var files = []struct{ Name, Body string }{
		{"root.txt", "root content"},
		{"/sub1/sub1.txt", "sub1 content"},
	}
	for _, file := range files {
		hdr := &tar.Header{
			Name: file.Name,
			Size: int64(len(file.Body)),
		}
		if err := tar_w_handle.WriteHeader(hdr); err != nil {
			panic("Write header")
			//t.Fail()
		}
		if _, err := tar_w_handle.Write([]byte(file.Body)); err != nil {
			panic("Write body")
			//t.Fail()
		}
	}
	// transition to read
	err := tar_w_handle.Close()
	if err != nil {
		panic("bad things happened")
	}
	zip_w_handle.Close()
	aes_w_handle.Close()

	// begin reading
	bytes_r_handle := bytes.NewReader(buffer.Bytes())
	aes_r_handle, _ := aes.NewReader(bytes_r_handle, key)
	gzip_r_handle, _ := gzip.NewReader(aes_r_handle)
	tar_r_handle := tar.NewReader(gzip_r_handle)

	for {
		header, err := tar_r_handle.Next()
		if err == io.EOF {
			break // End of archive
		}
		if err != nil {
			panic("bad reader")
		}
		fmt.Printf("\nContents of %s:\n", header.Name)
		if _, err := io.Copy(os.Stdout, tar_r_handle); err != nil {
			panic(err)
		}
		fmt.Println()
		fmt.Println()
	}

	gzip_r_handle.Close()

}

func readTheThings() {
	file_r_handle, _ := os.OpenFile("testing/testfile.tar.gz", os.O_RDWR, os.ModePerm)
	defer file_r_handle.Close()
	var zip_r_handle *gzip.Reader
	var err error
	if zip_r_handle, err = gzip.NewReader(file_r_handle); err != nil {
		return
	}
	defer zip_r_handle.Close()
	tar_r_handle := tar.NewReader(zip_r_handle)
	tar_r_handle.Next()
	/*
		for header, err := tar_r_handle.Next(); err == nil; tar_r_handle.Next() {
			fmt.Printf("File handle: %s\n", header.Name)
		}
	*/
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
