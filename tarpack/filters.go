package tarpack

import (
	"os"
	"path/filepath"
	"strings"
)

type FilePredicate func(os.FileInfo) bool

// Limit filter scope to only files
func OnFiles(filters ...FilePredicate) FilePredicate {
	return func(file os.FileInfo) bool {
		if file.IsDir() {
			return false
		}
		for _, filter := range filters {
			if !filter(file) {
				return false
			}
		}
		return true
	}
}

// Limit filter scope to only dirs
func OnDirs(filters ...FilePredicate) FilePredicate {
	return func(file os.FileInfo) bool {
		if !file.IsDir() {
			return false
		}
		for _, filter := range filters {
			if !filter(file) {
				return false
			}
		}
		return true
	}
}

// Make sure none of the contained filters matches
func Not(filters ...FilePredicate) FilePredicate {
	return func(file os.FileInfo) bool {
		for _, filter := range filters {
			if filter(file) {
				return false
			}
		}
		return true
	}
}

func NameContainsAny(whitelist ...string) FilePredicate {
	return func(file os.FileInfo) bool {
		if file == nil {
			return false
		}
		for _, entry := range whitelist {
			if strings.Contains(file.Name(), entry) {
				return true
			}
		}
		return false
	}
}

func PathMatchAny(whitelist ...string) FilePredicate {
	return func(file os.FileInfo) bool {
		if file == nil {
			return false
		}
		var match bool
		var err error
		for _, entry := range whitelist {
			if match, err = filepath.Match(entry, file.Name()); match && err == nil {
				return true
			}
		}
		return false
	}
}

func Where(files []os.FileInfo, filters ...FilePredicate) []os.FileInfo {
	result := make([]os.FileInfo, 0, len(files))
	for _, file := range files {
		for _, filter := range filters {
			if !filter(file) {
				continue
			}
			result = append(result, file)
		}
	}
	return result
}
