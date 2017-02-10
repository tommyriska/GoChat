import socket
from threading import Thread
from colorama import *
import time
import sys
import os
import signal



#Startup message
print (Fore.LIGHTMAGENTA_EX + """//////////////////////////////////////////
//  _____        _____ _           _    //
// |  __ \      / ____| |         | |   //
// | |__) |   _| |    | |__   __ _| |_  //
// |  ___/ | | | |    | '_ \ / _` | __| //
// | |   | |_| | |____| | | | (_| | |_  //
// |_|    \__, |\_____|_| |_|\__,_|\__| //
//         __/ |                        //
//        |___/                         //
//////////////////////////////////////////""")
print (Style.RESET_ALL)
print " "

nickName = "SERVER"

# Command functions
def helpCommand():
    print "--Commands: ", "\n--!help Shows this list", "\n--!list Lists connected users", "\n--!quit Disconnects from the chat"

def listCommand():
    for con in connections:
        print str(con)

def quitCommand():
    if len(connections) >= 1:
        s.shutdown(socket.SHUT_RDWR)
    s.close()
    print "Server is shut down"
    os.kill(os.getppid(), signal.SIGHUP)

# Command list
commandList = {"!help": helpCommand, "!list": listCommand, "!quit": quitCommand}

def inputChecker(input):
    try:
        if input[0] == "!":
            if input.split(" ")[0] in commandList:
                commandList.get(input.split(" ")[0])()
            return True
        else:
            return False
    except IndexError:
        return False

def onNewClient(clientsocket, address):
    while True:
        data = clientsocket.recv(BUFFERSIZE)
        # Check for connection code
        if data == "1a92?#qQ,=11":
            # Add IP to connections
            if address not in connections:
                connections.append(address)
                print address, " connected"

        elif data == "10n0m0001":
            if address in connections:
                connections.remove(address)
                print address, " disconnected"

        # If not, print message
        else:
            print address, " ", data

            # Send message to all connections
            for c in connections:
                if c != address:
                    s.sendall(data, c)
                    s.close()
# Connections list
connections = []

# Server IP and socket
TCP_IP = ""
TCP_PORT = 5005
BUFFERSIZE = 20


# Create and bind socket
s = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
s.setsockopt(socket.SOL_SOCKET, socket.SO_REUSEADDR, 1)
s.bind((TCP_IP, TCP_PORT))
s.listen(50000)

# Some more informative messages
print "Server is running at ", TCP_IP, ": ", TCP_PORT
print "Waiting for connections.."

def inputListener():
    while True:
        message = raw_input(nickName + "-> ")
        if not inputChecker(message):
            s.sendall(nickName + ": " + message)

# Server loop
while True:
    c, address = s.accept()

    connectionThread = Thread(target=onNewClient(c, address))
    connectionThread.start()