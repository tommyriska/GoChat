package main

import (
	"testing"
)

/* test will check if a string variable is altered after
method call (encrypted) */
func TestEncrypt(t *testing.T) {

	key := make([]byte, 32)
	var crypt string = "test crypto"
	text := encrypt(key, crypt)

	if text == crypt {
		t.Error("Expected new text, got same.")
	}
}

/* test will check if a string variable is altered after
method call (decrypted) */
func TestDecrypt(t *testing.T) {

	key := make([]byte, 32)
	var crypt string = "eL7Abmpadmk8nAaVcpT-a6MziwEKcL5z3ifS54SsnxA="
	text := decrypt(key, crypt)

	if text == crypt {
		t.Error("Expected new text, got same")
	}
}

/* test will check if a new server is added */
func TestNewServer(t *testing.T) {

	var a string = "testAdr"
	var p string = "testPort"
	var n string = "testName"
	var b bool = false

	storeNewServer(a, p, n)

	dat, err := ioutil.ReadFile("servers.txt")
	if err != nil {
		panic(err)
	}
	text := string(dat)
	if strings.Contains(text, a) {
		fmt.Println("Server added succesfully")
		b = true
	}
	if b == false {
		fmt.Println("Server not added")
	}
}
