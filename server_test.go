package main

import "testing"

/* test will check if the key array that is created
		on the server is the correct size (32 bytes) */
func TestCreateKey(t *testing.T) {

	key := createKey()
	var totalByte = 0

	for i := range key {
		fmt.Println(i)
		totalByte++
	}

	fmt.Println("Total byte i array --> ", totalByte)

	if totalByte < 32 {
		t.Error("Expected 32, got: ", totalByte)
	}
	if totalByte > 32 {
		t.Error("Expected 32. got: ", totalByte)
	}
}

/* test will check if a string variable is altered after
 		method call (encrypted) */
func TestEncrypt(t *testing.T) {

	key := make([]byte, 32)
	var crypt string = "test crypto"
	text := encrypt(key, crypt)

	if text == crypt {
		t.Error("Expected new text, got same.")
	}
}

/* test will check if a string variable is altered after
		method call (decrypted) */
func TestDecrypt(t *testing.T) {

	key := make([]byte, 32)
	var crypt string = "eL7Abmpadmk8nAaVcpT-a6MziwEKcL5z3ifS54SsnxA="
	text := decrypt(key, crypt)

	if text == crypt {
		t.Error("Expected new text, got same")
	}
}

/* test will check if a new room object is added correctly
	to the room array */
func TestMakeRoom(t *testing.T) {

	var n string = "testname"
	var m string = "testMessage"
	var b bool = false

	newRoom := Room{name: n, welcomeMsg: m}
	rooms = append(rooms, newRoom)

	for i := range rooms {
		if rooms[i] == newRoom {
			fmt.Println("New room added succesfully.")
			b = true
		}
	}
	if b != true {
		fmt.Println("New room not added to array.")
	}
}

/* test will check if lobby is empty at start up */
func TestLoadRooms(t *testing.T) {

	testList := loadRooms()
	if testList == nil {
		fmt.Println("Expected populated array, got an empty.")
	}
}

/* test will check if a new room is added to lobby array */
func TestSaveRoom(t *testing.T) {

	var s string = "testName"
	var w string = "testMsg"
	var d string = "testDesc"

	saveRoom(s, w, d)
	testList := loadRooms()

	for i := range testList {
		if testList[i].name == s {
			fmt.Println("New room succesfully added to array")
		}
	}
	fmt.Println("Room not added to array")
}
