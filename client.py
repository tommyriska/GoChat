import socket
import sys
from threading import Thread
from colorama import *
import time

# Server IP and socket
TCP_IP = "192.168.0.20"
TCP_PORT = 5006
BUFFERSIZE = 1024


# Create socket
s = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
s.connect((TCP_IP, TCP_PORT))
# Send connection code
s.send("1a92?#qQ,=11")

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
print "\nWelcome to fissechat!\nType !help to see commands.\n"
nickName = raw_input("Choose a nickname: ")

def helpCommand():
    print "--Commands: ",  "\n--!help shows this list", "\n--!list lists connected users", "\n--!name change your nickname", "\n--!quit quit"


def listCommand():
    print "list"


def nameCommand():
    global nickName
    oldName = nickName
    nickName = raw_input("Choose a nickname: ")
    s.sendall(oldName + " changed name to " + nickName)


def quitCommand():
    s.sendall("10n0m0001", (TCP_IP, TCP_PORT))
    sys.exit()

commandList = {"!help": helpCommand, "!list": listCommand, "!name": nameCommand, "!quit": quitCommand}

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


def chatListener():
    while True:
        try:
            data, address = s.recv(1024)
            print data
        except socket.error:
            return False

listenerThread = Thread(target=chatListener)
listenerThread.start()


def inputListener():
    while True:
        message = raw_input(nickName + "-> ")
        if not inputChecker(message):
            s.sendall(nickName + ": " + message)

inputThread = Thread(target=inputListener)
inputThread.start()
