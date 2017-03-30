package main

import (
	"bufio"
	"fmt"
	"net"
)

type Client struct {
	connection net.Conn
	clients    []Client
	key        []byte
	isInRoom   bool
	nick       string
	room       Room
}

func (c *Client) send(message []byte) {
	c.connection.Write(message)
}

func (c *Client) sendEncrypted(message string) {
	cryptMsg := encrypt(c.key, message)
	c.connection.Write([]byte(cryptMsg + "\n"))
}

func (c *Client) setAndSendKey(key []byte) {
	c.key = key
	c.send([]byte(string(key) + "\n"))
}

func (c Client) listener() {
	c.sendEncrypted(c.room.welcomeMsg + "\n")
	for {
		message, _ := bufio.NewReader(c.connection).ReadString('\n')
		if len(message) > 0 {
			msg := message[0 : len(message)-1]
			fmt.Print(c.connection.RemoteAddr().String() + ": " + decrypt(c.key, msg))
			if !checkForCmd(c, decrypt(c.key, msg)) {
				for _, element := range c.clients {
					if element.connection != c.connection {
						element.send([]byte(message + "\n"))
					}
				}
			}
		}

	}
}

func (c *Client) startThread() {
	go c.listener()
}
