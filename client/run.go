package main

import (
	"bufio"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"io"
	"net"
	"os"

	"github.com/jroimartin/gocui"
)

func main() {
	g, _ := gocui.NewGui(gocui.OutputNormal)

	var key string
	// connect to server
	conn, err := net.Dial("tcp", "localhost:8081")
	if err != nil {
		fmt.Println("Can't find server.")
		return
	}
	fmt.Println("Connected to server.")
	// get key
	key, _ = bufio.NewReader(conn).ReadString('\n')

	// get key
	keyMsg := []byte(key)
	byteKey := keyMsg[0 : len(keyMsg)-1]

	// start listener thread
	go listener(conn, byteKey)

	for {
		// read in input from stdin
		reader := bufio.NewReader(os.Stdin)
		text, _ := reader.ReadString('\n')
		cryptText := encrypt(byteKey, text)

		// send to socket
		fmt.Fprintf(conn, cryptText+"\n")
	}
}

func listener(conn net.Conn, key []byte) {
	for {
		// listen for message from server
		message, _ := bufio.NewReader(conn).ReadString('\n')
		msg := decrypt(key, message)
		fmt.Print(msg)
	}
}

func encrypt(key []byte, text string) string {
	plaintext := []byte(text)

	block, err := aes.NewCipher(key)
	if err != nil {
		panic(err)
	}

	ciphertext := make([]byte, aes.BlockSize+len(plaintext))
	iv := ciphertext[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		panic(err)
	}

	stream := cipher.NewCFBEncrypter(block, iv)
	stream.XORKeyStream(ciphertext[aes.BlockSize:], plaintext)

	return base64.URLEncoding.EncodeToString(ciphertext)
}

func decrypt(key []byte, cryptoText string) string {
	ciphertext, _ := base64.URLEncoding.DecodeString(cryptoText)

	block, err := aes.NewCipher(key)
	if err != nil {
		panic(err)
	}
	if len(ciphertext) < aes.BlockSize {
		panic("Ciphertext too short")
	}

	iv := ciphertext[:aes.BlockSize]
	ciphertext = ciphertext[aes.BlockSize:]

	stream := cipher.NewCFBDecrypter(block, iv)

	stream.XORKeyStream(ciphertext, ciphertext)

	return fmt.Sprintf("%s", ciphertext)
}
