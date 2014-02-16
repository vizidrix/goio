package goio

import (
	"archive/tar"
	"bytes"
	"errors"
	"fmt"
	. "github.com/vizidrix/goio/tarfile"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

type tarServer struct {
	cache map[string]ReadWriteContainer
}

func TarServer(path string) (*tarServer, error) {
	var buffer []byte
	var err error
	fmt.Printf("Opening tar server: [ %s ]\n", path)
	if buffer, err = ioutil.ReadFile(path); err != nil {
		return nil, err
	}
	return RawTarServer(buffer)
}

func RawTarServer(data []byte) (*tarServer, error) {
	var reader *bytes.Reader
	var handle *tar.Reader
	var cache map[string]ReadWriteContainer
	var header *tar.Header
	var err error

	reader = bytes.NewReader(data)
	if handle = tar.NewReader(reader); err != nil {
		return nil, err
	}
	modTime := time.Now().Add(-2 * time.Second)
	cache = make(map[string]ReadWriteContainer)
	cache["/"] = TarFile("/", true, modTime, make([]byte, 0))
	for {
		if header, err = handle.Next(); err != nil {
			if err == io.EOF {
				break // End of archive
			}
			return nil, err
		}
		parentFile := cache["/"]
		parts := strings.Split(header.Name, "/")
		partPath := "/"

		for i, part := range parts {
			if i == len(parts)-1 { // File
				b := new(bytes.Buffer)
				if _, err = io.Copy(b, handle); err != nil {
					return nil, err
				}
				partPath += part
				file := TarFile(partPath, false, modTime, b.Bytes())
				parentFile.AddChild(file)
				cache[partPath] = file
				break
			} // Dir

			partPath += part + "/"
			if tempFile, ok := cache[partPath]; ok {
				parentFile = tempFile
				continue
			} else { // Didn't find the dir in the cache
				// Make the dir, add it, cache it and set it to parent
				dir := TarFile(partPath, true, modTime, make([]byte, 0))
				parentFile.AddChild(dir)
				cache[partPath] = dir
				parentFile = dir
			}
		}
	}
	return &tarServer{
		cache: cache,
	}, nil
}

func (server *tarServer) Open(name string) (http.File, error) {
	var content ReadWriteContainer
	var ok bool

	if len(name) > 2 && name[0:2] == "//" {
		name = name[1:]
	}
	fmt.Printf("|")
	//fmt.Printf("\nREQ\t[%s]\n", name)
	if content, ok = server.cache[name]; !ok {
		modTime := time.Now().Add(-2 * time.Second)
		return TarFile(name, false, modTime, make([]byte, 0)), errors.New(fmt.Sprintf("File [%s] not found in cache\n", name))
	}
	return content, nil
}
