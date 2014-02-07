package main

import (
	"github.com/vizidrix/goio"
	"time"
)

func main() {
	MakeCert()
}

func MakeCert() {
	var err error
	var cert *goio.Cert
	private := "private.pem"
	public := "public.pem"
	if cert, err = goio.MakeCert("intel", 2048, []string{"localhost"}, 30*time.Minute, false); err != nil {
		panic(err)
	}
	cert.WritePrivate(private)
	cert.WritePublic(public)
}
