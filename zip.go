package gocert

import (
	//"archive/tar"
	//"compress/gzip"
	"errors"
	"fmt"
	"net/http"
	"os"
)

type tarPack struct {
	filePath   string
	fileLength int64
	fileData   []byte
}

func OpenTarPack(path string) (result *tarPack, err error) {
	var file *os.File
	var fileStats os.FileInfo
	// Load the tar pack file from disk
	if file, err = os.Open(path); err != nil {
		return nil, errors.New(fmt.Sprintf("Unable to open tar pack at '%s'\n\t- %s\n", path, err))
	}
	defer file.Close()
	// Get the size to create the buffer
	if fileStats, err = file.Stat(); err != nil {
		return nil, errors.New(fmt.Sprintf("Unable to retrieve tar pack stats\n\t- %s\n", err))
	}

	result = &tarPack{
		filePath:   path,
		fileLength: fileStats.Size(),
		fileData:   make([]byte, fileStats.Size()),
	}
	return
}

func CreateTarPack(path string) (result *tarPack, err error) {
	var file *os.File
	// Make the tar pack file on disk
	if file, err = os.Create(path); err != nil {
		return nil, errors.New(fmt.Sprintf("Unable to create tar pack at '%s'\n\t- %s\n", path, err))
	}
	defer file.Close()

	result = &tarPack{
		filePath:   path,
		fileLength: 0,
		fileData:   make([]byte, 0),
	}
	return
}

func (tarPack *tarPack) Open(name string) (http.File, error) {
	// TODO: verify existence
	return &file{
		tarPack: tarPack,
		path:    name,
	}, nil
}

type file struct {
	tarPack *tarPack
	path    string
	data    []byte
}

func (f *file) Close() error {
	return nil
}

func (f *file) Stat() (os.FileInfo, error) {
	return nil, nil
}

func (f *file) Readdir(count int) ([]os.FileInfo, error) {
	return make([]os.FileInfo, 0), nil
}

func (f *file) Read(buffer []byte) (int, error) {
	return 0, nil
}

func (f *file) Seek(offset int64, whence int) (int64, error) {
	return 0, nil
}

/*
	if filepath.Separator != '/' && strings.IndexRune(name, filepath.Separator) >= 0 ||
    	strings.Contains(name, "\x00") {
    		return nil, errors.New("http: invalid character in file path")
    	}
    	dir := string(d)
    	if dir == "" {
    37			dir = "."
    38		}
    39		f, err := os.Open(filepath.Join(dir, filepath.FromSlash(path.Clean("/"+name))))
    40		if err != nil {
    41			return nil, err
    42		}
    43		return f, nil
}
*/
