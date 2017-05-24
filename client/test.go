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
  "math"
	"strconv"
)

// keys
var clientPrivateKey float64
var clientPublicKey float64
var serverPublicKey float64
var commonKey float64
var sharedKey int
var connection net.Conn
var prime float64
var generator float64

// message codes
var publicKeyCode string

func setup(){
	prime = 11
	generator = 23
	publicKeyCode = "ssad990=+?A¡][ªsa)(asdª]ßðA=S)]"

	clientPrivateKey = 6
	fmt.Println("Client private key: ", clientPrivateKey)

	clientPublicKey = math.Mod(math.Pow(prime, clientPrivateKey), generator)
	fmt.Println("Client public key: ", clientPublicKey, "\n")
}

func dialServer(){
  conn, err := net.Dial("tcp", "localhost:8081")
  if err != nil {
    fmt.Println("Can't connect to server")
    return
  }
  fmt.Println("Connected to server")
  connection = conn
}

func exchangeKeys(){
	// send and reviece public key
	fmt.Println("Sending public key")
	msg := publicKeyCode + strconv.FormatFloat(clientPublicKey, 'f', -1, 64) + "\n"
	fmt.Fprintf(connection, msg)

	fmt.Println("Waiting for server public key..")

	for{
		message, _ := bufio.NewReader(connection).ReadString('\n')
		if len(message) > len(publicKeyCode){
			if message[0 : len(publicKeyCode)] == publicKeyCode{
				c, _ := strconv.ParseFloat((message[len(publicKeyCode) : len(message) - 1]), 64)
				serverPublicKey = c
				fmt.Println("Server public key recieved: " + message[len(publicKeyCode) : len(message) - 1])
				break
			}
		}
	}

	commonKey = math.Mod(math.Pow(serverPublicKey, clientPrivateKey), generator)
	commonKeyString := strconv.FormatFloat(commonKey, 'f', -1, 64)
	fmt.Println("Common key: " + commonKeyString)
}

func startClient(){
	setup()
  dialServer()
  exchangeKeys()
}

func main(){
  startClient()
}

/**
func main() {
	conn, err := net.Dial("tcp", "localhost:8081")
	if err != nil {
		fmt.Println("Can't find server.")
		return
	}
	fmt.Println("Connected to server.")

	go listener(conn, byteKey)

	for {
		// read in input from stdin
		reader := bufio.NewReader(os.Stdin)
		text, _ := reader.ReadString('\n')
		if !checkForCmd(conn, text) {
			cryptText := encrypt(byteKey, text)

			fmt.Fprintf(conn, cryptText+"\n")
		}
	}
} **/

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
