package main

import (
	"archive/tar"
	"bytes"
	//"compress/gzip"
	//"crypto/sha256"
	"fmt"
	//"github.com/vizidrix/goio/aes"
	//"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	var err error
	targetDirPath := ""

	fileWhiteList := []string{
		"*.js",
		"*.html",
		"*.png",
		"*.css",
		"*.ttf",
	}
	//decodeString := "intelinde"

	//hash := sha256.Sum256([]byte(decodeString))
	//key := hash[:]

	buffer := new(bytes.Buffer)
	//aes_w_handle, _ := aes.NewWriter(buffer, key)
	//zip_w_handle, _ := gzip.NewWriterLevel(aes_w_handle, gzip.BestCompression)

	//zip_w_handle, _ := gzip.NewWriterLevel(buffer, gzip.BestCompression)
	//tar_w_handle := tar.NewWriter(zip_w_handle)

	tar_w_handle := tar.NewWriter(buffer)
	writeDir(tar_w_handle, targetDirPath, "", fileWhiteList)

	/*
		files, _ := ioutil.ReadDir(targetDirPath)
		for _, file := range files {
			if !file.IsDir() {
				fmt.Printf("Writing file: %s\n", file.Name())
				if err = writeFile(tar_w_handle, targetDirPath+"/"+file.Name()); err != nil {
					fmt.Printf("Error writing file: %s\n", err)
					return
				}
			}
		}
	*/

	if err = tar_w_handle.Close(); err != nil {
		fmt.Printf("Error closing tar:\n\t- %s\n", err)
	}
	/*
		if err = zip_w_handle.Close(); err != nil {
			fmt.Printf("Error closing zip:\n\t- %s\n", err)
		}
	*/
	//aes_w_handle.Close()

	ioutil.WriteFile("indeweb.tar", buffer.Bytes(), 0666)
}

func isWhitelistedFile(whitelist []string, name string) bool {
	for _, entry := range whitelist {
		if match, _ := filepath.Match(entry, name); match {
			return true
		}
	}
	return false
}

func writeDir(handle *tar.Writer, rootPath, relPath string, fileWhiteList []string) (err error) {
	fmt.Printf("Writing dir: %s => %s\n", rootPath, relPath)
	var files []os.FileInfo
	if files, err = ioutil.ReadDir(rootPath); err != nil {
		return
	}
	for _, file := range files {
		newPath := relPath
		if len(relPath) != 0 && relPath[:len(relPath)] != "/" {
			newPath += "/"
		}
		if file.IsDir() { // File is a dir and not a file
			if strings.Contains(file.Name(), ".git") {
				continue // TODO: Improve this filter
			}
			writeDir(handle, rootPath+"/"+file.Name(), newPath+file.Name(), fileWhiteList)
		} else { // File is a file and not a dir
			if !isWhitelistedFile(fileWhiteList, file.Name()) {
				continue
			}
			/*
				if match, _ := filepath.Match(whitelist, file.Name()); !match {
					// Skip files that aren't in whitelist pattern
					continue
				}
			*/
			fmt.Printf("Writing file: %s -> %s\n", relPath, file.Name())
			if err = writeFile(handle, rootPath+"/"+file.Name(), newPath+file.Name()); err != nil {
				fmt.Printf("Error writing file: %s\n", err)
				return
			}
		}
	}
	return
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
	/*newPath := relPath
	if len(relPath) != 0 && relPath[:len(relPath)] != "/" {
		newPath += "/"
	}
	newPath += stat.Name()
	fmt.Printf("Writing to: %s\n", newPath)
	*/
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
