package goio

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"errors"
	"fmt"
	"math/big"
	"net"
	"os"
	"time"
)

type Cert struct {
	privateKey *rsa.PrivateKey
	data       []byte
}

type CertDefinition struct {
	Organization string
	Size         int
	Hosts        []string
	LifeSpan     time.Duration
	IsCA         bool
}

func intInSlice(value int, list []int) bool {
	for _, v := range list {
		if v == value {
			return true
		}
	}
	return false
}

func NewCertDefinition(organization string, size int, hosts []string, lifespan time.Duration, isCA bool) (*CertDefinition, error) {
	if organization == "" {
		return nil, errors.New("Invalid RSA organization")
	}
	validSizes := []int{512, 1024, 2048, 3072, 7680, 15360}
	if !intInSlice(size, validSizes) {
		return nil, errors.New("Invalid RSA key size")
	}
	return &CertDefinition{
		organization,
		size,
		hosts,
		lifespan,
		isCA,
	}, nil
}

func MakeCert(organization string, size int, hosts []string, lifespan time.Duration, isCA bool) (*Cert, error) {
	var err error
	var cert Cert
	if cert.privateKey, err = rsa.GenerateKey(rand.Reader, size); err != nil {
		return nil, errors.New(fmt.Sprintf("Unable to generate key [%s]\n", err))
	}

	// Setup certificate expiration
	notBefore := time.Now()
	notAfter := notBefore.Add(lifespan)

	// Setup certificate configuration
	template := x509.Certificate{
		SerialNumber: new(big.Int).SetInt64(0),
		Subject: pkix.Name{
			Organization: []string{organization},
		},
		NotBefore:             notBefore,
		NotAfter:              notAfter,
		KeyUsage:              x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		BasicConstraintsValid: true,
	}

	// Separate valid IP addresses from the named hosts provided
	for _, host := range hosts {
		if ip := net.ParseIP(host); ip != nil {
			template.IPAddresses = append(template.IPAddresses, ip)
		} else {
			template.DNSNames = append(template.DNSNames, host)
		}
	}

	// Append usage for Cert Signing if CA is specified
	if isCA {
		template.IsCA = true
		template.KeyUsage |= x509.KeyUsageCertSign
	}

	// Make certificate blob
	if cert.data, err = x509.CreateCertificate(rand.Reader, &template, &template, &cert.privateKey.PublicKey, cert.privateKey); err != nil {
		return nil, errors.New(fmt.Sprintf("Unable to create certificate [%s]\n", err))
	}

	return &cert, nil
}

//"certs/private.pem"
func (cert *Cert) WritePrivate(file string) error {
	var target *os.File
	var err error
	if target, err = os.OpenFile(file, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600); err != nil {
		return errors.New(fmt.Sprintf("Unable to create private cert file [%s]\n", err))
	}
	defer target.Close()
	if err := pem.Encode(target, &pem.Block{Type: "CERTIFICATE", Bytes: cert.data}); err != nil {
		return errors.New(fmt.Sprintf("Unable to encode private cert [%s]\n", err))
	}
	return nil
}

//"certs/public.pem"
func (cert *Cert) WritePublic(file string) error {
	var target *os.File
	var err error
	if target, err = os.OpenFile(file, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600); err != nil {
		return errors.New(fmt.Sprintf("Unable to create public cert file [%s]\n", err))
	}
	defer target.Close()
	if err = pem.Encode(target, &pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(cert.privateKey)}); err != nil {
		return errors.New(fmt.Sprintf("Unable to encode public cert [%s]\n", err))
	}
	return nil
}

/*
func GetWindowsMachineGuid() (guid string, err error) {
	var h syscall.Handle
	err = syscall.RegOpenKeyEx(syscall.HKEY_LOCAL_MACHINE, syscall.StringToUTF16Ptr(`SOFTWARE\Microsoft\Cryptography`), 0, syscall.KEY_READ, &h)
	if err != nil {
		return
	}
	defer syscall.RegCloseKey(h)
	var typ uint32
	var buf [74]uint16
	n := uint32(len(buf))
	err = syscall.RegQueryValueEx(h, syscall.StringToUTF16Ptr("MachineGuid"), nil, &typ, (*byte)(unsafe.Pointer(&buf[0])), &n)
	if err != nil {
		return
	}
	guid = syscall.UTF16ToString(buf[:])
	return
}
*/
