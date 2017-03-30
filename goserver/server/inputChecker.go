package main

import "strings"

func checkForCmd(client Client, msg string) bool {
	b := false
	if len(msg) > 1 {
		words := strings.Split(msg[0:len(msg)-1], " ")
		switch words[0] {
		case "!room":
			b = true
			if len(words) > 1 {
				for _, element := range client.room.roomList {
					if element.name == words[1] {
						client.sendEncrypted("Switching room: " + words[1] + "\n")
						break
					}
				}
			} else {
				for _, element := range client.room.roomList {
					client.sendEncrypted(" - " + element.name + "\n")
				}
			}
		}
	}
	return b
}
