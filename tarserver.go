package goio

import (
	"archive/tar"
	"bytes"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"time"
)

type tarServer struct {
	cache map[string]*tarFile
}

func TarServer(path string) (*tarServer, error) {
	var buffer []byte
	var reader *bytes.Reader
	var handle *tar.Reader
	var cache map[string]*tarFile
	var header *tar.Header
	var err error
	fmt.Printf("Opening tar server: %s\n", path)

	if buffer, err = ioutil.ReadFile(path); err != nil {
		return nil, err
	}
	reader = bytes.NewReader(buffer)
	if handle = tar.NewReader(reader); err != nil {
		return nil, err
	}
	modTime := time.Now().Add(-2 * time.Second)
	cache = make(map[string]*tarFile)
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
				//fmt.Printf("[%d]\t\\-- %s\n", len(cache), part)
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
				//fmt.Printf("\n[%d]__%s\n", len(cache), partPath)
			}
		}
	}
	return &tarServer{
		cache: cache,
	}, nil
}

func (server *tarServer) Open(name string) (http.File, error) {
	var content *tarFile
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

type tarFile struct {
	name     string
	isDir    bool
	modTime  time.Time
	offset   int64
	data     []byte
	children []*tarFile
}

func TarFile(name string, isDir bool, modTime time.Time, data []byte) *tarFile {
	return &tarFile{
		name:     name,
		isDir:    isDir,
		modTime:  modTime,
		offset:   0,
		data:     data,
		children: make([]*tarFile, 0),
	}
}

func (file *tarFile) String() string {
	fileType := "F"
	if file.isDir {
		fileType = "D"
	}
	return fmt.Sprintf("[%s] %s\t[%d]\tC[%d]", fileType, file.name, len(file.data), len(file.children))
}

func (file *tarFile) AddChild(child *tarFile) {
	//fmt.Printf("Adding:\n\t%s\nTo:\n\t%s\n\n", child, file)
	file.children = append(file.children, child)
}

func (file *tarFile) Close() error {
	file.offset = 0
	//fmt.Printf("CLOSING\t[%s]\n", file.name)
	return nil // N/A
}

func (file *tarFile) Stat() (os.FileInfo, error) {
	//fmt.Printf("STAT\t[%s]\n", file)
	return &tarFileInfo{
		file: file,
	}, nil
}

func (file *tarFile) Readdir(count int) ([]os.FileInfo, error) {
	//fmt.Printf("READDIR\t[%s]\n", file)

	result := make([]os.FileInfo, len(file.children))
	var err error
	for i, child := range file.children {
		//fmt.Printf("File[%d]: %s\n", i, file)
		if result[i], err = child.Stat(); err != nil {
			return nil, err
		}
	}
	return result, nil
}

func (file *tarFile) Read(data []byte) (int, error) {
	buffer := bytes.NewReader(file.data)
	buffer.Seek(file.offset, 0)
	count, err := buffer.Read(data)
	file.offset += int64(count)
	//fmt.Printf("READ\t[%d] %s\n", count, file.name)
	return count, err
}

func (file *tarFile) Seek(offset int64, whence int) (int64, error) {
	//fmt.Printf("SEEK\t[%s]\n", file.name)
	if offset >= int64(len(file.data)) || offset < 0 {
		return 0, errors.New("Invalid offset")
	}
	file.offset = offset
	return file.offset, nil
}

type tarFileInfo struct {
	file *tarFile
}

func (fileInfo *tarFileInfo) Name() string {
	//fmt.Printf("INFO NAME\t\t[%s]\n", fileInfo.file.name)
	return fileInfo.file.name
}

func (fileInfo *tarFileInfo) Size() int64 {
	size := int64(len(fileInfo.file.data))
	//fmt.Printf("INFO SIZE\t\t[%s] [%d]\n", fileInfo.file.name, size)
	return size
}

func (fileInfo *tarFileInfo) Mode() os.FileMode {
	//fmt.Printf("INFO MODE\t\t[%s]\n", fileInfo.file.name)
	return os.FileMode(os.O_RDONLY)
}

func (fileInfo *tarFileInfo) ModTime() time.Time {
	//fmt.Printf("INFO TIME\t\t[%s] [%s]\n", fileInfo.file.name, fileInfo.file.modTime)
	return fileInfo.file.modTime
}

func (fileInfo *tarFileInfo) IsDir() bool {
	//fmt.Printf("INFO DIR\t\t[%s] [%v]\n", fileInfo.file.name, fileInfo.file.isDir)
	return fileInfo.file.isDir
}

func (fileInfo *tarFileInfo) Sys() interface{} {
	//fmt.Printf("TarInfo sys\n")
	return nil
}
