goio
======

IO helper utils for Golang

## Getting Started ##

1\. Add the import to your project

````go
import (
	. "github.com/vizidrix/goio"
)
````

2\. Start using

````
cert, err := MakeCert("Acme Inc.", 1024, []string{}, time.Minute * 10, false)
````

````
// 32 char key code
key := []byte("1a2a3a4a5a 1a2a3a4a5a 1a2a3a4a5a")
data := []byte("some secret information to encode")

// Do encrypt and decrypt
encrypted, _ := AesEncrypt(key, data)
decrypted, _ := AesDecrypt(key, encrypted)
````

----

Version
----
0.1.1 ish

Tech
----

* [Go] - Golang.org
* [GOIO] - Utility methods for Golang to help with Zip, Crypto, etc

License
----

https://github.com/vizidrix/goio/blob/master/LICENSE

----
## Edited
* 27-January-2014		refactoring AES to include Reader and Writer interfaces
* 27-January-2014		initial release

----
## Credits
* Vizidrix <https://github.com/organizations/vizidrix>
* Perry Birch <https://github.com/PerryBirch>
