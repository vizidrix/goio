package goio

import (
	"errors"
	"testing"
	"time"
)

func Test_Should_generage_a_basic_cert_from_MakeCert(t *testing.T) {
	organization := "org"
	size := 1024
	var hosts []string
	lifespan := time.Hour * 10
	isCA := false

	cert, err := MakeCert(organization, size, hosts, lifespan, isCA)

	if err != nil {
		t.Fatalf("Error: %s", err)
	}
	if cert == nil {
		t.Fatalf("Invalid cert returned, should have been non nil err")
	}
}

var table_CertDefinitionTests = []struct {
	organization string
	size         int
	hosts        []string
	lifespan     time.Duration
	isCA         bool
	err          error
}{
	{"", 1024, []string{}, time.Minute, false, errors.New("Invalid RSA organization")},
	{"a", 1024, []string{}, time.Minute, false, nil},
	{"a", 10, []string{}, time.Minute, false, errors.New("Invalid RSA key size")},
}

func Test_Should_fail_for_invalid_cert_definition_fields(t *testing.T) {
	for i, tt := range table_CertDefinitionTests {
		def, err := NewCertDefinition(tt.organization, tt.size, tt.hosts, tt.lifespan, tt.isCA)
		bothNil := err == nil && tt.err == nil
		neitherNil := err != nil && tt.err != nil
		if bothNil || (neitherNil && (err.Error() == tt.err.Error())) {
			continue
		}
		t.Errorf("\n%d. NewCertDefinition(%q, %q, %q, %q, %q) => (%q, %q), want err %q\n", i, tt.organization, tt.size, tt.hosts, tt.lifespan, tt.isCA, def, err, tt.err)
	}
}
