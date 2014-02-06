package goio

import (
	"archive/tar"
	//"compress/gzip"
	//"errors"
	"bytes"
	"fmt"
	"io/ioutil"
	//"net/http"
	"os"
	"path/filepath"
	"strings"
)

func ContainsPredicate(whitelist []string, name string) bool {
	for _, entry := range whitelist {
		if strings.Contains(name, entry) {
			return true
		}
	}
	return false
}

func FileMatchPredicate(whitelist []string, name string) bool {
	for _, entry := range whitelist {
		if match, _ := filepath.Match(entry, name); match {
			return true
		}
	}
	return false
}

func DirFilter(whitelist, blacklist []string) func(string) bool {
	wLen := len(whitelist)
	bLen := len(blacklist)
	if wLen == 0 && bLen == 0 {
		return func(fileName string) bool { return true }
	}
	if wLen == 0 {
		return func(fileName string) bool { return !ContainsPredicate(blacklist, fileName) }
	}
	if bLen == 0 {
		return func(fileName string) bool { return ContainsPredicate(whitelist, fileName) }
	}
	return func(fileName string) bool {
		return ContainsPredicate(whitelist, fileName) && !ContainsPredicate(blacklist, fileName)
	}
}

func FileFilter(whitelist, blacklist []string) func(string) bool {
	wLen := len(whitelist)
	bLen := len(blacklist)
	if wLen == 0 && bLen == 0 {
		return func(fileName string) bool { return true }
	}
	if wLen == 0 {
		return func(fileName string) bool { return !FileMatchPredicate(blacklist, fileName) }
	}
	if bLen == 0 {
		return func(fileName string) bool { return FileMatchPredicate(whitelist, fileName) }
	}
	return func(fileName string) bool {
		return FileMatchPredicate(whitelist, fileName) && !FileMatchPredicate(blacklist, fileName)
	}
}

func TarDir(rootPath, relPath string, folderFilter func(string) bool, fileFilter func(string) bool) error {
	buffer := new(bytes.Buffer)
	handle := tar.NewWriter(buffer)
	return tarDir(handle, rootPath, relPath, folderFilter, fileFilter)
}

func tarDir(handle *tar.Writer, rootPath, relPath string, folderFilter func(string) bool, fileFilter func(string) bool) error {
	fmt.Printf("Tarring dir: %s => %s\n", rootPath, relPath)
	var err error
	var files []os.FileInfo
	if files, err = ioutil.ReadDir(rootPath); err != nil {
		return err
	}
	for _, file := range files {
		newPath := relPath
		if len(relPath) != 0 && relPath[:len(relPath)] != "/" {
			newPath += "/"
		}
		if file.IsDir() { // File is a dir and not a file
			if !folderFilter(file.Name()) {
				continue
			}
			tarDir(handle, rootPath+"/"+file.Name(), newPath+file.Name(), folderFilter, fileFilter)
		} else { // File is a file and not a dir
			if !fileFilter(file.Name()) {
				continue
			}
			fmt.Printf("Writing file: %s -> %s\n", relPath, file.Name())
			if err = writeFile(handle, rootPath+"/"+file.Name(), newPath+file.Name()); err != nil {
				fmt.Printf("Error writing file: %s\n", err)
				return err
			}
		}
	}
	return nil
}

func writeFile(handle *tar.Writer, rootPath, relPath string) error {
	var file *os.File
	var stat os.FileInfo
	var buffer []byte
	var err error
	if file, err = os.OpenFile(rootPath, os.O_RDONLY, os.ModePerm); err != nil {
		return err
	}
	if stat, err = file.Stat(); err != nil {
		return err
	}
	file.Close()
	header := &tar.Header{
		Name: relPath,
		Size: stat.Size(),
	}
	if err := handle.WriteHeader(header); err != nil {
		return err
	}
	if buffer, err = ioutil.ReadFile(rootPath); err != nil {
		return err
	}
	if _, err := handle.Write(buffer); err != nil {
		return err
	}
	return nil
}

/*

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
*/
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
