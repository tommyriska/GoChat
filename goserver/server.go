package main

import (
	"bufio"
	"fmt"
	"net"
)

var clients []Client
var key []byte

type Client struct {
	connection net.Conn
}

func (c *Client) send(message []byte) {
	c.connection.Write(message)
}

func (c *Client) listener() {
	c.send([]byte(key))
	for {
		message, _ := bufio.NewReader(c.connection).ReadString('\n')
		if len(message) > 0 {
			msg := message[0 : len(message)-1]
			fmt.Println(msg)
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
	ln, _ := net.Listen("tcp", ":8081")
	fmt.Println("Server is listening on " + ln.Addr().String())
	key = []byte("example key 1234") // Must be 16,24 or 32 bytes
	fmt.Println("A new key is created!")
	for {
		conn, _ := ln.Accept()
		newClient := Client{connection: conn}
		clients = append(clients, newClient)
		fmt.Println(conn.RemoteAddr().String(), "connected.")
		fmt.Println("Connected clients:", len(clients))
		newClient.startThread()
	}
}
