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
)

var key = []byte("example key 1234")

func main() {
	conn, _ := net.Dial("tcp", "localhost:8081")
	// fmt.Println("Dial finished, waiting for message from server")
	// message, _ := bufio.NewReader(conn).ReadString('\n')
	for {
		// fmt.Println("MESSAGE FROM SERVER: ", message)
		// cipherkey := []byte(message)
		// fmt.Println(cipherkey)

		// read in input from stdin
		reader := bufio.NewReader(os.Stdin)
		fmt.Print("Text to send: ")
		plaintext, _ := reader.ReadString('\n')
		fmt.Printf("Plaintext: %s\n", plaintext)
		encrypted := encrypt(key, plaintext)
		// send to socket
		fmt.Fprintf(conn, encrypted+"\n")
		fmt.Printf("Encrypted and sent to server: %s\n", encrypted)

		// listen for reply
		message, _ := bufio.NewReader(conn).ReadString('\n')
		cryptomsg := message[0 : len(message)-1]
		fmt.Printf("Servers encrypted reply: %s\n", cryptomsg)
		decmessage := decrypt(key, cryptomsg)
		fmt.Printf("Servers decrypted reply: %s\n", decmessage)
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
