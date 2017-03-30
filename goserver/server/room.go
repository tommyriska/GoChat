package main

import "fmt"

type Room struct {
	roomList    []Room
	clients     []Client
	name        string
	discription string
	password    string
	hasPassword bool
	maxClients  int
	key         []byte
	welcomeMsg  string
}

func (r Room) addClient(c Client) {
	r.clients = append(r.clients, c)
	c.room = r
	c.clients = r.clients
	c.setAndSendKey(r.key)
	c.startThread()
	fmt.Println(c.connection.RemoteAddr().String(), "connected.")
	fmt.Println("Connected clients:", len(r.clients))
}
