package main

import (
	"bufio"
	"fmt"
	"net"
)

var clients []Client

type Client struct {
	connection net.Conn
}

func (c *Client) send(message []byte) {
	c.connection.Write(message)
}

func (c *Client) listener() {
	c.send([]byte("test" + "\n"))
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

	for {
		conn, _ := ln.Accept()
		newClient := Client{connection: conn}
		clients = append(clients, newClient)
		fmt.Println(conn.RemoteAddr().String(), "connected.")
		fmt.Println("Connected clients:", len(clients))
		newClient.startThread()
	}
}
