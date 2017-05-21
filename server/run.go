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
	"strings"
)

// client list
var rooms []Room
var clientRoom map[Client]Room
var key []byte

type Room struct {
	name        string
	discription string
	password    string
	maxClients  int
	welcomeMsg  string
}

type Client struct {
	connection net.Conn
	nick       string
}

func (c *Client) send(message []byte) {
	c.connection.Write(message)
}

func (c *Client) sendEncrypted(message string) {
	cryptMsg := encrypt(key, message)
	c.connection.Write([]byte(cryptMsg + "\n"))
}

func (c *Client) setAndSendKey(key []byte) {
	c.send([]byte(string(key) + "\n"))
}

func (c Client) listener() {
	for {
		message, _ := bufio.NewReader(c.connection).ReadString('\n')
		if len(message) > 0 {
			msg := message[0 : len(message)-1]
			fmt.Print(c.connection.RemoteAddr().String() + ": " + decrypt(key, msg))
			if !checkForCmd(c, decrypt(key, msg)) {
				for key, value := range clientRoom {
					if key != c && value == clientRoom[c] {
						key.send([]byte(message + "\n"))
					}
				}
			}
		}
	}
}

func (c *Client) startThread() {
	go c.listener()
}

func startServer(){
	// create key
	key = createKey()

	// lobby
	makeRoom("Lobby", "Welcome to Lobby")
	makeRoom("TestRoom", "Welcome to TestRoom")

	// start server
	port := ":8081"
	ln, _ := net.Listen("tcp", port)
	fmt.Println("Server is listening on " + port)

	clientRoom = make(map[Client]Room)

	// listen loop
	for {
		conn, _ := ln.Accept()
		newClient := Client{connection: conn}
		newClient.setAndSendKey(key)
		newClient.startThread()
		fmt.Println(conn.RemoteAddr().String() + " connected.")
		switchRoom(newClient, rooms[0])
	}
}

func main() {
	startServer()
}

// check for command
func checkForCmd(client Client, msg string) bool {
	if len(msg) > 1 {
		words := strings.Split(msg[0:len(msg)-1], " ")
		switch words[0] {
		case "!room":
			if len(words) > 1 {
				for _, element := range rooms {
					if element.name == words[1] {
						client.sendEncrypted("Switching room: " + element.name + "\n")
						switchRoom(client, element)
					}
				}
			} else {
				roomMsg := ""
				for _, element := range rooms {
					roomMsg += " - " + element.name + "\n"
				}
				client.sendEncrypted(roomMsg)
			}
			return true
		}
	}
	return false
}

func switchRoom(client Client, room Room) {
	clientRoom[client] = room
	client.sendEncrypted(room.welcomeMsg + "\n")
}

func makeRoom(name string, welcomeMsg string) {
	newRoom := Room{name: name, welcomeMsg: welcomeMsg}
	rooms = append(rooms, newRoom)
}

// create key
func createKey() []byte {
	randKey := make([]byte, 32)
	_, err := rand.Read(randKey)
	if err != nil {
		panic(err)
	}
	return randKey
}

// encrypt message
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

// decrypt message
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
