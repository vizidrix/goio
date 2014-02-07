package tarpack

import (
	"archive/tar"
	//"compress/gzip"
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	//"net/http"
	"os"
)

/*
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
*/

/*
rootPath	- Path to read the contents from
relPath		- Initial internal path to begin writing into
filters		- List of filters to apply to the folder/file set
*/
func TarDir(rootPath, relPath string, filters ...FilePredicate) ([]byte, error) {
	var err error
	buffer := new(bytes.Buffer)
	handle := tar.NewWriter(buffer)
	if err = tarDir(handle, rootPath, relPath, filters...); err != nil {
		return nil, errors.New(fmt.Sprintf("\nTarDir: %s\n", err))
	}
	if err = handle.Close(); err != nil {
		fmt.Printf("Error closing TAR file")
		return nil, err
	}
	return buffer.Bytes(), nil
}

/*
handle		- Pointer to a Tar file being built
rootPath	- Path to read the contents from
relPath		- Path to write these files into
filters		- List of filters to apply to the folder/file set
*/
func tarDir(handle *tar.Writer, rootPath, relPath string, filters ...FilePredicate) error {
	//fmt.Printf("\nTarring dir:\n\t[%s]\n\t=>\n\t[%s]\n\n", rootPath, relPath)
	var err error
	var files []os.FileInfo
	if files, err = ioutil.ReadDir(rootPath); err != nil {
		return errors.New(fmt.Sprintf("Root path [%s] error: %s", rootPath, err))
	}
	for _, file := range Where(files, filters...) {
		newPath := relPath
		if len(relPath) != 0 && relPath[:len(relPath)] != "/" {
			newPath += "/"
		}
		/*
			for _, filter := range filters {
				if !filter(file) {
					continue
				}
			}
		*/
		if file.IsDir() {
			tarDir(handle, rootPath+"/"+file.Name(), newPath+file.Name(), filters...)
		} else { // File is a file
			//fmt.Printf("Writing file: %s -> %s\n", relPath, file.Name())
			if err = writeFile(handle, rootPath+"/"+file.Name(), newPath+file.Name()); err != nil {
				fmt.Printf("Error writing file: %s\n", err)
				return err
			}
		}
	}
	return nil
}

/*
	if file.IsDir() { // File is a dir and not a file
		for _, filter := range filters {

			if !filter(file) {
				continue
			}
		}
		if !folderFilter(file.Name()) {
			continue
		}
		tarDir(handle, rootPath+"/"+file.Name(), newPath+file.Name(), folderFilter, fileFilter)
	} else { // File is a file and not a dir
		if _, filter := range filters {
			if !filter(file) {
				continue
			}
		}

		if !fileFilter(file.Name()) {
			continue
		}

		fmt.Printf("Writing file: %s -> %s\n", relPath, file.Name())
		if err = writeFile(handle, rootPath+"/"+file.Name(), newPath+file.Name()); err != nil {
			fmt.Printf("Error writing file: %s\n", err)
			return err
		}
	}
*/
//}
//return nil
//}

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
