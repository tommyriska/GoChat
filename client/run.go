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
	"strings"
)

func main() {
	startClient()
}

func startClient(){
	var key string
	// connect to server!
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
		if !checkForCmd(conn, text) {
			cryptText := encrypt(byteKey, text)

			// send to socket
			fmt.Fprintf(conn, cryptText+"\n")
		}
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

func quit(conn net.Conn) {
	conn.Close()
	os.Exit(1)
}

// check for command
func checkForCmd(conn net.Conn, msg string) bool {
	if len(msg) > 1 {
		words := strings.Split(msg[0:len(msg)-1], " ")
		switch words[0] {
		case "!quit":
			quit(conn)
			return true
		}
	}
	return false
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
