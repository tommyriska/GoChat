import socket
from threading import Thread
from colorama import *
import time

#Startup message
print (Fore.LIGHTMAGENTA_EX + '//////////////////////////////////////////////////')
print ("//  ______ _                   _           _    //")
print ("// |  ____(_)                 | |         | |   //")
print ("// | |__   _ ___ ___  ___  ___| |__   __ _| |_  //")
print ("// |  __| | / __/ __|/ _ \/ __| '_ \ / _` | __| //")
print ("// | |    | \__ \__ \  __/ (__| | | | (_| | |_  //")
print ("// |_|    |_|___/___/\___|\___|_| |_|\__,_|\__| //")
print ("//                                              //")
print ("//////////////////////////////////////////////////")
print (Style.RESET_ALL)
print " "

nickName = "SERVER"

# Command functions
def helpCommand():
    print "--Commands: ", "\n--!help Shows this list", "\n--!list Lists connected users", "\n--!quit Disconnects from the chat"

def listCommand():
    printConn()

def quitCommand():
    print "Add some function here"
# Command list
commandList = {"!help": helpCommand, "!list": listCommand, "!quit": quitCommand}

def printConn():
    for c in connections:
        print address

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

def onNewClient(clientsocket, addr):
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
#Connections list
connections = []

# Server IP and socket
TCP_IP = ''
TCP_PORT = 5006
BUFFERSIZE = 20


# Create and bind socket
s = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
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


inputThread = Thread(target=inputListener)
inputThread.start()

#Server loop
while True:
    c, address = s.accept()
    connectionThread = Thread(target=onNewClient(c, address))
    connectionThread.start()


















# #Startup message
# print (Fore.LIGHTMAGENTA_EX + '//////////////////////////////////////////////////')
# print (Fore.LIGHTMAGENTA_EX + "//  ______ _                   _           _    //")
# print (Fore.LIGHTMAGENTA_EX + "// |  ____(_)                 | |         | |   //")
# print (Fore.LIGHTMAGENTA_EX + "// | |__   _ ___ ___  ___  ___| |__   __ _| |_  //")
# print (Fore.LIGHTMAGENTA_EX + "// |  __| | / __/ __|/ _ \/ __| '_ \ / _` | __| //")
# print (Fore.LIGHTMAGENTA_EX + "// | |    | \__ \__ \  __/ (__| | | | (_| | |_  //")
# print (Fore.LIGHTMAGENTA_EX + "// |_|    |_|___/___/\___|\___|_| |_|\__,_|\__| //")
# print (Fore.LIGHTMAGENTA_EX + "//                                              //")
# print (Fore.LIGHTMAGENTA_EX + "//////////////////////////////////////////////////")
# print (Style.RESET_ALL)
# print " "
#
# print "Server is running at ", host, ": ", port
#
# # List that holds connections
# connections = []

# # Server loop
# while True:
#     # Listen for message
#     data, address = s.recvfrom(1024)
#
#     # Check for connection code
#     if data == "1a92?#qQ,=11":
#         # Add IP to connections
#         if address not in connections:
#             connections.append(address)
#             print address, " connected"
#
#     elif data == "10n0m0001":
#         if address in connections:
#             connections.remove(address)
#             print address, " disconnected"
#
#     # If not, print message
#     else:
#         print address, " ", data
#
#         # Send message to all connections
#         for c in connections:
#             if c != address:
#                 s.sendto(data, c)
