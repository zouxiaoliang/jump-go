package main

import (
	"crypto/rc4"
	"fmt"
	"testing"
)

func TestMain(t *testing.T) {
	t.Fatal("This is a placeholder for the actual test function.")
}

func TestRc4(t *testing.T) {
	var content = []byte("1")
	fmt.Printf("%v\n", content)
	var key = []byte("fFa2dhdqBHZKpg61JsAX")
	var ic, _ = rc4.NewCipher(key)
	var oc, _ = rc4.NewCipher(key)

	var encrypted = make([]byte, len(content))
	ic.XORKeyStream(encrypted, content)
	fmt.Printf("%v\n", encrypted)

	var decrypted = make([]byte, len(content))
	oc.XORKeyStream(encrypted, encrypted)
	fmt.Printf("%v\n", decrypted)

}
