package tarpack

import (
	"fmt"
	"os"
	"testing"
	"time"
)

func Test_Should_return_false_on_NameContainsAny_for_nil_files(t *testing.T) {
	filter := NameContainsAny("")
	if filter(nil) {
		fmt.Printf("NameContainsAny should have returned false for nil file")
		t.Fail()
	}
}

var nameContainsAnyTests = []struct {
	name   string
	filter string
	result bool
}{
	{"test.tar", "*.tar", false},
	{"test.toml", "*.toml", false},
	{"test.toml", ".toml", true},
	{"test.tar", "*.toml", false},
	{"something.exe", "thing", true},
}

func Test_Should_filter_out_invalid_files_on_NameContainsAny(t *testing.T) {
	for i, test := range nameContainsAnyTests {
		fileInfo := &fakeFileInfo{
			name: test.name,
		}
		filter := NameContainsAny(test.filter)
		result := filter(fileInfo)
		if result != test.result {
			t.Errorf("\n%d. NameContainsAny(%q)(%q) =>\n\tGot [%q] but want [%q]\n", i, test.filter, test.name, result, test.result)
		}
	}
}

func Test_Should_return_false_on_PathMatchAny_for_nil_files(t *testing.T) {
	filter := PathMatchAny("")
	if filter(nil) {
		fmt.Printf("PathMatchAny should have returned false for nil file")
		t.Fail()
	}
}

var pathMatchAnyTests = []struct {
	name   string
	filter string
	result bool
}{
	{"test.tar", "*.tar", true},
	{"test.toml", "*.toml", true},
	{"test.tar", ".tar", false},
	{"test.tar", "tar", false},
	{"test.toml", ".toml", false},
	{"test.tar", "*.toml", false},
	{"a", "a", true},
}

func Test_Should_filter_out_invalid_files_on_PathMatchAny(t *testing.T) {
	for i, test := range pathMatchAnyTests {
		fileInfo := &fakeFileInfo{
			name: test.name,
		}
		filter := PathMatchAny(test.filter)
		result := filter(fileInfo)
		if result != test.result {
			t.Errorf("\n%d. PathMatchAny(%q)(%q) =>\n\tGot [%q] but want [%q]\n", i, test.filter, test.name, result, test.result)
		}
	}
}

type fakeFileInfo struct {
	name    string
	size    int64
	mode    os.FileMode
	modTime time.Time
	isDir   bool
	sys     interface{}
}

func (fi *fakeFileInfo) Name() string {
	return fi.name
}

func (fi *fakeFileInfo) Size() int64 {
	return fi.size
}

func (fi *fakeFileInfo) Mode() os.FileMode {
	return fi.mode
}

func (fi *fakeFileInfo) ModTime() time.Time {
	return fi.modTime
}

func (fi *fakeFileInfo) IsDir() bool {
	return fi.isDir
}

func (fi *fakeFileInfo) Sys() interface{} {
	return fi.sys
}
