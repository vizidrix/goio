package goio

import (
	"testing"
)

func Test_Should_extract_zipped_contents(t *testing.T) {
	var tarPack *tarPack
	var err error
	path := "testing/maketarpack.tar"
	if tarPack, err = CreateTarPack(path); err != nil {
		t.Errorf("Error opening tar pack [%s]\n", path)
		return
	}
	if tarPack == nil {
		t.Errorf("TarPack was nil but no err")
	}
}
