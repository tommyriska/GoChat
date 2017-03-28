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
)

var clients []Client
var key = []byte("example key 1234")

type Client struct {
	connection net.Conn
}

func main() {
	ln, _ := net.Listen("tcp", ":8081")
	fmt.Println("Server is listening on port 8081..")
	for {
		conn, _ := ln.Accept()
		newClient := Client{connection: conn}
		clients = append(clients, newClient)
		fmt.Println(conn.RemoteAddr().String(), "connected.")
		fmt.Println("Connected clients:", len(clients))
		newClient.startThread()
	}
}

func (c *Client) send(message []byte) {
	c.connection.Write(message)
}

func (c *Client) listener() {
	for {
		message, _ := bufio.NewReader(c.connection).ReadString('\n')
		if len(message) > 0 {
			cryptomsg := message[0 : len(message)-1]
			fmt.Printf("Encrypted message from %s: %s\n", c.connection.RemoteAddr().String(), cryptomsg)
			decmessage := decrypt(key, cryptomsg)
			fmt.Printf("Decrypted message: %s\n", decmessage)
			encryptedmessage := encrypt(key, decmessage)
			fmt.Printf("Message sent to clients: %s\n", encryptedmessage)
			for _, element := range clients {
				element.send([]byte(encryptedmessage + "\n"))
			}
		}
	}
}

func (c *Client) startThread() {
	go c.listener()
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
