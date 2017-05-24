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
	"strconv"
	"math"
)

var rooms []Room
var clientRoom map[Client]Room
var key []byte

var publicKeyCode string
var serverPrivateKey float64
var serverPublicKey float64
var clientPublicKey float64
var commonKey float64
var prime float64
var generator float64

type Room struct {
	name        string
	discription string
	password    string
	maxClients  int
	welcomeMsg  string
}

func setup(){
	publicKeyCode = "ssad990=+?A¡][ªsa)(asdª]ßðA=S)]"
	serverPrivateKey = 3
	prime = 11
	generator = 23
	serverPublicKey = math.Mod(math.Pow(prime, serverPrivateKey), generator)
	fmt.Println("Server private key: ", serverPrivateKey)
	fmt.Println("Server public key: ", serverPublicKey)
}

func exchangeKeys(c Client){

	fmt.Println("\n" + c.connection.RemoteAddr().String() + " is trying to connect")
	fmt.Println("Waiting for client public key..")

	for{
		message, _ := bufio.NewReader(c.connection).ReadString('\n')
		if len(message) > len(publicKeyCode){
			if message[0 : len(publicKeyCode)] == publicKeyCode{
				c, _ := strconv.ParseFloat((message[len(publicKeyCode) : len(message) - 1]), 64)
				clientPublicKey = c
				fmt.Println("Client public key recieved: " + message[len(publicKeyCode) : len(message) - 1])
				break
			}
		}
	}

	fmt.Println("Sending server public key")
	msg := publicKeyCode + strconv.FormatFloat(serverPublicKey, 'f', -1, 64) + "\n"
	fmt.Fprintf(c.connection, msg)

	commonKey = math.Mod(math.Pow(clientPublicKey, serverPrivateKey), generator)
	commonKeyString := strconv.FormatFloat(commonKey, 'f', -1, 64)
	fmt.Println("Common key: " + commonKeyString)
}

func main() {
	key = createKey()

	makeRoom("Lobby", "Welcome to Lobby")
	makeRoom("TestRoom", "Welcome to TestRoom")

	port := ":8081"
	ln, _ := net.Listen("tcp", port)
	fmt.Println("Server is listening on " + port)
	setup()

	clientRoom = make(map[Client]Room)

	for {
		conn, _ := ln.Accept()
		newClient := Client{connection: conn}
		exchangeKeys(newClient)
		/**
		newClient.setAndSendKey(key)
		newClient.startThread()
		fmt.Println(conn.RemoteAddr().String() + " connected.")
		switchRoom(newClient, rooms[0])**/
	}
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
