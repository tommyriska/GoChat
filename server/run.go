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

// Room is a struct keeping track of room details
type Room struct {
	name        string
	discription string
	password    string
	maxClients  int
	welcomeMsg  string
}

// Client represents a connected client
type Client struct {
	connection net.Conn
	nick       string
}

// send is a function for sending a plaintext message
func (c *Client) send(message []byte) {
	c.connection.Write(message)
}

// sendEncrypted is a function for sending an encrypted message
func (c *Client) sendEncrypted(message string) {
	cryptMsg := encrypt(key, message)
	c.connection.Write([]byte(cryptMsg + "\n"))
}

// setAndSendKey is a function for sending the
// encryption key to a new connection
func (c *Client) setAndSendKey(key []byte) {
	c.send([]byte(string(key) + "\n"))
}

// listener is a function which runs in the background as a
// go routine which listens for incoming messages.
// As of now, for display purposes, the message gets decrypted and displayed
// in plaintext in the server window. Then it is checked for command words with
// the checkForCmd function. If that returns false the message gets sent
// to all clients in the same room as the sender. At the end, variable "message"
// is still decrypted, therefore it is sent using only the send function.
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

// Starts a goroutine that listens for incoming messages.
func (c *Client) startThread() {
	go c.listener()
}

func main() {
	// Upon startup we create a key which is used to encrypt and decrypt all messages
	// with. This is distributed to all new connections upon connection. This key
	// is 32 bytes long(256 bits) which is still concidered safe and not broken.
	key = createKey()

	// Then we set up default rooms. Every new conncetion will be set to the
	// lobby room.
	makeRoom("Lobby", "Welcome to Lobby")
	makeRoom("TestRoom", "Welcome to TestRoom")

	// Then we set up the socket to be a TCP connection and to listen to a
	// specified port. This port can be set to all port from 1024 and up.
	port := ":8081"
	ln, _ := net.Listen("tcp", port)
	fmt.Println("Server is listening on " + port)

	// Here we set up a map which keeps track of which room all clients are connected to.
	clientRoom = make(map[Client]Room)

	// This loop listens for new connections. It accepts the new connection, creates
	// a new Client struct, sends the encryption/decryption key and starts a listener
	// thread which listens for incoming messages. Then it prints out in the
	// server window that a new client has connected and sets the clients room
	// to lobby, which is the default landing room.
	for {
		conn, _ := ln.Accept()
		newClient := Client{connection: conn}
		newClient.setAndSendKey(key)
		newClient.startThread()
		fmt.Println(conn.RemoteAddr().String() + " connected.")
		switchRoom(newClient, rooms[0])
	}
}

// Function to check for commands in a message sent by a client
func checkForCmd(client Client, msg string) bool {
	if len(msg) > 1 {
		words := strings.Split(msg[0:len(msg)-1], " ")
		switch words[0] {
		// The room command is used to both list all rooms and to change rooms.
		// !room sends a list over all rooms to the client
		// !room [roomname] is a command used to change room to the selected room
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
		// The !newroom command lets clients create new rooms
		case "!newroom":
			if len(words) > 1 {
				name := words[1]
				welcomeMsg := words[2]
				for _, element := range rooms {
					if element.name == name {
						client.sendEncrypted("A room with this name already exists!\n")
					} else {
						makeRoom(name, welcomeMsg)
						for _, element := range rooms {
							if element.name == words[1] {
								client.sendEncrypted("Switching room: " + element.name + "\n")
								switchRoom(client, element)
							}
						}
					}
				}

			}
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
