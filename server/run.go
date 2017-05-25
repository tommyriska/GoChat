package main

import (
	"bufio"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"github.com/monnand/dhkx"
	"encoding/base64"
	"fmt"
	"io"
	"net"
	"strings"
)

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
	clientKey  string
}

var rooms []Room
var clientRoom map[Client]Room
var publicKeyCode string

func setup() {
	publicKeyCode = "ssd990=+?¡][ªs)(sdª]ßð=S)]"
	clientRoom = make(map[Client]Room)
	makeRoom("Lobby", "Welcome to Lobby")
	makeRoom("TestRoom", "Welcome to TestRoom")
}

func contains(s []byte, e byte) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

func exchangeKeys(c Client) {
	fmt.Println("\n" + c.connection.RemoteAddr().String() + " is trying to connect")

	var clientPublicKey []byte

	// generate private key
	g, _ := dhkx.GetGroup(0)
	serverPrivateKey, _ := g.GeneratePrivateKey(nil)

	// make sure the key does not contain '\n' or '%'
	for {
		if contains(serverPrivateKey.Bytes(), byte('\n')) || contains(serverPrivateKey.Bytes(), byte('%')) {
			newKey, _ := g.GeneratePrivateKey(nil)
			serverPrivateKey = newKey
		} else {
			break
		}
	}

	// generate public key
	serverPublicKey := serverPrivateKey.Bytes()

	fmt.Println("Waiting for client public key..")
	// listening for client public key
	for {
		message, _ := bufio.NewReader(c.connection).ReadString('\n')
		if len(message) > len(publicKeyCode) {
			if message[0:len(publicKeyCode)] == publicKeyCode {
				clientPublicKey = []byte(message[len(publicKeyCode) : len(message)-1])
				fmt.Println("Client public key recieved")
				break
			}
		}
	}

	// sending server public key
	fmt.Println("Sending server public key")
	msg := publicKeyCode + string(serverPublicKey) + "\n"
	fmt.Fprintf(c.connection, msg)

	// finding common key
	fmt.Println("Finding common key")
	pubKey := dhkx.NewPublicKey(clientPublicKey)
	k, _ := g.ComputeKey(pubKey, serverPrivateKey)
	c.clientKey = string(k.Bytes()[0:32])
	fmt.Println("Key exchange complete")
	fmt.Println("Common key: ", c.clientKey)

	c.startThread()
	fmt.Println(c.connection.RemoteAddr().String() + " connected")
	switchRoom(c, rooms[0])
}

func main() {
	setup()
	startServer()
}

func startServer() {
	port := ":8081"
	ln, _ := net.Listen("tcp", port)
	fmt.Println("Server is listening on " + port)

	for {
		conn, _ := ln.Accept()
		newClient := Client{connection: conn}
		go exchangeKeys(newClient)
	}
}

func (c Client) listener() {
	for {
		message, _ := bufio.NewReader(c.connection).ReadString('\n')
		if len(message) > 0 {
			msg := message[0 : len(message)-1]
			msgDecrypted := decrypt([]byte(c.clientKey), msg)
			fmt.Print(c.connection.RemoteAddr().String() + ": " + msgDecrypted)

			if !checkForCmd(c, msgDecrypted) {
				for mapKey, value := range clientRoom {
					if mapKey != c && value == clientRoom[c] {
						mapKey.sendEncrypted(msgDecrypted)
					}
				}
			}
		}
	}
}

func (c *Client) send(message []byte) {
	c.connection.Write(message)
}

func (c *Client) sendEncrypted(message string) {
	cryptMsg := encrypt([]byte(c.clientKey), message)
	c.connection.Write([]byte(cryptMsg + "\n"))
}

func (c *Client) startThread() {
	go c.listener()
}

func checkForCmd(client Client, msg string) bool {
	if len(msg) > 1 {
		words := strings.Split(msg[0:len(msg)-1], " ")
		switch words[0] {
		case "!room":
			if len(words) > 1 {
				if words[1] == clientRoom[client].name {
					message := "You are already in this room!\nType !room to get a list of other available chatrooms"
					client.sendEncrypted(message)
				} else {
					for _, element := range rooms {
						if element.name == words[1] {
							client.sendEncrypted("Switching room: " + element.name + "\n")
							switchRoom(client, element)
						}
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

func createKey() []byte {
	randKey := make([]byte, 32)
	_, err := rand.Read(randKey)
	if err != nil {
		panic(err)
	}
	return randKey
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
