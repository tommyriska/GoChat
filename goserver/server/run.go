package main

import (
	"fmt"
	"net"
)

func main() {
	msg := "Welcome to Lobby"
	lobby := Room{name: "Lobby", hasPassword: false, maxClients: 0, key: createKey(), welcomeMsg: msg}
	test := Room{name: "Test", hasPassword: false, key: createKey()}
	lobby.roomList = append(lobby.roomList, test)

	// start server
	port := ":8081"
	ln, _ := net.Listen("tcp", port)
	fmt.Println("Server is listening on " + port)

	for {
		conn, _ := ln.Accept()
		newClient := Client{connection: conn, isInRoom: false}
		lobby.addClient(newClient)
	}
}
