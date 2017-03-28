package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
)

func main() {
	var key string
	// connect to server
	conn, _ := net.Dial("tcp", "158.37.63.27:8081")

	// get key
	key, _ = bufio.NewReader(conn).ReadString('\n')
	fmt.Print("Key: " + key)

	// get key
	keyMsg := []byte(key)
	byteKey := keyMsg[0 : len(keyMsg)-1]

	// start listener thread
	go listener(conn, byteKey)

	for {
		// read in input from stdin
		reader := bufio.NewReader(os.Stdin)
		text, _ := reader.ReadString('\n')
		cryptText := encrypt(byteKey, text)

		// send to socket
		fmt.Fprintf(conn, cryptText+"\n")
	}
}

func listener(conn net.Conn, key []byte) {
	for {
		// listen for message from server
		message, _ := bufio.NewReader(conn).ReadString('\n')
		msg := decrypt(key, message)
		fmt.Print(msg)
	}
}
