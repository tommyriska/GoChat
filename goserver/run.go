package main

import (
	"bufio"
	"fmt"
	"net"
)

var clients []Client

type Client struct {
	connection net.Conn
	key        []byte
}

func (c *Client) send(message []byte) {
	c.connection.Write(message)
}

func (c *Client) listener() {
	// send key
	c.send([]byte(string(c.key) + "\n"))
	// listen for msg
	for {
		message, _ := bufio.NewReader(c.connection).ReadString('\n')
		if len(message) > 0 {
			msg := message[0 : len(message)-1]
			fmt.Println(c.connection.RemoteAddr().String() + ": " + msg)
			fmt.Print(c.connection.RemoteAddr().String() + ": " + decrypt(c.key, msg))
			for _, element := range clients {
				if element.connection != c.connection {
					element.send([]byte(message + "\n"))
				}
			}
		}
	}
}

func (c *Client) startThread() {
	go c.listener()
}

func main() {
	// start server
	port := ":8081"
	ln, _ := net.Listen("tcp", port)
	fmt.Println("Server is listening on " + port)

	// generate key
	key := createKey()

	// accept clients
	for {
		conn, _ := ln.Accept()
		newClient := Client{connection: conn, key: key}
		clients = append(clients, newClient)
		fmt.Println(conn.RemoteAddr().String(), "connected.")
		fmt.Println("Connected clients:", len(clients))
		newClient.startThread()
	}
}
