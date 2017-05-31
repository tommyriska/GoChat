package main

import (
	"bufio"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"os"
	"os/exec"
	"runtime"
	"strconv"
	"strings"

	"github.com/monnand/dhkx"
)

// struct for holding server info
type Server struct {
	address string
	port    string
	name    string
}

var commonKey []byte
var connection net.Conn
var publicKeyCode string
var nickCode string
var nick string

// init variables
func setup() {
	publicKeyCode = "ssd990=+?¡][ªs)(sdª]ßð=S)]"
	nickCode = "!#28jKas>zzx'**!+?,>lzc012"
}

// clear terminal screen
func clear() {
	switch runtime.GOOS {
	case "linux":
		cmd := exec.Command("clear")
		cmd.Stdout = os.Stdout
		cmd.Run()
	case "windows":
		cmd := exec.Command("cmd", "/c", "cls")
		cmd.Stdout = os.Stdout
		cmd.Run()
	default:
		fmt.Println("Attempted to clear terminal, but OS is not supported.")
	}
}

// prompt user for nickname and return it
func chooseNick() string {
	fmt.Print("Nickname: ")
	reader := bufio.NewReader(os.Stdin)
	text, _ := reader.ReadString('\n')
	nickname := text[0 : len(text)-1]

	return nickname
}

// print welcome message and return chosen action
func welcomePrompt() string {
	clear()

	// print user choices
	fmt.Println("Welcome to GoChat!\n")
	fmt.Println("1 Direct connection")
	fmt.Println("2 Choose from stored servers")
	fmt.Println("3 Add new server\n")

	// prompt for input
	reader := bufio.NewReader(os.Stdin)
	text, _ := reader.ReadString('\n')

	// return the chosen action
	return string(text[0])
}

// choose server, return address and port
func chooseServer(choice string) (string, string) {
	clear()

	var address string
	var port string

	switch choice {
	case "1":
		address, port = chooseDirectServer()
	case "2":
		address, port = chooseStoredServer()
	case "3":
		address, port = chooseDirectServer()

		// promt for server name
		fmt.Print("Server name: ")
		reader := bufio.NewReader(os.Stdin)
		text, _ := reader.ReadString('\n')
		name := text[0 : len(text)-1]

		// save new server
		storeNewServer(address, port, name)
	}

	return address, port
}

// print saved servers and let user choose one, return the chosen servers address and port
func chooseStoredServer() (string, string) {
	var address string
	var port string
	var serverArray []Server

	clear()

	// read data from server file
	dat, err := ioutil.ReadFile("servers.txt")
	if err != nil {
		panic(err)
	}

	// split servers and save in string slice
	servers := strings.Split(string(dat), "|")

	// print servers, and create a struct for each server, add them to the serverArray
	fmt.Println("Servers:\n")
	for i, e := range servers {
		if len(e) > 1 {
			data := strings.Split(e, "-")
			s := Server{address: data[0], port: data[1], name: data[2]}
			fmt.Println(i, s.address+":"+s.port+" "+s.name)
			serverArray = append(serverArray, s)
		}
	}

	// prompt for server choice
	fmt.Print("\nChoose server: ")
	reader := bufio.NewReader(os.Stdin)
	text, _ := reader.ReadString('\n')
	chosen := text[0 : len(text)-1]
	i, _ := strconv.ParseInt(chosen, 10, 64)

	// check if choice is valid
	if int(i) < len(serverArray) {
		address = serverArray[i].address
		port = serverArray[i].port
	}

	return address, port
}

// add new server info to servers.txt
func storeNewServer(address string, port string, name string) {
	// read data from server file
	dat, err := ioutil.ReadFile("servers.txt")
	if err != nil {
		panic(err)
	}

	// convert []byte to string
	text := string(dat)
	// create string representing the new server
	// "|" seperates servers and "-" seperates the server attributes
	newServer := address + "-" + port + "-" + name + "|"
	// append the new server to the servers.txt-content
	text += newServer

	// write the updated content to file
	err2 := ioutil.WriteFile("servers.txt", []byte(text), 0644)
	if err2 != nil {
		panic(err2)
	}
}

// lets user choose address and port to connect to
func chooseDirectServer() (string, string) {
	var address string
	var port string

	clear()

	// prmpt for address
	fmt.Print("Server address: ")
	reader := bufio.NewReader(os.Stdin)
	text, _ := reader.ReadString('\n')
	address = text[0 : len(text)-1]

	// prompt for port
	fmt.Print("Server port: ")
	reader2 := bufio.NewReader(os.Stdin)
	text2, _ := reader2.ReadString('\n')
	port = text2[0 : len(text2)-1]

	return address, port
}

// dial server
func dialServer(address string, port string) bool {
	conn, err := net.Dial("tcp", address+":"+port)
	if err != nil {
		fmt.Println("Can't connect to server")
		return false
	}

	connection = conn
	return true
}

// check if byte slice contains a specific byte
func contains(s []byte, e byte) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

// exchange keys with server
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
}

// send nickname to server
func sendNick() {
	fmt.Fprintf(connection, encrypt(commonKey, nickCode+nick)+"\n")
}

// make string bold
func makeBold(text string) string {
	return "\033[1m" + text + "\033[0m"
}

// start the client
func startClient() {
	// init
	setup()
	// find wich address and port to connect to
	serverChoice := welcomePrompt()
	address, port := chooseServer(serverChoice)
	clear()
	// find the chosen nickname
	nick = chooseNick()
	clear()

	// if establish connection to server
	if dialServer(address, port) {
		// exchange keys
		exchangeKeys()
		sendNick()
		fmt.Println("Connected to: " + address + ":" + port)

		// start thread to listen for messages from server
		go listener(connection, commonKey)

		// listen for input from user
		for {
			reader := bufio.NewReader(os.Stdin)
			text, _ := reader.ReadString('\n')

			// check user input for commands, if no commands encrypt and send to server
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

// listens for messages from server, decrypts and prints them
func listener(conn net.Conn, key []byte) {
	for {
		message, _ := bufio.NewReader(conn).ReadString('\n')
		msg := decrypt(key, message)
		fmt.Print(msg)
	}
}

// close connection to server and quit client
func quit(conn net.Conn) {
	conn.Close()
	os.Exit(1)
}

// check for commands in user input, return true if command is found
func checkForCmd(conn net.Conn, msg string) bool {
	if len(msg) > 1 {
		words := strings.Split(msg[0:len(msg)-1], " ")
		switch words[0] {
		// !quit is the command for quitting
		case "!quit":
			fmt.Fprintf(connection, encrypt(commonKey, "!quit")+"\n")
			quit(conn)
			return true
		}
	}
	return false
}

// encrypt message
func encrypt(key []byte, text string) string {
	// []byte of text to encrypt
	plaintext := []byte(text)

	// get cipher from key
	block, err := aes.NewCipher(key)
	if err != nil {
		panic(err)
	}

	// generating initialization vector
	ciphertext := make([]byte, aes.BlockSize+len(plaintext))
	iv := ciphertext[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		panic(err)
	}

	// encrypting
	stream := cipher.NewCFBEncrypter(block, iv)
	stream.XORKeyStream(ciphertext[aes.BlockSize:], plaintext)

	// return encrypted text
	return base64.URLEncoding.EncodeToString(ciphertext)
}

// decrypt message
func decrypt(key []byte, cryptoText string) string {
	// get ciphertext
	ciphertext, _ := base64.URLEncoding.DecodeString(cryptoText)

	// get cipher from key
	block, err := aes.NewCipher(key)
	if err != nil {
		panic(err)
	}

	// check length
	if len(ciphertext) < aes.BlockSize {
		panic("Ciphertext too short")
	}

	// generate initialization vector
	iv := ciphertext[:aes.BlockSize]
	ciphertext = ciphertext[aes.BlockSize:]

	// decrypting
	stream := cipher.NewCFBDecrypter(block, iv)
	stream.XORKeyStream(ciphertext, ciphertext)

	// return decrypted text
	return fmt.Sprintf("%s", ciphertext)
}
