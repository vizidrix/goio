package tarfile

import (
	"bytes"
	"errors"
	"fmt"
	//"io"
	"net/http"
	"os"
	"time"
)

type ReadWriteContainer interface {
	//io.ReadCloser
	http.File
	AddChild(ReadWriteContainer)
}

type tarFile struct {
	name    string
	isDir   bool
	modTime time.Time
	offset  int64
	data    []byte
	//children []*tarFile
	children []ReadWriteContainer
}

func TarFile(name string, isDir bool, modTime time.Time, data []byte) *tarFile {
	return &tarFile{
		name:    name,
		isDir:   isDir,
		modTime: modTime,
		offset:  0,
		data:    data,
		//children: make([]*tarFile, 0),
		children: make([]ReadWriteContainer, 0),
	}
}

func (file *tarFile) String() string {
	fileType := "F"
	if file.isDir {
		fileType = "D"
	}
	return fmt.Sprintf("[%s] %s\t[%d]\tC[%d]", fileType, file.name, len(file.data), len(file.children))
}

//func (file *tarFile) AddChild(child *tarFile) {
func (file *tarFile) AddChild(child ReadWriteContainer) {
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
