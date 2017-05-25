package main

import (
	"bufio"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"dhkx"
	"encoding/base64"
	"fmt"
	"io"
	"net"
	"os"
	"strings"
)

<<<<<<< HEAD
var commonKey []byte
var connection net.Conn
var publicKeyCode string

func setup() {
	publicKeyCode = "ssd990=+?¡][ªs)(sdª]ßð=S)]"
}

func dialServer() bool {
	conn, err := net.Dial("tcp", "localhost:8081")
	if err != nil {
		fmt.Println("Can't connect to server")
		return false
	}

	fmt.Println("\nConnected to server")
	connection = conn
	return true
}

func contains(s []byte, e byte) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

func exchangeKeys() {
	var serverPublicKey []byte

	// generate private key
	g, _ := dhkx.GetGroup(0)
	clientPrivateKey, _ := g.GeneratePrivateKey(nil)

	// make sure the key does not contain '\n' or '%'
	for {
		if contains(clientPrivateKey.Bytes(), byte('\n')) || contains(clientPrivateKey.Bytes(), byte('%')) {
			newKey, _ := g.GeneratePrivateKey(nil)
			clientPrivateKey = newKey
		} else {
			break
		}
	}

	// generate public key
	clientPublicKey := clientPrivateKey.Bytes()

	// sending client public key
	fmt.Println("Sending client public key")
	msg := publicKeyCode + string(clientPublicKey) + "\n"

	fmt.Fprintf(connection, msg)

	// listening for server public key
	fmt.Println("Waiting for server public key..")
	for {
		message, _ := bufio.NewReader(connection).ReadString('\n')
		if len(message) > len(publicKeyCode) {
			if message[0:len(publicKeyCode)] == publicKeyCode {
				serverPublicKey = []byte(message[len(publicKeyCode) : len(message)-1])
				fmt.Println("Server public key recieved")
				break
			}
		}
	}

	// finding common key
	fmt.Println("Finding common key")
	pubKey := dhkx.NewPublicKey(serverPublicKey)
	k, _ := g.ComputeKey(pubKey, clientPrivateKey)
	commonKey = k.Bytes()[0:32]
	fmt.Println("Key exchange complete")
	fmt.Println("Common key: " + string(commonKey))
	fmt.Println("")
}

func startClient() {
	setup()
	if dialServer() {
		exchangeKeys()

		go listener(connection, commonKey)

		for {
			reader := bufio.NewReader(os.Stdin)
			text, _ := reader.ReadString('\n')
			if !checkForCmd(connection, text) {
				cryptText := encrypt(commonKey, text)

				fmt.Fprintf(connection, cryptText+"\n")
			}
		}
	}
}

func main() {
	startClient()
}

func listener(conn net.Conn, key []byte) {
	for {
		message, _ := bufio.NewReader(conn).ReadString('\n')
		msg := decrypt(key, message)
		fmt.Print(msg)
	}
}

func quit(conn net.Conn) {
	conn.Close()
	os.Exit(1)
}

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
