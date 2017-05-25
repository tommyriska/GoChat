package main

import (
	"bufio"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"os"
	"strconv"
	"strings"

	"github.com/monnand/dhkx"
)

var commonKey []byte
var connection net.Conn
var publicKeyCode string

func setup() {
	publicKeyCode = "ssd990=+?¡][ªs)(sdª]ßð=S)]"
}

func welcome() (string, string) {
	fmt.Println("\nWelcome to GoChat!")
	fmt.Println("1. Direct connection")
	fmt.Println("2. Choose from stored servers")
	fmt.Println("3. Add new server")

	reader := bufio.NewReader(os.Stdin)
	text, _ := reader.ReadString('\n')

	var address string
	var port string

	switch string(text[0]) {
	case "1":
		address, port = chooseServer()
	case "2":
		address, port = chooseStoredServer()
	case "3":
		address, port = chooseServer()
		fmt.Print("Server name: ")

		reader := bufio.NewReader(os.Stdin)
		text, _ := reader.ReadString('\n')
		name := text[0 : len(text)-1]

		storeNewServer(address, port, name)
	}

	return address, port
}

func chooseStoredServer() (string, string) {
	var address string
	var port string
	var servers []Server

	json.Unmarshal([]byte("servers.json"), &servers)
	for i, e := range servers {
		fmt.Println(i, ". "+e.address+":"+e.port+": "+e.name)
	}

	fmt.Print("\nChoose server: ")
	reader := bufio.NewReader(os.Stdin)
	text, _ := reader.ReadString('\n')
	chosen := text[0 : len(text)-1]
	i, _ := strconv.ParseInt(chosen, 10, 64)

	address = servers[i].address
	port = servers[i].port

	return address, port
}

func storeNewServer(address string, port string, name string) {
	dat, err := ioutil.ReadFile("servers.txt")
	if err != nil {
		panic(err)
	}

	text := string(dat)
	newServer := address + "|" + port + "|" + name + "\n"
	text += newServer

	err2 := ioutil.WriteFile("servers.txt", []byte(text), 0644)
	if err2 != nil {
		panic(err2)
	}
}

func chooseServer() (string, string) {
	var address string
	var port string

	fmt.Print("Server address: ")
	reader := bufio.NewReader(os.Stdin)
	text, _ := reader.ReadString('\n')
	address = text[0 : len(text)-1]

	fmt.Print("Server port: ")
	reader2 := bufio.NewReader(os.Stdin)
	text2, _ := reader2.ReadString('\n')
	port = text2[0 : len(text2)-1]

	return address, port
}

func dialServer(address string, port string) bool {
	conn, err := net.Dial("tcp", address+":"+port)
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
	msg := publicKeyCode + string(clientPublicKey) + "\n"
	fmt.Fprintf(connection, msg)

	// listening for server public key
	for {
		message, _ := bufio.NewReader(connection).ReadString('\n')
		if len(message) > len(publicKeyCode) {
			if message[0:len(publicKeyCode)] == publicKeyCode {
				serverPublicKey = []byte(message[len(publicKeyCode) : len(message)-1])
				break
			}
		}
	}

	// finding common key
	pubKey := dhkx.NewPublicKey(serverPublicKey)
	k, _ := g.ComputeKey(pubKey, clientPrivateKey)
	commonKey = k.Bytes()[0:32]
	fmt.Println("")
}

func startClient() {
	setup()
	address, port := welcome()

	if dialServer(address, port) {
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
